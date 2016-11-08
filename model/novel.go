package model

import "github.com/jinzhu/gorm"

type Novel struct {
	gorm.Model
	Title        string `sql:"title"`
	Auth         string `sql:"auth"`
	Style        string `sql:"style"`
	Status       int    `sql:"status"`
	Introduction string `sql:"introduction"`

	From string `sql:"from"`
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
