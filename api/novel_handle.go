package api

import (
	"net/http"
	"strconv"

	"git.oschina.net/cnjack/novel-spider/model"
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
