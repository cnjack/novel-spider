package engine

import (
	"strconv"

	"git.oschina.net/cnjack/novel-spider/model"
	"github.com/labstack/echo"
	"fmt"
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
