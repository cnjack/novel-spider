package api

import (
	"net/http"
	"net/url"
	"spider/spider"
	"spider/spider/snwx"
	"strconv"

	"spider/model"

	"github.com/labstack/echo"
)

func GetNovelDetails(c echo.Context) error {
	idString := c.Param("id")
	id, err := strconv.Atoi(idString)
	if err != nil || id == 0 {
		return ParamError
	}
	db := model.MustGetDB()
	novel, err := model.FirstNovelByID(db, uint(id))
	if err != nil {
		return ServerError
	}
	if novel.Style != "" {
		novel.Style = "其他"
	}

	if novel == nil {
		return RecodeNotFound
	}
	return c.JSON(http.StatusOK, struct {
		Code     int              `json:"code"`
		Data     *model.NovelData `json:"data"`
		FirstCid uint             `json:"first_cid"`
	}{
		Code: 0,
		Data: novel.Todata(false),
	})
}

func GetNovelDetailsFromUrl(c echo.Context) error {
	urlString := c.Param("url")
	urlString, err := url.QueryUnescape(urlString)
	if err != nil {
		return ParamError
	}
	if _, err := url.Parse(urlString); err != nil {
		return ParamError
	}
	db := model.MustGetDB()
	novel, err := model.FirstNovelByUrl(db, urlString)
	if err != nil && novel != nil {
		return c.JSON(http.StatusOK, struct {
			Code     int              `json:"code"`
			Data     *model.NovelData `json:"data"`
			FirstCid uint             `json:"first_cid"`
		}{
			Code: 0,
			Data: novel.Todata(false),
		})
	}
	novel, err = model.GetNovelFromUrl(urlString)
	if err != nil {
		return ServerError
	}
	return c.JSON(http.StatusOK, struct {
		Code     int              `json:"code"`
		Data     *model.NovelData `json:"data"`
		FirstCid uint             `json:"first_cid"`
	}{
		Code: 0,
		Data: novel.Todata(false),
	})
}

func DeleteNovel(c echo.Context) error {
	idString := c.Param("id")
	id, err := strconv.Atoi(idString)
	if err != nil || id == 0 {
		return ParamError
	}
	db := model.MustGetDB()
	novel, err := model.FirstNovelByID(db, uint(id))
	if err != nil {
		return ServerError
	}
	if novel == nil {
		return RecodeNotFound
	}
	//处理删除任务
	tx := db.Begin()
	if tx.Error != nil {
		return ServerError
	}
	defer func() {
		if err != nil {
			tx.Commit()
		} else {
			tx.Rollback()
		}
	}()
	if err = db.Exec("DELETE FROM novels WHERE id = ?", novel.ID).Error; err != nil {
		return ServerError
	}
	if err = db.Exec("DELETE FROM chapters WHERE novel_id = ?", novel.ID).Error; err != nil {
		return ServerError
	}
	if err = db.Exec("DELETE FROM tasks WHERE t_type = 0 AND url = ?", novel.Url).Error; err != nil {
		return ServerError
	}
	return c.JSON(http.StatusOK, struct {
		Code int    `json:"code"`
		Data string `json:"data"`
	}{
		Code: 0,
		Data: "ok",
	})
}

type PostAddNovelParam struct {
	URL string `json:"url"`
}

func AddNovel(c echo.Context) error {
	postAddNovelParam := &PostAddNovelParam{}
	err := c.Bind(postAddNovelParam)
	if err != nil {
		return ParamError
	}
	if postAddNovelParam.URL == "" {
		return ParamError
	}
	searchers := []spider.Spider{
		&snwx.Novel{},
	}
	novelData := spider.Novel{}
	var ok bool
	for _, s := range searchers {
		if s.Match(postAddNovelParam.URL) {
			sRespInterface, err := s.Gain()
			if err != nil {
				return GainError
			}
			novelData, ok = sRespInterface.(spider.Novel)
			if !ok {
				return SpiderError
			}
			break
		}
	}
	if novelData.From == "" {
		return GainEmptyError
	}
	novel := &model.Novel{
		Title:        novelData.Title,
		Auth:         novelData.Auth,
		Style:        novelData.Style,
		Status:       novelData.Status,
		Cover:        novelData.Cover,
		Introduction: novelData.Introduction,
		Url:          novelData.From,
	}
	db := model.MustGetDB()
	err = novel.Add(db)
	if err != nil {
		return ServerError
	}
	return c.JSON(http.StatusOK, struct {
		Code     int              `json:"code"`
		Data     *model.NovelData `json:"data"`
		FirstCid uint             `json:"first_cid"`
	}{
		Code: 0,
		Data: novel.Todata(false),
	})
}

func GetNovelChapters(c echo.Context) error {
	idString := c.Param("id")
	id, err := strconv.Atoi(idString)
	if err != nil || id == 0 {
		return ParamError
	}
	db := model.MustGetDB()
	chapters, err := model.FirstChapterByID(db, uint(id))
	if err != nil {
		return ServerError
	}

	if chapters == nil {
		return RecodeNotFound
	}
	return c.JSON(http.StatusOK, struct {
		Code int         `json:"code"`
		Data interface{} `json:"data"`
	}{
		Code: 0,
		Data: chapters,
	})
}
