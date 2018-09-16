package kxs

import (
	"bytes"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"spider/spider"
	"sync"

	"golang.org/x/text/transform"

	"spider/downloader"

	"github.com/PuerkitoBio/goquery"
	"golang.org/x/text/encoding/simplifiedchinese"
)

type Search struct {
	NovelName string
	Data      interface{}
}

func (s *Search) Name() string {
	return "00kxs.com"
}

func (s *Search) Match(name string) bool {
	s.NovelName = name
	return true
}

func (s *Search) Gain() (interface{}, error) {
	postParams := url.Values{}
	postNovelName, _ := ioutil.ReadAll(transform.NewReader(bytes.NewReader([]byte(s.NovelName)), simplifiedchinese.GBK.NewEncoder()))
	postParams.Set("searchkey", string(postNovelName))
	postParams.Set("searchtype", "articlename")
	req, _ := http.NewRequest(http.MethodPost, "http://www.00kxs.com/modules/article/search.php", bytes.NewBufferString(postParams.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	d := downloader.NewHttpDownloaderFromRequest(req).Download()
	if err := d.Error(); err != nil {
		return "", err
	}
	doc, err := d.Resource().Document()
	if err != nil {
		return "", err
	}
	var searchs = []*spider.Search{}
	var wg = sync.WaitGroup{}
	doc.Find("#nr").Each(func(i int, selection *goquery.Selection) {
		var b bool
		d := selection.Find(".odd a")
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
