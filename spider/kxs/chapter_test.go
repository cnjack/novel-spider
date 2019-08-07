package kxs_test

import (
	"testing"

	"spider/spider/kxs"

	"github.com/stretchr/testify/assert"
)

const testChapterUrl = "http://www.00kxs.com/html/23/23770/"

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
