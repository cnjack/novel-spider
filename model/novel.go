package model

import "github.com/jinzhu/gorm"

type Novel struct {
	gorm.Model
	Title        string      `sql:"title"`
	Auth         string      `sql:"auth"`
	Style        string      `sql:"style"`
	Status       NovelStatus `sql:"status"`
	Introduction string      `sql:"introduction" gorm:"type:text"`
	Chapter      string      `sql:"chapter" gorm:"type:longtext"`
	Url          string      `sql:"url"`
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

func FindNovelByAuth(db *gorm.DB, auth string, op *PageOption) (ns []*Novel, err error) {
	if op == nil {
		op = defaultPageOption
	}
	if err = db.Model(&Novel{}).Where("auth = ?", auth).Order(op.OrderBy).Limit(op.Count).Offset(op.Page * op.Count).Find(ns).Error; err != nil {
		return nil, err
	}
	return
}
