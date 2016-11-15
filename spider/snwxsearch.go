package spider

import (
	"net/url"

	"fmt"

	"git.oschina.net/cnjack/downloader"
	"github.com/PuerkitoBio/goquery"
)

type SnwxSearch struct {
	NovelName string
	Data      interface{}
}

func (s *SnwxSearch) Name() string {
	return "snwx.com"
}

func (s *SnwxSearch) Match(name string) bool {
	s.NovelName = name
	return true
}

func (s *SnwxSearch) Gain() (interface{}, error) {
	u, _ := url.Parse(fmt.Sprintf("http://zhannei.baidu.com/cse/search?q=%s&click=1&s=5516249222499057291&nsid=", s.NovelName))
	d := downloader.NewHttpDownloaderFromUrl(u).Download()
	if err := d.Error(); err != nil {
		return "", err
	}
	doc, err := d.Resource().Document()
	if err != nil {
		return "", err
	}
	var searchs = []*Search{}
	doc.Find(".result").Each(func(i int, selection *goquery.Selection) {
		var b bool
		d := selection.Find(".c-title a")
		search := Search{}
		from, b := d.Attr("href")

		if !b {
			return
		}
		filter := &SnwxNovel{}
		if !filter.Match(from) {
			return
		}
		search.SearchName = d.Text()
		search.From = filter.Url.String()
		search.Name = s.NovelName
		searchs = append(searchs, &search)
	})
	return searchs, nil
}
