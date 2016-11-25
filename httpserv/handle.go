package httpserv

import (
	"net/http"

	"github.com/labstack/echo"
)

func WarpRouter(g *echo.Group) {
	g.GET("/", indexHandle)
	g.GET("/status", getStatus)
	g.GET("/statuss", getStatuss)
	g.GET("/novel/:id", getNovelDetails)
	g.DELETE("/novel/:id", deleteNovel)
	g.GET("/novels", getNovels, ParseParam)
	g.GET("/novels/style/:style", getStyleNovels, ParseParam)
	g.GET("/styles", getStyles)
	g.POST("/search/remote", postSearchRemote)
	g.POST("/search/local", postSearchLocal, ParseParam)
	g.POST("/novel/task", postNovelTask)
	g.GET("/novel/:id/chapters", getNovelChapters)
	g.GET("/chapter/:id", getChapter)
}

func indexHandle(c echo.Context) error {
	return c.String(http.StatusOK, "")
}

func getStatus(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]interface{}{
		"code": 0,
		"data": status,
	})
}

func getStatuss(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]interface{}{
		"code": 0,
		"data": statuss,
	})
}
