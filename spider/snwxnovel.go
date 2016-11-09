package spider

import (
	"errors"
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/hu17889/go_spider/core/common/request"
)

type SnwxNovel struct {
	Url    string
	BookID string
	Data   interface{}
}

func (s *SnwxNovel) Name() string {
	return "snwx.com"
}

func (s *SnwxNovel) Match(urlString string) bool {
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
	if len(paths) == 0 {
		return false
	}
	if len(paths) == 3 && paths[0] == "book" {
		s.BookID = paths[1] + "/" + paths[2]
		if err != nil {
			return false
		}
		return true
	}
	return false
}

func (snwx *SnwxNovel) Gain() (interface{}, error) {
	urlString := "http://www.snwx.com/book/" + snwx.BookID + "/"
	page := d.Download(request.NewRequest(urlString, "html", "", "GET", "", nil, nil, nil, nil))
	if page.Errormsg() != "" {
		return "", errors.New(page.Errormsg())
	}
	doc := page.GetHtmlParser()
	var novel Novel
	novel.Title = doc.Find("div .infotitle h1").Text()
	doc.Find("div .infotitle i1").Each(func(i int, s *goquery.Selection) {
		ss := strings.Split(s.Text(), "：")
		if len(ss) == 2 && i == 0 {
			novel.Auth = ss[1]
		}
		if len(ss) == 2 && i == 1 {
			novel.Style = s.Text()
		}
	})

	is := strings.Split(doc.Find(".intro").Text(), "：")
	if len(is) == 3 {
		iss := strings.Split(is[1], "\n")
		if len(iss) > 1 {
			novel.Introduction = iss[0] + iss[1]
		}
	}

	doc.Find("div#list dl dd").Each(func(i int, s *goquery.Selection) {
		cp := &Chapter{}
		cp.Title = s.Find("a").Text()

		from, b := s.Find("a").Attr("href")
		if !b {
			return
		}
		cp.From = urlString + from
		novel.Chapter = append(novel.Chapter, cp)
	})
	return novel, nil
}
