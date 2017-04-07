package httpserv

import (
	"fmt"
	"net/http"
	"strconv"

	"git.oschina.net/cnjack/novel-spider/model"
	"github.com/labstack/echo"
	"git.oschina.net/cnjack/novel-spider/spider"
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
	tag, err := model.FirstTagsByID(db, novel.TagID)
	if err != nil {
		novel.Style = "其他"
	} else {
		novel.Style = tag.TagName
	}

	if novel == nil {
		return RecodeNotFound
	}
	chapters, err := novel.ChapterTodata()
	if err != nil {
		return ServerError
	}
	if len(*chapters) == 0 {
		return ServerError
	}
	return c.JSON(http.StatusOK, struct {
		Code     int              `json:"code"`
		Data     *model.NovelData `json:"data"`
		FirstCid uint             `json:"first_cid"`
	}{
		Code:     0,
		Data:     novel.Todata(false),
		FirstCid: (*chapters)[0].ChapterID,
	})
}

func deleteNovel(c echo.Context) error {
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
	chapters, err := novel.ChapterTodata()
	if err != nil {
		return ServerError
	}
	return c.JSON(http.StatusOK, struct {
		Code int                   `json:"code"`
		Data *[]model.NovelChapter `json:"data"`
	}{
		Code: 0,
		Data: chapters,
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

	if chapter.Data == "" {
		searchs := []spider.Spider{
			&spider.SnwxChapter{},
		}
		for _,s := range searchs {
			if s.Match(chapter.Url) {
				sRespInterface, err := s.Gain()
				if err != nil {
					return err
				}
				sResp, ok := sRespInterface.(string)
				if !ok || sResp == "" {
					return ServerError
				}
				chapter.Data = sResp
				//插入数据库
				go func(id uint, data string) {
					model.UpdateChapterDataByID(db, id, data)
				}(chapter.ID, sResp)
			}
		}
	}
	var next, prev = 0, 0
	var pages = []model.Chapter{}
	//TODO 性能问题 需要添加索引
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
	return c.JSON(http.StatusOK, struct {
		Code int                `json:"code"`
		Data *model.ChapterData `json:"data"`
		Next int                `json:"next"`
		Prev int                `json:"prev"`
	}{
		Code: 0,
		Data: chapter.Todata(),
		Next: next,
		Prev: prev,
	})
}
