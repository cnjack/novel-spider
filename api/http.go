package api

import (
	_ "net/http/pprof"

	"git.oschina.net/cnjack/novel-spider/config"
	"github.com/cnjack/echo-binder"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

func Start() {
	e := echo.New()
	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowMethods: []string{echo.GET, echo.HEAD, echo.PUT, echo.PATCH, echo.POST, echo.DELETE},
		AllowOrigins: []string{"*"},
	}))
	e.Use(middleware.Logger())
	e.Use(ErrorHandle)
	e.Use(binder.BindBinder(e))
	v1 := e.Group("v1")
	WarpRouter(v1)
	SpiderWarpRouter(v1)
	port := config.GetHttpConfig().Port
	if err := e.Start(port); err != nil {
		e.Logger.Fatal(err.Error())
	}
}
