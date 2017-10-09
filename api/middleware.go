package api

import (
	"fmt"
	"strconv"

	"net/http"

	"git.oschina.net/cnjack/novel-spider/model"
	"github.com/labstack/echo"
)

const PageOptionKey = `Page-Option-Key`

func ParseParam(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		op := &model.PageOption{
			Page:  0,
			Count: 25,
			All:   false,
			Sort:  "desc",
		}
		pageString := c.FormValue("page")
		page, _ := strconv.Atoi(pageString)
		if page != 0 {
			op.Page = page
		}
		countString := c.FormValue("count")
		count, _ := strconv.Atoi(countString)
		if count != 0 {
			op.Count = count
		}
		AllString := c.FormValue("count")
		all, err := strconv.ParseBool(AllString)
		if err == nil {
			op.All = all
		}
		SortString := c.FormValue("orderby")
		if SortString != "" {
			op.Sort = SortString
		}
		fmt.Println(op.Page)
		c.Set(PageOptionKey, op)
		return next(c)
	}
}

func ErrorHandle(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		err := next(c)
		if err == nil {
			return nil
		}
		switch err.(type) {
		case *NightcErr:
			return c.JSON(err.(*NightcErr).HttpCode, map[string]interface{}{
				"code": err.(*NightcErr).Code,
				"data": err.(*NightcErr).Data,
			})
		default:
			return c.JSON(http.StatusInternalServerError, map[string]interface{}{
				"code": 1,
				"data": err.Error(),
			})
		}
	}
}
