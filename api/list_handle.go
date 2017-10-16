package api

import (
	"net/http"

	"gitee.com/cnjack/novel-spider/model"
	"github.com/labstack/echo"
)

func GetNovels(c echo.Context) error {
	db := model.MustGetDB()
	op := c.Get(PageOptionKey).(*model.PageOption)
	novels, err := model.FindNovels(db, op)
	if err != nil {
		return ServerError
	}
	if novels == nil {
		return RecodeNotFound
	}
	var data = []*model.NovelData{}
	if err != nil {
		return ServerError
	}
	for _, v := range novels {
		data = append(data, v.Todata(false))
	}
	nextPage := 0
	if len(novels) >= op.Count {
		nextPage = op.Page + 1
	}
	return c.JSON(http.StatusOK, struct {
		Code int                `json:"code"`
		Next int                `json:"next"`
		Data []*model.NovelData `json:"data"`
	}{
		Code: 0,
		Next: nextPage,
		Data: data,
	})
}

func GetStyles(c echo.Context) error {
	db := model.MustGetDB()
	tags, err := model.GetStyle(db)
	if err != nil {
		return ServerError
	}
	return c.JSON(http.StatusOK, struct {
		Code int      `json:"code"`
		Data []string `json:"data"`
	}{
		Code: 0,
		Data: tags,
	})
}

func GetStyleNovels(c echo.Context) error {
	style := c.Param("style")
	if len(style) < 0 {
		return ParamError
	}
	db := model.MustGetDB()
	op := c.Get(PageOptionKey).(*model.PageOption)
	novels, err := model.FindNovelsWithStyle(db, style, op)
	if err != nil {
		return ServerError
	}
	if novels == nil {
		return RecodeNotFound
	}
	var data = []*model.NovelData{}
	for _, v := range novels {
		data = append(data, v.Todata(false))
	}
	nextPage := 0
	if len(novels) >= op.Count {
		nextPage = op.Page + 1
	}
	return c.JSON(http.StatusOK, struct {
		Code int                `json:"code"`
		Next int                `json:"next"`
		Data []*model.NovelData `json:"data"`
	}{
		Code: 0,
		Next: nextPage,
		Data: data,
	})
}
