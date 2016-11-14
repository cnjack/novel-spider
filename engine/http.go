package engine

import (
	"git.oschina.net/cnjack/novel-spider/config"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

func Http() {
	e := echo.New()

	e.Use(middleware.Recover())
	e.Use(middleware.CORS())
	e.Use(middleware.Gzip())
	e.Use(middleware.BodyLimit("1M"))
	e.Use(middleware.Logger())

	e.Static("/", config.GetHttpConfig().StaticPath)
	e.GET("/", IndexHandle)

	port := config.GetHttpConfig().Port
	if err := e.Start(port); err != nil {
		e.Logger.Fatal(err.Error())
	}
}
