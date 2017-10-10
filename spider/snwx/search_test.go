package snwx_test

import (
	"testing"

	"git.oschina.net/cnjack/novel-spider/spider/snwx"
	"github.com/stretchr/testify/assert"
	"git.oschina.net/cnjack/novel-spider/spider"
)

const testSearch = "校园绝品狂徒"

func TestSnwxSearch_Match(t *testing.T) {
	s := snwx.Search{}
	b := s.Match(testSearch)
	assert.Equal(t, true, b)
}

func TestSnwxSearch_Gain(t *testing.T) {
	s := snwx.Search{}
	b := s.Match(testSearch)
	searchs, err := s.Gain()

	if assert.NoError(t, err) {
		searchsStruct, b2 := searchs.([]*spider.Search)
		assert.Equal(t, true, b)
		assert.Equal(t, true, b2)
		assert.NotNil(t, searchs)
		assert.Equal(t, "http://www.snwx.com/book/0/381/", searchsStruct[0].From)
		assert.Equal(t, "校园绝品狂徒", searchsStruct[0].Novel.Title)
		assert.Equal(t, "连载中", searchsStruct[0].Novel.Status)
	}

}
