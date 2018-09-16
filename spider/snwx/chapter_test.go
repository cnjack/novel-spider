package snwx_test

import (
	"spider/spider/snwx"
	"testing"

	"github.com/stretchr/testify/assert"
)

const testChapterUrl = "https://www.snwx8.com/book/0/381/155205.html"

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
