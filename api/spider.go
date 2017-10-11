package api

import (
	"net/http"
	"net/url"

	"git.oschina.net/cnjack/novel-spider/spider"
	"git.oschina.net/cnjack/novel-spider/spider/kxs"
	"git.oschina.net/cnjack/novel-spider/spider/snwx"
	"github.com/labstack/echo"
)

func GetSpiderNovelChapters(c echo.Context) error {
	urlString := c.Param("url")
	urlString, err := url.QueryUnescape(urlString)
	if err != nil {
		return ParamError
	}
	if _, err := url.Parse(urlString); err != nil {
		return ParamError
	}
	var s spider.Spider
	chaptersSpider := []spider.Spider{
		&snwx.Chapter{},
		&kxs.Chapter{},
	}
	for _, v := range chaptersSpider {
		if v.Match(urlString) {
			s = v
		}
	}
	if s == nil {
		return ParamError
	}
	data, err := s.Gain()
	if err != nil {
		return NewNightcErr(http.StatusNotFound, 1, err.Error())
	}
	return c.JSON(http.StatusOK, data)
}

func GetSpiderNovel(c echo.Context) error {
	urlString := c.Param("url")
	urlString, err := url.QueryUnescape(urlString)
	if err != nil {
		return ParamError
	}
	if _, err := url.Parse(urlString); err != nil {
		return ParamError
	}
	var s spider.Spider
	chaptersSpider := []spider.Spider{
		&snwx.Novel{},
		&kxs.Novel{},
	}
	for _, v := range chaptersSpider {
		if v.Match(urlString) {
			s = v
		}
	}
	if s == nil {
		return ParamError
	}
	data, err := s.Gain()
	if err != nil {
		return NewNightcErr(http.StatusNotFound, 1, err.Error())
	}
	return c.JSON(http.StatusOK, data)
}
