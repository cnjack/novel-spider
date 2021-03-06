package repository

import (
	"spider/internal/config"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var db *gorm.DB

type PageOption struct {
	Page  int
	Count int
	All   bool
	Sort  string
}

var defaultPageOption = &PageOption{
	Page:  0,
	Count: 25,
	All:   false,
	Sort:  "desc",
}

func MustGetDB() *gorm.DB {
	return db
}

func InitDatabase() (*gorm.DB, error) {
	var err error

	db, err = gorm.Open("mysql", config.GetConfig().MysqlConfig.DSN)
	if err != nil {
		return nil, err
	}
	return db, nil
}

func InitDB() {
	db.LogMode(true)
	db.CreateTable(&Novel{})
	db.Model(&Novel{}).AddIndex("idx_style", "style")
	db.Model(&Novel{}).AddIndex("idx_title", "title")
	db.Model(&Novel{}).AddIndex("idx_auth", "auth")
}
