package api

import (
	"net/http"
	"net/url"

	"git.oschina.net/cnjack/novel-spider/model"
	"github.com/labstack/echo"
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
