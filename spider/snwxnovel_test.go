package spider_test

import (
	"log"
	"testing"

	"git.oschina.net/cnjack/novel-spider/spider"
	"github.com/stretchr/testify/assert"
)

const testNovelUrl = "http://www.snwx.com/book/0/381/"

func TestSnwxNovel_Match(t *testing.T) {
	s := spider.SnwxNovel{}
	b := s.Match(testNovelUrl)
	assert.Equal(t, true, b)
}

func TestSnwxNovel_Gain(t *testing.T) {
	s := spider.SnwxNovel{}
	b := s.Match(testNovelUrl)
	novel, err := s.Gain()
	novelStruct, b2 := novel.(spider.Novel)
	log.Println(novelStruct.Introduction)
	if assert.NoError(t, err) {
		assert.Equal(t, true, b)
		assert.Equal(t, true, b2)
		assert.NotNil(t, novel)
		assert.Equal(t, "柳江南", novelStruct.Auth)
		assert.Equal(t, "http://www.snwx.com/files/article/image/0/381/381s.jpg", novelStruct.Cover)
		assert.Equal(t, "校园绝品狂徒", novelStruct.Title)
		assert.Equal(t, "其他小说", novelStruct.Style)
		assert.Equal(t, "连载中", novelStruct.Status)
		assert.NotEmpty(t, novelStruct.Introduction)
		assert.NotNil(t, novelStruct.Chapter)
	}
}
