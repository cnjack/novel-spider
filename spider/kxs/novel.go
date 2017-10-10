package kxs

import (
	"net/url"
	"strings"

	"git.oschina.net/cnjack/downloader"
	"git.oschina.net/cnjack/novel-spider/spider"
	"github.com/PuerkitoBio/goquery"
)

type Novel struct {
	Url             *url.URL
	BookID          string
	StyleMap        *map[string]string
	WithOutChapters bool
	Data            interface{}
}

func (s *Novel) Name() string {
	return "00kxs.com"
}

func (s *Novel) Match(urlString string) bool {
	u, err := url.Parse(urlString)
	if err != nil {
		return false
	}
	s.Url = u
	if u.Host != "www.00kxs.com" {

		return false
	}
	path := strings.TrimRight(u.Path, ".html")
	path = strings.Trim(path, `/`)
	paths := strings.Split(path, `/`)
	if len(paths) == 0 {
		return false
	}
	if len(paths) == 3 && paths[0] == "html" {
		s.BookID = paths[1] + "/" + paths[2]
		if err != nil {
			return false
		}
		return true
	}
	return false
}

func (snwx *Novel) Gain() (interface{}, error) {
	u, _ := url.Parse("http://www.00kxs.com/html/" + snwx.BookID + "/")
	d := downloader.NewHttpDownloaderFromUrl(u).Download()
	if err := d.Error(); err != nil {
		return "", err
	}
	doc, err := d.Resource().Document()
	if err != nil {
		return "", err
	}
	var novel spider.Novel
	novel.Title = doc.Find("div #info h1").Text()
	novel.Cover, _ = doc.Find("#fmimg img").Attr("src")
	if novel.Cover == "/modules/article/images/nocover.jpg" {
		novel.Cover = ""
	}
	novel.Cover = "http://www.00kxs.com" + novel.Cover
	novel.From = u.String()
	doc.Find("#info p").Each(func(i int, s *goquery.Selection) {
		ss := strings.Split(s.Text(), "：")
		strings.Trim(ss[1], " ")
		sss := ss[1]
		if len(ss) == 2 && i == 0 {
			novel.Auth = sss
		}
	})
	doc.Find(".con_top a").Each(func(i int, s *goquery.Selection) {
		if i == 7 {
			novel.Style = s.Text()
		}
	})
	if snwx.StyleMap != nil {
		style, ok := (*snwx.StyleMap)[novel.Style]
		if ok {
			novel.Style = style
		} else {
			novel.Style = "其他"
		}
	}
	introString, err := doc.Find("#intro").Html()
	if err != nil {
		return nil, err
	}
	novel.Introduction = introString

	if !snwx.WithOutChapters {
		doc.Find("div#list ul li").Each(func(i int, s *goquery.Selection) {
			cp := &spider.Chapter{}
			cp.Title = s.Find("a").Text()
			cp.Index = uint(i)
			from, b := s.Find("a").Attr("href")
			if !b {
				return
			}
			cp.From = u.String() + from
			novel.Chapter = append(novel.Chapter, cp)
		})
	}
	return novel, nil
}
