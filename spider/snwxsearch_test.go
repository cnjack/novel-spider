package spider_test

import (
	"testing"

	"git.oschina.net/cnjack/novel-spider/spider"
	"github.com/stretchr/testify/assert"
)

const testSearch = "校园绝品狂徒"

func TestSnwxSearch_Match(t *testing.T) {
	s := spider.SnwxSearch{}
	b := s.Match(testSearch)
	assert.Equal(t, true, b)
}

func TestSnwxSearch_Gain(t *testing.T) {
	s := spider.SnwxSearch{}
	b := s.Match(testSearch)
	searchs, err := s.Gain()
	searchsStruct, b2 := searchs.([]*spider.Search)
	if assert.NoError(t, err) {
		assert.Equal(t, true, b)
		assert.Equal(t, true, b2)
		assert.NotNil(t, searchs)
		assert.Equal(t, "http://www.snwx.com/book/0/381/", searchsStruct[0].From)
		assert.Equal(t, "校园绝品狂徒", searchsStruct[0].Novel.Title)
		assert.Equal(t, "连载中", searchsStruct[0].Novel.Status)
	}
}
