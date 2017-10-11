package model

import (
	"git.oschina.net/cnjack/novel-spider/config"
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

func MustGetDB() (*gorm.DB, error) {
	if db != nil {
		return db, nil
	}
	return Connect()
}

func Connect() (*gorm.DB, error) {
	var err error

	db, err = gorm.Open("mysql", config.GetMysqlConfig().DSN)
	if err != nil {
		return nil, err
	}
	return db, nil
}

func InitDB() error {
	tx := db.Begin()
	tx.LogMode(true)
	if err := tx.CreateTable(&Novel{}).Error; err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit().Error
}

func init() {
	_, err := Connect()
	if err != nil {
		panic(err)
	}
}
