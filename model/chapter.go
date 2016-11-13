package model

import "github.com/jinzhu/gorm"

type Chapter struct {
	gorm.Model
	NovelID uint   `sql:"novel_id"`
	Index   uint   `sql:"index"`
	Title   string `sql:"title"`
	Data    string `sql:"data" gorm:"type:longtext"`
	Status  uint8  `sql:"status"`
	Url     string `sql:"url"`
}
