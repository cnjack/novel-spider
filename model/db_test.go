package model_test

import (
	"testing"

	"git.oschina.net/cnjack/novel-spider/model"
)

func TestInitDB(t *testing.T) {
	err := model.InitDB()
	if err != nil {
		t.Error(err)
	}
}
