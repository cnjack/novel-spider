package engine

import (
	"net/http"
	"net/url"
	"strconv"

	"fmt"

	"git.oschina.net/cnjack/novel-spider/model"
	"git.oschina.net/cnjack/novel-spider/spider"
	"github.com/jinzhu/gorm"
	"github.com/labstack/echo"
)

func getNovelDetails(c echo.Context) error {
	idString := c.Param("id")
	id, err := strconv.Atoi(idString)
	if err != nil || id == 0 {
		return ParamError
	}
	db, err := model.MustGetDB()
	if err != nil {
		return ServerError
	}
	novel, err := model.FirstNovelByID(db, uint(id))
	if err != nil {
		return ServerError
	}
	if novel == nil {
		return RecodeNotFound
	}
	return c.JSON(http.StatusOK, map[string]interface{}{
		"code": 0,
		"data": novel.Todata(),
	})
}

func getNovels(c echo.Context) error {
	db, err := model.MustGetDB()
	if err != nil {
		return ServerError
	}
	op := c.Get(PageOptionKey).(*model.PageOption)
	novels, err := model.FindNovels(db, op)
	if err != nil {
		return ServerError
	}
	if novels == nil {
		return RecodeNotFound
	}
	var data = []interface{}{}
	for _, v := range novels {
		data = append(data, v.Todata())
	}
	nextPage := 0
	if len(novels) >= op.Count {
		nextPage = op.Page + 1
	}
	return c.JSON(http.StatusOK, map[string]interface{}{
		"code": 0,
		"next": nextPage,
		"data": data,
	})
}

func getNovelChapters(c echo.Context) error {
	idString := c.Param("id")
	id, err := strconv.Atoi(idString)
	if err != nil || id == 0 {
		return ParamError
	}
	db, err := model.MustGetDB()
	if err != nil {
		return ServerError
	}
	novel, err := model.FirstNovelByID(db, uint(id))
	if err != nil {
		return ServerError
	}
	if novel == nil {
		return RecodeNotFound
	}
	return c.JSON(http.StatusOK, map[string]interface{}{
		"code": 0,
		"data": novel.ChapterTodata(),
	})
}

func getChapter(c echo.Context) error {
	idString := c.Param("id")
	id, err := strconv.Atoi(idString)
	if err != nil || id == 0 {
		return ParamError
	}
	db, err := model.MustGetDB()
	if err != nil {
		return ServerError
	}
	chapter, err := model.FirstChapterByID(db, uint(id))
	if err != nil {
		return ServerError
	}
	var next, prev = 0, 0
	var pages = []model.Chapter{}
	if err := db.Where("`index` IN (?, ?) AND `novel_id` = ?", chapter.Index-1, chapter.Index+1, chapter.NovelID).Find(&pages).Error; err != nil {
		fmt.Println(err)
		return ServerError
	}
	if len(pages) == 2 {
		prev = int(pages[0].ID)
		next = int(pages[1].ID)
	} else if len(pages) == 1 {
		if chapter.Index == 0 {
			next = int(pages[0].ID)
		} else {
			prev = int(pages[0].ID)
		}
	}
	if chapter == nil {
		return RecodeNotFound
	}
	return c.JSON(http.StatusOK, map[string]interface{}{
		"code": 0,
		"data": chapter.Todata(),
		"next": next,
		"prev": prev,
	})
}

type PostSearchParam struct {
	Title string `json:"title"`
}

func postSearch(c echo.Context) error {
	postSearchParam := &PostSearchParam{}
	err := c.Bind(postSearchParam)
	if err != nil {
		return ParamError
	}
	if postSearchParam.Title == "" {
		return ParamError
	}
	searchs := []spider.Spider{
		&spider.SnwxSearch{},
	}
	var data = []*spider.Search{}
	for _, s := range searchs {
		if s.Match(postSearchParam.Title) {
			sRespInterface, err := s.Gain()
			if err != nil {
				return err
			}
			sResp, ok := sRespInterface.([]*spider.Search)
			if !ok {
				return ServerError
			}
			for _, v := range sResp {
				data = append(data, v)
			}
		}
	}
	return c.JSON(http.StatusOK, map[string]interface{}{
		"code": 0,
		"data": data,
	})
}

type PostTaskParam struct {
	Url string `json:"url"`
}

func postNovelTask(c echo.Context) error {
	postTaskParam := &PostTaskParam{}
	if err := c.Bind(postTaskParam); err != nil {
		return ParamError
	}
	_, err := url.Parse(postTaskParam.Url)
	if err != nil {
		return ParamError
	}
	searchs := []spider.Spider{
		&spider.SnwxNovel{},
	}
	match := false
	for _, s := range searchs {
		if s.Match(postTaskParam.Url) {
			match = true
		}
	}
	if !match {
		return ParamError
	}
	db, err := model.MustGetDB()
	if err != nil {
		return ServerError
	}
	task := &model.Task{
		TType:  model.NovelTask,
		Url:    postTaskParam.Url,
		Status: model.TaskStatusPrepare,
		Times:  -1,
	}
	var exist = &model.Task{}
	err = db.Model(exist).Where("url = ?", postTaskParam.Url).First(exist).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return ServerError
	}
	if err != gorm.ErrRecordNotFound {
		return TaskIsRepeated
	}
	if err := db.Model(task).Create(task).Error; err != nil {
		return ServerError
	}
	return c.JSON(http.StatusOK, map[string]interface{}{
		"code": 0,
		"data": db.RowsAffected,
	})
}
