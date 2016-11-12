package spider

import (
	"net/url"
	"strings"

	"git.oschina.net/cnjack/downloader"
)

type SnwxChapter struct {
	Url  *url.URL
	Data interface{}
}

func (s *SnwxChapter) Name() string {
	return "snwx.com"
}

func (s *SnwxChapter) Match(urlString string) bool {
	u, err := url.Parse(urlString)
	s.Url = u
	if err != nil {
		return false
	}
	if u.Host != "www.snwx.com" {
		return false
	}
	path := strings.TrimRight(u.Path, ".html")
	path = strings.Trim(path, `/`)
	paths := strings.Split(path, `/`)
	if len(paths) != 4 {
		return false
	}
	if paths[0] == "book" {
		return true
	}
	return false
}

func (s *SnwxChapter) Gain() (interface{}, error) {
	d := downloader.NewHttpDownloaderFromUrl(s.Url).Download()
	if err := d.Error(); err != nil {
		return "", err
	}
	doc, err := d.Resource().Document()
	if err != nil {
		return "", err
	}
	html, err := doc.Find("div#BookText").Html()
	if err != nil {
		return "", nil
	}
	return html, nil
}
