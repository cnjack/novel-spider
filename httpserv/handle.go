package httpserv

import (
	"net/http"

	"github.com/labstack/echo"
)

func WarpRouter(g *echo.Group) {
	g.GET("/status", getStatus)
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

func getStatus(c echo.Context) error {
	return c.JSON(http.StatusOK, struct {
		Code int    `json:"code"`
		Data Status `json:"data"`
	}{
		Code: 0,
		Data: status,
	})
}
