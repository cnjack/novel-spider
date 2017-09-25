package spider

import (
	"errors"
	"fmt"
	"net/url"
	"sync"

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
	var wg = sync.WaitGroup{}
	doc.Find(".result").Each(func(i int, selection *goquery.Selection) {
		var b bool
		d := selection.Find(".c-title a")
		search := Search{}
		from, b := d.Attr("href")
		if !b {
			return
		}
		filter := &SnwxNovel{WithOutChapters: true}
		if !filter.Match(from) {
			return
		}
		search.SearchName = d.Text()
		search.From = from
		search.Name = s.NovelName

		wg.Add(1)
		go func() {
			defer func() {
				wg.Done()
			}()
			n, err := filter.Gain()
			if err != nil {
				return
			}
			novel := n.(Novel)
			search.Novel = &novel
		}()
		searchs = append(searchs, &search)
	})
	wg.Wait()
	if searchs == nil || len(searchs) == 0 {
		return nil, errors.New("nothing search")
	}
	out := make([]*Search, 0)
	mark := make(map[string]bool)
	for _, v := range searchs {
		if _, ok := mark[v.Novel.From]; !ok {
			mark[v.Novel.From] = true
			out = append(out, v)
		}
	}
	return out, nil
}
