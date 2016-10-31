package model

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
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
	db, err = gorm.Open("sqlite3", "data.db?_txlock=exclusive")
	if err != nil {
		return nil, err
	}
	return db, nil
}

func InitDB() error {
	tx := db.Begin()
	if err := tx.CreateTable(&Novel{}).Error; err != nil {
		tx.Rollback()
		return err
	}
	if err := tx.CreateTable(&Task{}).Error; err != nil {
		tx.Rollback()
		return err
	}
	if err := tx.CreateTable(&Chapter{}).Error; err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit().Error
}

func init() {
	Connect()
}
