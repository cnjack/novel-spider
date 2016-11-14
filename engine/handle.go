package engine

import (
	"github.com/labstack/echo"
	"net/http"
)

func IndexHandle(c echo.Context) error {
	return c.String(http.StatusOK, "")
}
