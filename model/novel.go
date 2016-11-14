package model

import (
	"encoding/json"
	"github.com/jinzhu/gorm"
	"time"
)

type Novel struct {
	gorm.Model
	Title        string      `sql:"title" json:"title"`
	Auth         string      `sql:"auth" json:"auth"`
	Style        string      `sql:"style" json:"style"`
	Status       NovelStatus `sql:"status" json:"status"`
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

func FindNovelByAuth(db *gorm.DB, auth string, op *PageOption) (ns []*Novel, err error) {
	if op == nil {
		op = defaultPageOption
	}
	if err = db.Model(&Novel{}).Where("auth = ?", auth).Order(op.OrderBy).Limit(op.Count).Offset(op.Page * op.Count).Find(ns).Error; err != nil {
		return nil, err
	}
	return
}

func (n *Novel) Todata() interface{} {
	resp := map[string]interface{}{
		"id":           n.ID,
		"create_at":    n.CreatedAt.Format(time.RFC3339),
		"title":        n.Title,
		"auth":         n.Auth,
		"style":        n.Style,
		"status":       n.Status,
		"introduction": n.Introduction,
		"url":          n.Url,
	}
	return resp
}

type NovelChapter struct {
	Title     string `json:"title"`
	Index     uint   `json:"index"`
	ChapterID uint   `json:"chapter_id"`
	Url       string `json:"url"`
}

func (n *Novel) ChapterTodata() interface{} {
	novelChapters := []NovelChapter{}
	if n.Chapter != "" {
		err := json.Unmarshal([]byte(n.Chapter), &novelChapters)
		if err != nil {
			return err
		}
	}
	return novelChapters
}
