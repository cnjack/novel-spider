package spider

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/hu17889/go_spider/core/common/request"
	"github.com/hu17889/go_spider/core/downloader"
)

var d = downloader.NewHttpDownloader()

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
	u.Path = strings.Trim(u.Path, `/`)
	paths := strings.Split(u.Path, `/`)
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
	fmt.Println(paths)
	return false
}

func (snwx *Snwx) Gain() (interface{}, error) {
	urlString := "http://www.snwx.com/txt/" + strconv.Itoa(snwx.BookID) + ".html"
	page := d.Download(request.NewRequest(urlString, "html", "", "GET", "", nil, nil, nil, nil))
	doc := page.GetHtmlParser()
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
