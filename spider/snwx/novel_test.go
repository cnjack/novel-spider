package snwx_test

import (
	"testing"

	"git.oschina.net/cnjack/novel-spider/spider/snwx"
	"github.com/stretchr/testify/assert"
	"git.oschina.net/cnjack/novel-spider/spider"
)

const testNovelUrl = "http://www.snwx8.com/book/0/760/"

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
		assert.Equal(t, "风凌天下", novelStruct.Auth)
		assert.Equal(t, "http://www.snwx8.com/files/article/image/0/760/760s.jpg", novelStruct.Cover)
		assert.Equal(t, "傲世九重天", novelStruct.Title)
		assert.Equal(t, "玄幻", novelStruct.Style)
		assert.Equal(t, "连载中", novelStruct.Status)
		assert.NotEmpty(t, novelStruct.Introduction)
		assert.NotNil(t, novelStruct.Chapter)
	}
}
