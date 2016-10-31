package spider

import (
	"github.com/PuerkitoBio/goquery"
	"net/url"
	"strconv"
	"strings"
)

type Snwx struct {
	Url    string
	BookID int
	Data   interface{}
}

func (s *Snwx) Name() string {
	return "snwx.com"
}

func (s *Snwx) Match(urlString string) bool {
	s.Url = urlString
	u, err := url.Parse(urlString)
	if err != nil {
		return false
	}
	if u.Host != "www.snwx.com" {
		return false
	}
	u.Path = strings.TrimRight(u.Path, ".html")
	paths := strings.Split(u.Path, "/")
	if len(paths) == 0 {
		return false
	}
	if len(paths) == 3 && paths[0] == "book" {
		s.BookID, err = strconv.Atoi(paths[2])
		if err != nil {
			return false
		}
		return true
	}
	if len(paths) == 2 && paths[0] == "txt" {
		s.BookID, err = strconv.Atoi(paths[1])
		if err != nil {
			return false
		}
		return true
	}
	return false
}

func (snwx *Snwx) Gain() (interface{}, error) {
	doc, err := goquery.NewDocument("http://www.snwx.com/txt/" + strconv.Itoa(snwx.BookID) + ".html")
	if err != nil {
		return nil, err
	}
	var novel Novel
	novel.Title = strings.TrimRight(doc.Find("title").Text(), "TXT下载")
	doc.Find("#Chapters ul li").Each(func(i int, s *goquery.Selection) {
		cp := &Chapter{}
		cp.Title = s.Find("a").Text()
		from, b := s.Find("a").Attr("href")
		if !b {
			return
		}
		cp.From = from
		novel.Chapter = append(novel.Chapter, cp)
	})
	return novel, nil
}
