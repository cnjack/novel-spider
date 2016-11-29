package httpserv

import (
	"net/http"
	"net/url"

	"git.oschina.net/cnjack/novel-spider/job"
	"git.oschina.net/cnjack/novel-spider/model"
	"git.oschina.net/cnjack/novel-spider/spider"
	"github.com/labstack/echo"
)

type PostSearchParam struct {
	Title string `json:"title"`
}

func postSearchLocal(c echo.Context) error {
	postSearchParam := &PostSearchParam{}
	err := c.Bind(postSearchParam)
	if err != nil {
		return ParamError
	}
	op := c.Get(PageOptionKey).(*model.PageOption)
	if postSearchParam.Title == "" {
		return ParamError
	}
	db, err := model.MustGetDB()
	if err != nil {
		return ServerError
	}
	novels, err := model.SearchByTitleOrAuth(db, postSearchParam.Title, postSearchParam.Title, op)
	if err != nil {
		return ServerError
	}
	nextPage := 0
	if len(novels) >= op.Count {
		nextPage = op.Page + 1
	}
	return c.JSON(http.StatusOK, map[string]interface{}{
		"code": 0,
		"data": novels,
		"next": nextPage,
	})
}

func postSearchRemote(c echo.Context) error {
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
	var exist = 0
	if err := db.Model(&model.Novel{}).Where("url = ?", postTaskParam.Url).Count(&exist).Error; err != nil {
		return ServerError
	}
	if exist > 0 {
		return TaskIsRepeated
	}
	job.PublishTask(task)
	return c.JSON(http.StatusOK, map[string]interface{}{
		"code": 0,
		"data": "ok",
	})
}
