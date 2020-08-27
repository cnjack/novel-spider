package kxs

import (
	"net/url"
	"strings"

	"spider/internal/downloader"
)

type Chapter struct {
	Url  *url.URL
	Data interface{}
}

func (s *Chapter) Name() string {
	return "00kxs.com"
}

func (s *Chapter) Match(urlString string) bool {
	u, err := url.Parse(urlString)
	s.Url = u
	if err != nil {
		return false
	}
	if strings.Index(u.Host, "00kxs.com") == -1 {
		return false
	}
	path := strings.TrimRight(u.Path, ".html")
	path = strings.Trim(path, `/`)
	paths := strings.Split(path, `/`)
	if len(paths) != 4 {
		return false
	}
	if paths[0] == "html" {
		return true
	}
	return false
}

func (s *Chapter) Gain() (interface{}, error) {
	d := downloader.NewHttpDownloaderFromUrl(s.Url).Download()
	if err := d.Error(); err != nil {
		return "", err
	}
	doc, err := d.Resource().Document()
	if err != nil {
		return "", err
	}
	html, err := doc.Find("div#content").Html()
	if err != nil {
		return "", nil
	}
	return html, nil
}
