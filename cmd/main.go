package main

import (
	"log"
	"os"
	"spider/internal/api"
	"spider/internal/config"
	"spider/internal/repository"
)

func init() {
	log.SetOutput(os.Stdout)

	config.Init()
	_, err := repository.InitDatabase()
	if err != nil {
		panic(err)
	}
	repository.InitRedis()
	db := repository.MustGetDB()
	db.LogMode(true)
	if err := db.CreateTable(&repository.Novel{}).Error; err != nil {
		log.Println(db)
	}
	db.LogMode(false)
}

func main() {
	api.Start()
}
