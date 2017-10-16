package snwx_test

import (
	"testing"

	"gitee.com/cnjack/novel-spider/spider/snwx"
	"github.com/stretchr/testify/assert"
)

const testChapterUrl = "http://www.snwx.com/book/0/381/155205.html"

func TestSnwxChapter_Match(t *testing.T) {
	s := snwx.Chapter{}
	b := s.Match(testChapterUrl)
	assert.Equal(t, true, b)
}

func TestSnwxChapter_Gain(t *testing.T) {
	s := snwx.Chapter{}
	b := s.Match(testChapterUrl)
	chapter, err := s.Gain()
	chapterString, b2 := chapter.(string)
	if assert.NoError(t, err) {
		assert.Equal(t, true, b)
		assert.Equal(t, true, b2)
		assert.NotEmpty(t, chapterString)
	}
}
