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
