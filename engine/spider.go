package engine

import "errors"

import (
	"encoding/json"
	"log"
	"sync"
	"time"

	"git.oschina.net/cnjack/novel-spider/config"
	"git.oschina.net/cnjack/novel-spider/model"
	"git.oschina.net/cnjack/novel-spider/spider"
	"git.oschina.net/cnjack/novel-spider/tool"
	"github.com/jinzhu/gorm"
)

var w = sync.WaitGroup{}

func Spider() {
	if config.GetSpiderConfig().StopSingle {
		return
	}
	for i := 0; i < config.GetSpiderConfig().MaxProcess; i++ {
		w.Add(1)
		go RunATask()
	}
	w.Wait()
}

func RunATask() {
	defer func() {
		w.Done()
	}()

	t, err := getTask()
	if err != nil {
		if err != gorm.ErrRecordNotFound {
			log.Printf("ERROR: getTask ERROR; ERRstring: %s", err.Error())
		}
		//
		time.Sleep(30 * time.Second)
		RunATask()
	}
	log.Printf("INFO: getTask ok; task id: %d", t.ID)
	runTask(t)
	if config.GetSpiderConfig().StopSingle {
		return
	}
	RunATask()
}

func runTask(t *model.Task) (err error) {
	defer func() {
		if err != nil {
			t.ChangeTaskStatus(model.TaskStatusFail)
		}
	}()
	spiders := []spider.Spider{
		&spider.SnwxChapter{},
		&spider.SnwxNovel{},
	}
	for _, s := range spiders {
		if !s.Match(t.Url) {
			continue
		}
		resp, err := s.Gain()
		if err != nil {
			log.Printf("ERROR: Gain Task ERROR; ERRstring: %s; task id: %d", err.Error(), t.ID)
			return err
		}
		//更新任务状态
		err = flashTask(t, resp)
		if err != nil {
			log.Printf("ERROR: flashTask ERROR; ERRstring: %s; task id: %d", err.Error(), t.ID)
			return err
		}
		log.Printf("INFO: runTask OK;task id: %d", t.ID)
		return nil
	}
	return errors.New("have not match spider")
}

func flashTask(t *model.Task, data interface{}) (err error) {
	switch t.TType {
	case model.NovelTask:
		err = flashNovelTask(t, data)
	case model.ChapterTask:
		err = flashChapterTask(t, data)
	default:
		return errors.New("unknown task")
	}
	return nil
}

type NovelChapters []NovelChapter

func (n NovelChapters) Has(from string) bool {
	for _, v := range n {
		if v.Url == from {
			return true
		}
	}
	return false
}

type NovelChapter struct {
	Title     string `json:"title"`
	Index     uint   `json:"index"`
	ChapterID uint   `json:"chapter_id"`
	Url       string `json:"url"`
}

func flashNovelTask(t *model.Task, data interface{}) error {
	novel, ok := data.(spider.Novel)
	if !ok {
		return errors.New("get the data error")
	}
	dbNovel := &model.Novel{}
	db, err := model.MustGetDB()
	if err != nil {
		return err
	}
	if err := db.Model(dbNovel).Where("title = ? AND url = ? AND auth = ?", novel.Title, novel.From, novel.Auth).Find(dbNovel).Error; err != nil {
		if err != gorm.ErrRecordNotFound {
			return err
		}
	}
	if dbNovel.ID == 0 {
		cover, err := tool.UploadFromUrl(novel.Cover)
		if err != nil {
			cover = novel.Cover
		}
		dbNovel = &model.Novel{
			Title:        novel.Title,
			Auth:         novel.Auth,
			Cover:        cover,
			Style:        novel.Style,
			Status:       model.String2NovelStatus(novel.Status),
			Introduction: novel.Introduction,
			Url:          novel.From,
		}
		if err := db.Model(dbNovel).Create(dbNovel).Error; err != nil {
			return err
		}
	}
	NovelChapters := NovelChapters{}
	if dbNovel.Chapter != "" {
		err := json.Unmarshal([]byte(dbNovel.Chapter), &NovelChapters)
		if err != nil {
			return err
		}
	}
	NewNovelChapters := []NovelChapter{}
	//对比之前的章节
	if len(NovelChapters) != len(novel.Chapter) {
		for _, c := range novel.Chapter {
			if !NovelChapters.Has(c.From) {
				ncp := &model.Chapter{
					NovelID: dbNovel.ID,
					Index:   c.Index,
					Title:   c.Title,
					Status:  0,
					Url:     c.From,
				}

				if err := db.Model(ncp).Create(ncp).Error; err != nil {
					return err
				}
				ntask := &model.Task{
					TType:    model.ChapterTask,
					Url:      c.From,
					Status:   model.TaskStatusPrepare,
					Times:    -1,
					TargetID: ncp.ID,
				}
				//创建新的任务
				if err := db.Model(ntask).Create(ntask).Error; err != nil {
					return err
				}
				NewNovelChapters = append(NewNovelChapters, NovelChapter{
					Title:     c.Title,
					Index:     c.Index,
					ChapterID: ncp.ID,
					Url:       c.From,
				})
			}
		}
	}
	for _, nc := range NewNovelChapters {
		NovelChapters = append(NovelChapters, nc)
	}
	NovelChaptersJsonString, _ := json.Marshal(NovelChapters)
	dbNovel.Chapter = string(NovelChaptersJsonString)
	//更新章节
	if err := db.Model(dbNovel).Update(dbNovel).Error; err != nil {
		return err
	}
	return t.ChangeTaskStatus(model.TaskStatusOk)
}

func flashChapterTask(t *model.Task, data interface{}) error {
	chapterString, ok := data.(string)
	if !ok {
		return errors.New("get the data error")
	}
	db, err := model.MustGetDB()
	if err != nil {
		return err
	}
	if err := db.Model(&model.Chapter{}).Where("id = ?", t.TargetID).Update(map[string]interface{}{"data": chapterString, "status": 1}).Error; err != nil {
		return err
	}
	return t.ChangeTaskStatus(model.TaskStatusOk)
}

func getTask() (task *model.Task, err error) {
	db, err := model.MustGetDB()
	if err != nil {
		return nil, err
	}
	tx := db.Begin()
	if err = tx.Error; err != nil {
		return nil, err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()
	task = &model.Task{}
	if err = db.Raw("SELECT * FROM tasks WHERE status IN (?, ?) AND deleted_at IS NULL ORDER BY status, id ASC LIMIT 0, 1 FOR UPDATE", model.TaskStatusPrepare, model.TaskStatusFail).Scan(task).Error; err != nil {
		return nil, err
	}
	if err = tx.Model(&model.Task{}).Where("id = ?", task.ID).Updates(map[string]interface{}{"status": model.TaskStatusRunning}).Error; err != nil {
		return nil, err
	}
	return
}
