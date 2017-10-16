package main

import (
	"log"
	"os"

	"gitee.com/cnjack/novel-spider/api"
	"gitee.com/cnjack/novel-spider/model"
)

func init() {
	log.SetOutput(os.Stdout)

	db := model.MustGetDB()
	db.LogMode(true)
	if err := db.CreateTable(&model.Novel{}).Error; err != nil {
		log.Println(db)
	}
}

func main() {
	api.Start()
}
