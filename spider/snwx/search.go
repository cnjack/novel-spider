package snwx

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"net/url"
	"spider/downloader"
	"spider/spider"
	"sync"

	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"

	"github.com/PuerkitoBio/goquery"
)

type Search struct {
	NovelName string
	Data      interface{}
}

func (s *Search) Name() string {
	return "snwx.com"
}

func (s *Search) Match(name string) bool {
	s.NovelName = name
	return true
}

func (s *Search) Gain() (interface{}, error) {
	postNovelName, _ := ioutil.ReadAll(transform.NewReader(bytes.NewReader([]byte(s.NovelName)), simplifiedchinese.GBK.NewEncoder()))
	u, _ := url.Parse(fmt.Sprintf("https://www.snwx8.com/modules/article/search.php?searchkey=%s", url.QueryEscape(string(postNovelName))))
	d := downloader.NewHttpDownloaderFromUrl(u).Download()
	if err := d.Error(); err != nil {
		return "", err
	}
	doc, err := d.Resource().Document()
	if err != nil {
		return "", err
	}
	var searchs = []*spider.Search{}
	var wg = sync.WaitGroup{}
	doc.Find("#newscontent ul li").Each(func(i int, selection *goquery.Selection) {
		var b bool
		d := selection.Find(".s2 a")
		search := spider.Search{}
		from, b := d.Attr("href")
		if !b {
			return
		}
		filter := &Novel{WithOutChapters: true}
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
			novel := n.(spider.Novel)
			search.Novel = &novel
		}()
		searchs = append(searchs, &search)
	})
	wg.Wait()
	if searchs == nil || len(searchs) == 0 {
		return nil, errors.New("nothing search")
	}
	out := make([]*spider.Search, 0)
	mark := make(map[string]bool)
	for _, v := range searchs {
		if _, ok := mark[v.Novel.From]; !ok {
			mark[v.Novel.From] = true
			out = append(out, v)
		}
	}
	return out, nil
}
