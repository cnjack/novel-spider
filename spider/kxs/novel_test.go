package kxs_test

import (
	"testing"

	"git.oschina.net/cnjack/novel-spider/spider"
	"git.oschina.net/cnjack/novel-spider/spider/kxs"
	"github.com/stretchr/testify/assert"
)

const testNovelUrl = "http://www.00kxs.com/html/26/26058/"

func TestSnwxNovel_Match(t *testing.T) {
	s := kxs.Novel{}
	b := s.Match(testNovelUrl)
	assert.Equal(t, true, b)
}

func TestSnwxNovel_Gain(t *testing.T) {
	s := kxs.Novel{}
	b := s.Match(testNovelUrl)
	novel, err := s.Gain()
	novelStruct, b2 := novel.(spider.Novel)
	if assert.NoError(t, err) {
		assert.Equal(t, true, b)
		assert.Equal(t, true, b2)
		assert.NotNil(t, novel)
		assert.Equal(t, "萧瑾瑜", novelStruct.Auth)
		assert.Equal(t, "http://www.00kxs.com/img/26/26058/26058s.jpg", novelStruct.Cover)
		assert.Equal(t, "天骄战纪", novelStruct.Title)
		assert.Equal(t, "玄幻魔法", novelStruct.Style)
		assert.Equal(t, "", novelStruct.Status)
		assert.NotEmpty(t, novelStruct.Introduction)
		assert.NotNil(t, novelStruct.Chapter)
	}
}
