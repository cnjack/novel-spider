package api

import (
	_ "net/http/pprof"

	"spider/internal/config"

	binder "github.com/cnjack/echo-binder"
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

	{
		v1 := e.Group("v1")
		v1.GET("/novel/:id", GetNovelDetails)
		v1.GET("/novel/from_url/:url", GetNovelDetailsFromUrl)
		v1.DELETE("/novel/:id", DeleteNovel)
		v1.GET("/novel/:id/chapters", GetNovelChapters)

		v1.GET("/novels", GetNovels, ParseParam)
		v1.GET("/novels/style/:style", GetStyleNovels, ParseParam)

		v1.GET("/styles", GetStyles)

		v1.POST("/search/local", PostSearchLocal, ParseParam)
		v1.POST("/novel/add", AddNovel)

		v1.GET("/chapter/:url", GetChapter)
		v1.GET("/image_proxy/:url", GetImage)
	}

	port := config.GetConfig().HttpConfig.Port
	if err := e.Start(port); err != nil {
		e.Logger.Fatal(err.Error())
	}
}
