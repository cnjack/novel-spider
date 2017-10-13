package main

import (
	"log"
	"os"

	"git.oschina.net/cnjack/novel-spider/api"
	"git.oschina.net/cnjack/novel-spider/model"
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
