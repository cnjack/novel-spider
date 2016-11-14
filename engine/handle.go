package engine

import (
	"net/http"

	"github.com/labstack/echo"
)

func IndexHandle(c echo.Context) error {
	return c.String(http.StatusOK, "")
}
