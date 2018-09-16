package main

import (
	"log"
	"os"

	"spider/api"
	"spider/model"
)

func init() {
	log.SetOutput(os.Stdout)

	db := model.MustGetDB()
	db.LogMode(true)
	if err := db.CreateTable(&model.Novel{}).Error; err != nil {
		log.Println(db)
	}
	db.LogMode(false)
}

func main() {
	api.Start()
}
