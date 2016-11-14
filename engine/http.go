package engine

import (
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

	{
		e.Static("/", ".")
		e.GET("/", IndexHandle)
	}

	if err := e.Start(":1314"); err != nil {
		e.Logger.Fatal(err.Error())
	}
}
