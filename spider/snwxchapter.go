package spider

import (
	"errors"
	"net/url"
	"strings"

	"github.com/hu17889/go_spider/core/common/request"
)

type SnwxChapter struct {
	Url  string
	Data interface{}
}

func (s *SnwxChapter) Name() string {
	return "snwx.com"
}

func (s *SnwxChapter) Match(urlString string) bool {
	s.Url = urlString
	u, err := url.Parse(urlString)
	if err != nil {
		return false
	}
	if u.Host != "www.snwx.com" {
		return false
	}
	u.Path = strings.TrimRight(u.Path, ".html")
	u.Path = strings.Trim(u.Path, `/`)
	paths := strings.Split(u.Path, `/`)
	if len(paths) != 3 {
		return false
	}
	return false
}

func (s *SnwxChapter) Gain() (interface{}, error) {
	page := d.Download(request.NewRequest(s.Url, "html", "", "GET", "", nil, nil, nil, nil))
	if page.Errormsg() != "" {
		return "", errors.New(page.Errormsg())
	}
	doc := page.GetHtmlParser()
	//return doc.Find("div#BookText").Text(), nil
	html, err := doc.Find("div#BookText").Html()
	if err != nil {
		return "", nil
	}
	return html, nil
}
