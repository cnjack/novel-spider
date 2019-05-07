package api

import (
	"net/http"
	"spider/model"
	"spider/spider"
	"spider/spider/snwx"

	"github.com/labstack/echo"
)

type PostSearchParam struct {
	Title string `json:"title"`
}

func PostSearchLocal(c echo.Context) error {
	postSearchParam := &PostSearchParam{}
	err := c.Bind(postSearchParam)
	if err != nil {
		return ParamError
	}
	op := c.Get(PageOptionKey).(*model.PageOption)
	if postSearchParam.Title == "" {
		return ParamError
	}
	db := model.MustGetDB()
	novels, err := model.SearchByTitleOrAuth(db, postSearchParam.Title, postSearchParam.Title, op)
	if err != nil {
		return ServerError
	}
	nextPage := 0
	if len(novels) >= op.Count {
		nextPage = op.Page + 1
	}
	return c.JSON(http.StatusOK, struct {
		Code int                  `json:"code"`
		Data []*model.SearchNovel `json:"data"`
		Next int                  `json:"next"`
	}{
		Code: 0,
		Data: novels,
		Next: nextPage,
	})
}

func PostSearchRemote(c echo.Context) error {
	postSearchParam := &PostSearchParam{}
	err := c.Bind(postSearchParam)
	if err != nil {
		return ParamError
	}
	if postSearchParam.Title == "" {
		return ParamError
	}
	searchers := []spider.Spider{
		&snwx.Search{},
	}
	var data = []*spider.Search{}
	for _, s := range searchers {
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
	return c.JSON(http.StatusOK, struct {
		Code int `json:"code"`
		Data []*spider.Search
	}{
		Code: 0,
		Data: data,
	})
}
