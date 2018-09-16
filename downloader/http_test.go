package downloader_test

import (
	"net/url"
	"testing"

	"git.oschina.net/cnjack/downloader"
	"github.com/stretchr/testify/assert"
)

func TestNewHttpDownloaderFromUrl(t *testing.T) {
	url, err := url.Parse("https://blog.nightc.com/")
	assert.NoError(t, err)
	d := downloader.NewHttpDownloaderFromUrl(url)
	assert.Nil(t, d.Download().Error())
	doc, err := d.Resource().Document()
	if assert.NoError(t, err) {
		assert.Equal(t, "Nightc", doc.Find("title").Text())
	}

}
