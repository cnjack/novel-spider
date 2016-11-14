package engine

import (
	"git.oschina.net/cnjack/novel-spider/config"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"time"
)

var (
	startTime time.Time = time.Now()
)

func Http() {
	e := echo.New()

	e.Use(middleware.Recover())
	e.Use(middleware.CORS())
	e.Use(middleware.Gzip())
	e.Use(middleware.BodyLimit("1M"))
	e.Use(middleware.Logger())

	e.Static("/", config.GetHttpConfig().StaticPath)
	v1 := e.Group("v1")
	WarpRouter(v1)

	port := config.GetHttpConfig().Port
	if err := e.Start(port); err != nil {
		e.Logger.Fatal(err.Error())
	}
}

func init() {
	ReloadStatus()
}
