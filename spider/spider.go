package spider

import (
	"encoding/json"
	"errors"
	"log"
	"sync"
	"time"

	"git.oschina.net/cnjack/novel-spider/model"
	"github.com/hu17889/go_spider/core/downloader"
	"github.com/jinzhu/gorm"
)

var MaxRunningTask = 20 //最大线程数
var w = sync.WaitGroup{}
var StopSingle = false

var d = downloader.NewHttpDownloader()

type Novel struct {
	Title        string
	Auth         string
	Style        string
	Introduction string
	From         string
	Status       string
	Chapter      []*Chapter
}

type Chapter struct {
	Novel *Novel
	Index uint
	Title string
	From  string
	Data  string
}

type Search struct {
	Name       string
	SearchName string
	From       string
}

type Spider interface {
	Name() string
	Match(string) bool
	Gain() (interface{}, error)
}

var spiders []Spider

func Run() {
	for i := 0; i < MaxRunningTask; i++ {
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
		log.Println("get task err:", err)
		return
	}
	if t == nil {
		//如果没有任务就暂停一下
		time.Sleep(30 * time.Second)
		RunATask()
	}
	err = runTask(t)
	if StopSingle {
		return
	}
	RunATask()
}

func runTask(t *model.Task) error {
	for _, s := range spiders {
		if !s.Match(t.Url) {
			continue
		}
		log.Printf("task id %d is running", t.ID)
		resp, err := s.Gain()
		if err != nil {
			return err
		}
		log.Printf("task id %d Gain down", t.ID)
		//更新任务状态
		err = flashTask(t, resp)
		if err != nil {
			log.Printf("flashTask ID: %d error.%s", t.ID, err.Error())
			return err
		}
		log.Println(s.Name(), " deal ", t.Url, " done")
		return nil
	}
	return errors.New("have not match spider")
}

func flashTask(t *model.Task, data interface{}) (err error) {
	switch t.TType {
	case model.NovelTask:
		flashNovelTask(t, data)
	case model.ChapterTask:
		flashChapterTask(t, data)
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
	log.Println("flash novel task start id", t.ID)
	novel, ok := data.(Novel)
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
			log.Println("find novel error", err)
			return err
		}
	}
	if dbNovel.ID == 0 {
		dbNovel = &model.Novel{
			Title:        novel.Title,
			Auth:         novel.Auth,
			Style:        novel.Style,
			Status:       model.String2NovelStatus(novel.Status),
			Introduction: novel.Introduction,
			Url:          novel.From,
		}
		if err := db.Model(dbNovel).Create(dbNovel).Error; err != nil {
			return err
		}
		log.Println("insert novel id", dbNovel.ID)
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
	log.Println("chapterString", chapterString)
	if err := db.Model(&model.Chapter{}).Where("id = ?", t.TargetID).Update(map[string]interface{}{"data": chapterString, "status": 1}).Error; err != nil {
		log.Println("update Novel err", chapterString)
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
	task, err = model.FisrtTask(db)
	if err != nil {
		return nil, err
	}
	if task == nil {
		return nil, nil
	}
	err = tx.Model(&model.Task{}).Where("id = ?", task.ID).Updates(map[string]interface{}{"status": model.TaskStatusRunning}).Error
	if err != nil {
		return nil, err
	}
	return
}

func init() {
	spiders = []Spider{
		&SnwxChapter{},
		&SnwxNovel{},
	}
}
