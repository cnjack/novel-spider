package kxs_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"spider/spider/kxs"
)

const testChapterUrl = "http://www.00kxs.com/html/26/26058/11909975.html"

func TestSnwxChapter_Match(t *testing.T) {
	s := kxs.Chapter{}
	b := s.Match(testChapterUrl)
	assert.Equal(t, true, b)
}

func TestSnwxChapter_Gain(t *testing.T) {
	s := kxs.Chapter{}
	b := s.Match(testChapterUrl)
	chapter, err := s.Gain()
	chapterString, b2 := chapter.(string)
	if assert.NoError(t, err) {
		assert.Equal(t, true, b)
		assert.Equal(t, true, b2)
		assert.NotEmpty(t, chapterString)
	}
}
