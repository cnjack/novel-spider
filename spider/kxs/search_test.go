package kxs_test

import (
	"testing"

	"spider/spider"
	"spider/spider/kxs"

	"github.com/stretchr/testify/assert"
)

const testSearch = "大主宰"

func TestSnwxSearch_Match(t *testing.T) {
	s := kxs.Search{}
	b := s.Match(testSearch)
	assert.Equal(t, true, b)
}

func TestSnwxSearch_Gain(t *testing.T) {
	s := kxs.Search{}
	b := s.Match(testSearch)
	searchs, err := s.Gain()

	if assert.NoError(t, err) {
		searchsStruct, b2 := searchs.([]*spider.Search)
		assert.Equal(t, true, b)
		assert.Equal(t, true, b2)
		assert.NotNil(t, searchs)
		assert.Equal(t, "http://www.00kxs.com/html/0/59/", searchsStruct[0].From)
		assert.Equal(t, "大主宰", searchsStruct[0].Novel.Title)
	}

}
