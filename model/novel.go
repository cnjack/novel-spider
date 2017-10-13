package model

import (
	"encoding/json"
	"time"

	"net/url"

	"github.com/jinzhu/gorm"
)

type Novel struct {
	gorm.Model
	Title        string `sql:"title" json:"title"`
	Auth         string `sql:"auth" json:"auth"`
	Style        string `sql:"style" json:"style"`
	Status       string `sql:"status" json:"status"`
	Cover        string `sql:"cover" json:"cover"`
	Introduction string `sql:"introduction" gorm:"type:text" json:"intrduction"`
	Chapter      string `sql:"chapter" gorm:"type:longtext" json:"-"`
	Url          string `sql:"url" json:"from"`
}

type SearchNovel struct {
	ID    uint   `sql:"id" json:"id"`
	Title string `sql:"title" json:"title"`
	Auth  string `sql:"auth" json:"auth"`
}

func SearchByTitleOrAuth(db *gorm.DB, title, auth string, op *PageOption) (*[]SearchNovel, error) {
	var ns []SearchNovel
	var err error
	if op == nil {
		op = defaultPageOption
	}
	if err = db.Table("novels").Where("title LIKE ? OR auth = ?", "%"+title+"%", auth).Select([]string{"title", "id", "auth"}).Limit(op.Count).Offset(op.Page * op.Count).Order("id desc").Find(&ns).Error; err != nil {
		return nil, err
	}
	return &ns, nil
}

func (n *Novel) Add(db *gorm.DB) error {
	return db.Create(n).Error
}

func FirstNovelByID(db *gorm.DB, id uint) (n *Novel, err error) {
	n = &Novel{}
	if err = db.Model(n).Where("id = ?", id).First(n).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return
}

func FirstNovelByIDWithoutChapters(db *gorm.DB, id uint) (n *Novel, err error) {
	n = &Novel{}
	if err = db.Model(n).Select([]string{"id", "title", "auth", "tag_id", "cover", "status", "introduction", "url"}).Where("id = ?", id).First(n).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return
}

func FindNovelByAuth(db *gorm.DB, auth string, op *PageOption) (ns []Novel, err error) {
	if op == nil {
		op = defaultPageOption
	}
	if err = db.Model(&Novel{}).Where("auth = ?", auth).Order("id " + op.Sort).Limit(op.Count).Offset(op.Page * op.Count).Find(&ns).Error; err != nil {
		return nil, err
	}
	return
}

func FindNovels(db *gorm.DB, op *PageOption) (ns []Novel, err error) {
	if op == nil {
		op = defaultPageOption
	}
	if err = db.Model(&Novel{}).Order("id " + op.Sort).Limit(op.Count).Offset(op.Page * op.Count).Find(&ns).Error; err != nil {
		return nil, err
	}
	return
}

func FindNovelsWithStyle(db *gorm.DB, style string, op *PageOption) (ns []Novel, err error) {
	if op == nil {
		op = defaultPageOption
	}
	if err = db.Model(&Novel{}).Where("tag_id LIKE '%" + style + "%'").Order("id " + op.Sort).Limit(op.Count).Offset(op.Page * op.Count).Find(&ns).Error; err != nil {
		return nil, err
	}
	return
}

func FirstChapterByID(db *gorm.DB, id uint) ([]*NovelChapter, error) {
	n := &Novel{}
	if err := db.Model(&Novel{}).Select([]string{"chapters"}).Where("id = ?", id).Limit(1).Find(n).Error; err != nil {
		return nil, err
	}
	return n.ChapterTodata()
}

type NovelData struct {
	ID           uint            `json:"id"`
	CreateAt     string          `json:"create_at"`
	Title        string          `json:"title"`
	Auth         string          `json:"auth"`
	Style        string          `json:"style"`
	Status       string          `json:"status"`
	Cover        string          `json:"cover"`
	Introduction string          `json:"introduction"`
	Chapter      []*NovelChapter `json:"chapters"`
	Url          string          `json:"from"`
}

func (n *Novel) Todata(more bool) *NovelData {
	cover, err := url.Parse(n.Cover)
	coverStr := ""
	if err != nil || cover.Host != "spider-img.nightc.com" {
		coverStr = "http://spider-img.nightc.com/cover.jpg"
	} else {
		coverStr = cover.String()
	}
	resp := NovelData{
		ID:           n.ID,
		CreateAt:     n.CreatedAt.Format(time.RFC3339),
		Title:        n.Title,
		Auth:         n.Auth,
		Style:        n.Style,
		Status:       n.Status,
		Introduction: n.Introduction,
		Cover:        coverStr,
		Url:          n.Url,
	}
	if more {
		resp.Chapter, _ = n.ChapterTodata()
	}
	return &resp
}

type NovelChapter struct {
	Title string `json:"title"`
	Index uint   `json:"index"`
	Url   string `json:"url"`
}

func (n *Novel) ChapterTodata() ([]*NovelChapter, error) {
	novelChapters := make([]*NovelChapter, 0)
	if n.Chapter != "" {
		err := json.Unmarshal([]byte(n.Chapter), &novelChapters)
		if err != nil {
			return nil, err
		}
	}
	return novelChapters, nil
}
