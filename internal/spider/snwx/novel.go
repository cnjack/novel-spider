package snwx

import (
	"net/url"
	"strings"

	"spider/internal/downloader"
	"spider/internal/spider"

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
	return "snwx3.com"
}

func (s *Novel) Match(urlString string) bool {
	u, err := url.Parse(urlString)
	if err != nil {
		return false
	}
	s.Url = u
	if u.Host != "www.snwx.com" && u.Host != "www.snwx3.com" {
		return false
	}
	path := strings.TrimRight(u.Path, ".html")
	path = strings.Trim(path, `/`)
	paths := strings.Split(path, `/`)
	if len(paths) == 0 {
		return false
	}
	if len(paths) == 2 && paths[0] == "txt" {
		d := downloader.NewHttpDownloaderFromUrl(u).Download()
		if err := d.Error(); err != nil {
			return false
		}
		doc, err := d.Resource().Document()
		if err != nil {
			return false
		}
		nurl := ""
		doc.Find("a").EachWithBreak(func(i int, d *goquery.Selection) bool {
			if i == 1 {
				u, ok := d.Attr("href")
				if !ok {
					return true
				}
				nurl = u
				return false
			}
			return true
		})
		return s.Match(nurl)
	}
	if (len(paths) == 3 || len(paths) == 4) && paths[0] == "book" {
		s.BookID = paths[1] + "/" + paths[2]
		if err != nil {
			return false
		}
		return true
	}
	return false
}

func (snwx *Novel) Gain() (interface{}, error) {
	u, _ := url.Parse("https://www.snwx3.com/book/" + snwx.BookID + "/")
	d := downloader.NewHttpDownloaderFromUrl(u).Download()
	if err := d.Error(); err != nil {
		return "", err
	}
	doc, err := d.Resource().Document()
	if err != nil {
		return "", err
	}
	var novel spider.Novel
	novel.Title = doc.Find("div .infotitle h1").Text()
	novel.Cover, _ = doc.Find("#fmimg img").Attr("src")
	if novel.Cover == "/modules/article/images/nocover.jpg" {
		novel.Cover = ""
	}
	novel.From = u.String()
	doc.Find(".infotitle i").Each(func(i int, s *goquery.Selection) {
		ss := strings.Split(s.Text(), "：")
		strings.Trim(ss[1], " ")
		sss := ss[1]
		if len(ss) == 2 && i == 0 {
			novel.Auth = sss
		}
		if len(ss) == 2 && i == 1 {
			novel.Style = sss
		}
		if len(ss) == 2 && i == 2 {
			novel.Status = sss
		}
	})
	novel.Style = strings.Replace(novel.Style, "小说", "", -1)
	introString, err := doc.Find(".intro").Html()
	if err != nil {
		return nil, err
	}
	novel.Introduction = introString

	if !snwx.WithOutChapters {
		doc.Find("div#list dl dd").Each(func(i int, s *goquery.Selection) {
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
