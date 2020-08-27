package api

import (
	"net/http"
	"spider/internal/repository"

	"github.com/labstack/echo"
)

type PostSearchParam struct {
	Title string `json:"title"`
}

func PostSearchLocal(c echo.Context) error {
	postSearchParam := &PostSearchParam{}
	err := c.Bind(postSearchParam)
	if err != nil {
		return ParamError
	}
	op := c.Get(PageOptionKey).(*repository.PageOption)
	if postSearchParam.Title == "" {
		return ParamError
	}
	db := repository.MustGetDB()
	novels, err := repository.SearchByTitleOrAuth(db, postSearchParam.Title, postSearchParam.Title, op)
	if err != nil {
		return ServerError
	}
	nextPage := 0
	if len(novels) >= op.Count {
		nextPage = op.Page + 1
	}
	return c.JSON(http.StatusOK, struct {
		Code int                       `json:"code"`
		Data []*repository.SearchNovel `json:"data"`
		Next int                       `json:"next"`
	}{
		Code: 0,
		Data: novels,
		Next: nextPage,
	})
}
