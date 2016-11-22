package httpserv

import (
	"net/http"
	"strconv"

	"git.oschina.net/cnjack/novel-spider/model"
	"github.com/labstack/echo"
)

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
	tags, err := model.GetTags(db)
	if err != nil {
		return ServerError
	}
	for k, v := range novels {
		for _, vv := range *tags {
			if vv.ID == v.TagID {
				novels[k].Style = vv.TagName
				break
			}
		}
		if novels[k].Style == "" {
			novels[k].Style = "其他"
		}
	}
	for _, v := range novels {
		data = append(data, v.Todata(false))
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

func getStyles(c echo.Context) error {
	db, err := model.MustGetDB()
	if err != nil {
		return ServerError
	}
	tags, err := model.GetTags(db)
	if err != nil {
		return ServerError
	}
	return c.JSON(http.StatusOK, map[string]interface{}{
		"code": 0,
		"data": tags,
	})
}

func getStyleNovels(c echo.Context) error {
	styleString := c.Param("style")
	styleID, err := strconv.Atoi(styleString)
	if err != nil {
		return ParamError
	}
	if styleID <= 0 {
		return ParamError
	}
	db, err := model.MustGetDB()
	if err != nil {
		return ServerError
	}
	tag, err := model.FirstTagsByID(db, styleID)
	if err != nil {
		return ServerError
	}
	op := c.Get(PageOptionKey).(*model.PageOption)
	novels, err := model.FindNovelsWithStyle(db, tag.ID, op)
	if err != nil {
		return ServerError
	}
	if novels == nil {
		return RecodeNotFound
	}
	var data = []interface{}{}
	tags, err := model.GetTags(db)
	if err != nil {
		return ServerError
	}
	for k, v := range novels {
		for _, vv := range *tags {
			if vv.ID == v.TagID {
				novels[k].Style = vv.TagName
				break
			}
		}
		if novels[k].Style == "" {
			novels[k].Style = "其他"
		}
	}
	for _, v := range novels {
		data = append(data, v.Todata(false))
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
