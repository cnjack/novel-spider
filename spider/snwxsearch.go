package spider

import (
	"errors"
	"fmt"

	"github.com/PuerkitoBio/goquery"
	"github.com/hu17889/go_spider/core/common/request"
)

var filter = SnwxNovel{}

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
	url := fmt.Sprintf("http://zhannei.baidu.com/cse/search?q=%s&click=1&s=5516249222499057291&nsid=", s.NovelName)
	page := d.Download(request.NewRequest(url, "html", "", "GET", "", nil, nil, nil, nil))
	if page.Errormsg() != "" {
		return "", errors.New(page.Errormsg())
	}
	doc := page.GetHtmlParser()
	var searchs = []*Search{}
	doc.Find(".result").Each(func(i int, selection *goquery.Selection) {
		var b bool
		d := selection.Find(".c-title a")
		search := Search{}
		search.From, b = d.Attr("href")
		if !b {
			return
		}
		if !filter.Match(search.From) {
			return
		}
		search.SearchName = d.Text()
		search.Name = s.NovelName
		searchs = append(searchs, &search)
	})
	return searchs, nil
}
