package middle

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// IsAdmin middleware evalutes if the user is admin - super admin
func GetAuthorizedApi(next echo.HandlerFunc) echo.HandlerFunc {

	return func(c echo.Context) error {
		bizname := c.Request().Header.Get("bizname")
		if bizname == "" {
			return echo.NewHTTPError(http.StatusForbidden, "unable to get business")
		}
		c.Set("bizname", bizname)
		// fmt.Println(">>>>>>>>>>>>>>>", c.FormValue("name"))
		return next(c)
	}
}
