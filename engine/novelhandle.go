package engine

import (
	"git.oschina.net/cnjack/novel-spider/model"
	"github.com/labstack/echo"
	"net/http"
	"strconv"
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
	if chapter == nil {
		return RecodeNotFound
	}
	return c.JSON(http.StatusOK, map[string]interface{}{
		"code": 0,
		"data": chapter.Todata(),
	})
}
