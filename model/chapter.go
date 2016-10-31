package model

import "github.com/jinzhu/gorm"

type Chapter struct {
	gorm.Model
	NovelID int64  `sql:"novel_id"`
	Title   string `sql:"title"`
	Data    string `sql:"data"`
	From    string `sql:"from"`
}
