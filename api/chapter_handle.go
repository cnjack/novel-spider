package api

import (
	"net/http"
	"net/url"

	"github.com/labstack/echo"
	"spider/model"
	"spider/spider"
)

func GetChapter(c echo.Context) error {
	urlString := c.Param("url")
	urlString, err := url.QueryUnescape(urlString)
	if err != nil {
		return ParamError
	}
	if _, err := url.Parse(urlString); err != nil {
		return ParamError
	}

	chapter, err := model.GetChapter(urlString)
	if err != nil {
		return ServerError
	}
	if err != nil {
		return ServerError
	}
	return c.JSON(http.StatusOK, struct {
		Code int         `json:"code"`
		Data interface{} `json:"data"`
	}{
		Code: 0,
		Data: chapter,
	})
}

var imageSpider = spider.NewImageSpider(model.MustGetRedisClient())

func GetImage(c echo.Context) error {
	urlString := c.Param("url")
	urlString, err := url.QueryUnescape(urlString)
	if err != nil {
		return ParamError
	}
	if _, err := url.Parse(urlString); err != nil {
		return ParamError
	}
	c.Response().Header().Add("Cache-Control", "max-age=3600")
	if err := imageSpider.WriteWithUrl(urlString, c.Response()); err != nil {
		c.Response().Header().Add("Content-Type", "image/gif")
		c.Response().Write(imageSpider.DefaultImage())
	}
	return nil
}
