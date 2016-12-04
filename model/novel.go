package model

import (
	"encoding/json"
	"time"

	"github.com/jinzhu/gorm"
)

type Novel struct {
	gorm.Model
	Title        string      `sql:"title" json:"title"`
	Auth         string      `sql:"auth" json:"auth"`
	Style        string      `sql:"-" json:"style"`
	TagID        int         `sql:"tag_id" json:"tag_id"`
	Status       NovelStatus `sql:"status" json:"status"`
	Cover        string      `sql:"cover" json:"cover"`
	Introduction string      `sql:"introduction" gorm:"type:text" json:"intrduction"`
	Chapter      string      `sql:"chapter" gorm:"type:longtext" json:"-"`
	Url          string      `sql:"url" json:"from"`
}

type NovelStatus uint8

const (
	NovelSerializing NovelStatus = iota
	NovelCompleted
)

func String2NovelStatus(statusString string) NovelStatus {
	if statusString == "连载中" {
		return NovelSerializing
	}
	if statusString == "已完成" {
		return NovelCompleted
	}
	return NovelCompleted
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

func (s NovelStatus) Tostring() string {
	switch s {
	case NovelSerializing:
		return "连载中"
	case NovelCompleted:
		return "已完成"
	}
	return "未知"
}

func (n *Novel) Add(db *gorm.DB) error {
	return db.Create(n).Error
}

func (user *Novel) BeforeCreate(scope *gorm.Scope) error {
	return nil
}

func CountNovel() (count int, err error) {
	err = db.Model(&Novel{}).Count(&count).Error
	return
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

func FindNovelsWithStyle(db *gorm.DB, styleID int, op *PageOption) (ns []Novel, err error) {
	if op == nil {
		op = defaultPageOption
	}
	if err = db.Model(&Novel{}).Where("tag_id = ?", styleID).Order("id " + op.Sort).Limit(op.Count).Offset(op.Page * op.Count).Find(&ns).Error; err != nil {
		return nil, err
	}
	return
}

type NovelData struct {
	ID           uint            `json:"id"`
	CreateAt     string          `json:"create_at"`
	Title        string          `json:"title"`
	Auth         string          `json:"auth"`
	Style        string          `json:"style"`
	Status       NovelStatus     `json:"status"`
	Cover        string          `json:"cover"`
	Introduction string          `json:"intrduction"`
	Chapter      *[]NovelChapter `json:"chapters"`
	Url          string          `json:"from"`
}

func (n *Novel) Todata(more bool) *NovelData {
	resp := NovelData{
		ID:           n.ID,
		CreateAt:     n.CreatedAt.Format(time.RFC3339),
		Title:        n.Title,
		Auth:         n.Auth,
		Style:        n.Style,
		Status:       n.Status,
		Introduction: n.Introduction,
		Cover:        n.Cover,
		Url:          n.Url,
	}
	if more {
		resp.Chapter, _ = n.ChapterTodata()
	}
	return &resp
}

type NovelChapter struct {
	Title     string `json:"title"`
	Index     uint   `json:"index"`
	ChapterID uint   `json:"chapter_id"`
	Url       string `json:"url"`
}

func (n *Novel) ChapterTodata() (*[]NovelChapter, error) {
	novelChapters := []NovelChapter{}
	if n.Chapter != "" {
		err := json.Unmarshal([]byte(n.Chapter), &novelChapters)
		if err != nil {
			return nil, err
		}
	}
	return &novelChapters, nil
}
