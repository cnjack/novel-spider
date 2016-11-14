package model

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var db *gorm.DB

type PageOption struct {
	Page    int
	Count   int
	All     bool
	OrderBy string
}

var defaultPageOption = &PageOption{
	Page:    0,
	Count:   25,
	All:     false,
	OrderBy: "id desc",
}

func MustGetDB() (*gorm.DB, error) {
	if db != nil {
		return db, nil
	}
	return Connect()
}

func Connect() (*gorm.DB, error) {
	var err error
	db, err = gorm.Open("mysql", "root:root@tcp(127.0.0.1:3306)/novel?charset=utf8&parseTime=True&loc=Local")
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
	if err := tx.CreateTable(&Chapter{}).Error; err != nil {
		tx.Rollback()
		return err
	}
	if err := tx.CreateTable(&Task{}).Error; err != nil {
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
