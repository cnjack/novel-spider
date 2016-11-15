package engine

import (
	"strconv"

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
		pageString := c.Param("page")
		page, _ := strconv.Atoi(pageString)
		if page != 0 {
			op.Page = page
		}
		countString := c.Param("count")
		count, _ := strconv.Atoi(countString)
		if count != 0 {
			op.Count = count
		}
		AllString := c.Param("count")
		all, err := strconv.ParseBool(AllString)
		if err == nil {
			op.All = all
		}
		SortString := c.Param("orderby")
		if SortString != "" {
			op.Sort = SortString
		}
		c.Set(PageOptionKey, op)
		return next(c)
	}
}
