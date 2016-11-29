package job

import "errors"

import (
	"encoding/json"
	"log"

	"git.oschina.net/cnjack/novel-spider/model"
	"git.oschina.net/cnjack/novel-spider/spider"
	"git.oschina.net/cnjack/novel-spider/tool"
	"github.com/jinzhu/gorm"
)

func Spider() {
	go Run()
	UpdateNovelTask()
}

var StyleMap = map[string]string{
	"玄幻小说": "玄幻",
	"修真小说": "修真",
	"都市小说": "都市",
	"穿越小说": "穿越",
	"网游小说": "网游",
	"科幻小说": "科幻",
	"其他小说": "其他",
}

func RunTask(t *model.Task) error {
	spiders := []spider.Spider{
		&spider.SnwxChapter{},
		&spider.SnwxNovel{StyleMap: &StyleMap},
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
	return err
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

var stylemap map[string]int

func flashNovelTask(t *model.Task, data interface{}) (err error) {
	var db *gorm.DB
	defer func() {
		if err != nil {
			t.ChangeTaskStatus(model.TaskStatusFail)
		}
	}()
	novel, ok := data.(spider.Novel)
	if !ok {
		return errors.New("get the data error")
	}
	dbNovel := &model.Novel{}
	db, err = model.MustGetDB()
	if err != nil {
		return err
	}
	if stylemap == nil {
		var tags *[]model.Tags
		tags, err = model.GetTags(db)
		if err != nil {
			return err
		}
		stylemap = map[string]int{}
		for _, v := range *tags {
			stylemap[v.TagName] = v.ID
		}
	}

	if t.TargetID == 0 {
		cover, err := tool.UploadFromUrl(novel.Cover)
		if err != nil {
			cover = novel.Cover
		}
		dbNovel = &model.Novel{
			Title:        novel.Title,
			Auth:         novel.Auth,
			Cover:        cover,
			Status:       model.String2NovelStatus(novel.Status),
			Introduction: novel.Introduction,
			Url:          novel.From,
			TagID:        0,
		}
		if stylemap != nil {
			tag, ok := stylemap[dbNovel.Style]
			if ok {
				dbNovel.TagID = tag
			}
		}
		dbNovel.Status = model.NovelCompleted
		if novel.Status == "连载中" {
			dbNovel.Status = model.NovelSerializing
		}
		if err = db.Model(dbNovel).Create(dbNovel).Error; err != nil {
			return err
		}
	}
	NovelChapters := NovelChapters{}
	if dbNovel.Chapter != "" {
		err = json.Unmarshal([]byte(dbNovel.Chapter), &NovelChapters)
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
				q.PutNoWait(ntask)
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
	if err = db.Model(dbNovel).Select([]string{"update_at", "chapter"}).Update(dbNovel).Error; err != nil {
		return err
	}
	if t.ID != 0 {
		if dbNovel.Status == model.NovelCompleted {
			return t.ChangeTaskStatus(model.TaskStatusOk)
		} else {
			return t.ChangeTaskStatus(model.TaskStatusPrepare)
		}
	} else {
		if dbNovel.Status != model.NovelCompleted {
			t.TargetID = dbNovel.ID
			err = db.Model(&model.Task{}).Create(t).Error
			if err != nil {
				log.Println("create task err", err)
			}
			return nil
		}
	}

	return nil
}

func PublishTask(t *model.Task) {
	q.PutNoWait(t)
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
	return nil
}
