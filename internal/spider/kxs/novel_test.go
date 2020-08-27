package kxs_test

import (
	"spider/internal/spider"
	"spider/internal/spider/kxs"
	"testing"

	"github.com/stretchr/testify/assert"
)

const testNovelUrl = "http://www.00kxs.com/html/23/23770/"

func TestKxsNovel_Match(t *testing.T) {
	s := kxs.Novel{}
	b := s.Match(testNovelUrl)
	assert.Equal(t, true, b)
}

func TestKxsNovel_Gain(t *testing.T) {
	s := kxs.Novel{}
	b := s.Match(testNovelUrl)
	novel, err := s.Gain()
	novelStruct, b2 := novel.(spider.Novel)
	if assert.NoError(t, err) {
		assert.Equal(t, true, b)
		assert.Equal(t, true, b2)
		assert.NotNil(t, novel)
		assert.Equal(t, "无赖是我", novelStruct.Auth)
		assert.Equal(t, "https://img00kxs.cdn.bcebos.com/img/23/23770/23770s.jpg", novelStruct.Cover)
		assert.Equal(t, "僵尸无赖", novelStruct.Title)
		assert.Equal(t, "", novelStruct.Style)
		assert.Equal(t, "", novelStruct.Status)
		assert.NotEmpty(t, novelStruct.Introduction)
		assert.NotNil(t, novelStruct.Chapter)
		t.Log(novelStruct.Chapter[0].Title)
		t.Log(novelStruct.Chapter[0].From)
	}
}
