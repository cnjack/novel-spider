package httpserv

import (
	"fmt"
	"net/http"
	"strconv"

	"git.oschina.net/cnjack/novel-spider/model"
	"github.com/labstack/echo"
)

func getNovelDetails(c echo.Context) error {
	idString := c.Param("id")
	hasChapter := false
	if len(c.QueryParam("chapter")) > 0 {
		hasChapter = true
	}
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
	tag, err := model.FirstTagsByID(db, novel.TagID)
	if err != nil {
		novel.Style = "其他"
	}
	novel.Style = tag.TagName
	if novel == nil {
		return RecodeNotFound
	}
	return c.JSON(http.StatusOK, map[string]interface{}{
		"code": 0,
		"data": novel.Todata(hasChapter),
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
