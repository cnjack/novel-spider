package model_test

import (
	"testing"

	"git.oschina.net/cnjack/novel-spider/model"
)

func TestMustGetDB(t *testing.T) {
	db, err := model.MustGetDB()
	if err != nil {
		t.Error(err)
	}
	task := &model.Task{
		TType:  model.NovelTask,
		Url:    "http://www.snwx.com/book/0/381/",
		Status: model.TaskStatusPrepare,
		Times:  -1,
	}
	c := db.Model(task).Create(task)
	if c.Error != nil {
		t.Error(c.Error)
	}
}

func TestInitDB(t *testing.T) {
	err := model.InitDB()
	if err != nil {
		t.Error(err)
	}
}
