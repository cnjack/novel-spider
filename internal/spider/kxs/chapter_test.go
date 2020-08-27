package kxs_test

import (
	"spider/internal/spider/kxs"
	"testing"

	"github.com/stretchr/testify/assert"
)

const testChapterUrl = "http://www.00kxs.com/html/23/23770/2773140.html"

func TestKxsChapter_Match(t *testing.T) {
	s := kxs.Chapter{}
	b := s.Match(testChapterUrl)
	assert.Equal(t, true, b)
}

func TestKxsChapter_Gain(t *testing.T) {
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
