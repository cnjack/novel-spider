package snwx_test

import (
	"testing"

	"spider/spider"
	"spider/spider/snwx"

	"github.com/stretchr/testify/assert"
)

const testNovelUrl = "http://www.snwx8.com/book/392/392460/"

func TestSnwxNovel_Match(t *testing.T) {
	s := snwx.Novel{}
	b := s.Match(testNovelUrl)
	assert.Equal(t, true, b)
}

func TestSnwxNovel_Gain(t *testing.T) {
	s := snwx.Novel{}
	b := s.Match(testNovelUrl)
	novel, err := s.Gain()
	novelStruct, b2 := novel.(spider.Novel)
	if assert.NoError(t, err) {
		assert.Equal(t, true, b)
		assert.Equal(t, true, b2)
		assert.NotNil(t, novel)
		assert.Equal(t, "唐家三少", novelStruct.Auth)
		assert.Equal(t, "https://www.snwx8.com/files/article/image/392/392460/392460s.jpg", novelStruct.Cover)
		assert.Equal(t, "斗罗大陆Ⅳ终极斗罗", novelStruct.Title)
		assert.Equal(t, "玄幻", novelStruct.Style)
		assert.Equal(t, "连载中", novelStruct.Status)
		assert.NotEmpty(t, novelStruct.Introduction)
		assert.NotNil(t, novelStruct.Chapter)
	}
}
