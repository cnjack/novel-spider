package httpserv

import (
	"net/http"
	"net/url"

	"git.oschina.net/cnjack/novel-spider/model"
	"git.oschina.net/cnjack/novel-spider/spider"
	"github.com/jinzhu/gorm"
	"github.com/labstack/echo"
)

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
