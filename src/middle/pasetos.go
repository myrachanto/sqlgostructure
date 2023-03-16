package middle

import (
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/myrachanto/sqlgostructure/src/pasetos"
)

const (
	authorisationHeaderKey = "Authorization"
	authorisationType      = "Bearer"
)

func PasetoAuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		authorizationHeader := c.Request().Header.Get(authorisationHeaderKey)
		if len(authorizationHeader) == 0 {
			return echo.NewHTTPError(http.StatusUnauthorized, "Authorization header not provided")
		}
		fields := strings.Fields(authorizationHeader)
		if len(fields) < 2 {
			return echo.NewHTTPError(http.StatusUnauthorized, "Invalid Authorization format provided")
		}
		authtype := fields[0]
		if authtype != authorisationType {
			return echo.NewHTTPError(http.StatusUnauthorized, "That type of Authorization is not allowed here!")
		}
		accessToken := fields[1]
		maker, _ := pasetos.NewPasetoMaker()
		_, err := maker.VerifyToken(accessToken)
		if err != nil {
			return echo.NewHTTPError(http.StatusUnauthorized, "That token is invalid!")
		}
		// if payload.Bizname == "" {
		// 	return echo.NewHTTPError(http.StatusUnauthorized, "Invalid shop selection!")
		// }
		// fmt.Println(" wwwhhhhhhh", c.FormValue("name"))
		// c.Set("bizname", payload.Bizname)
		return next(c)
	}
}

func PasetoIsAdmin(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		authorizationHeader := c.Request().Header.Get(authorisationHeaderKey)
		if len(authorizationHeader) == 0 {
			return echo.NewHTTPError(http.StatusUnauthorized, "Authorization header not provided")
		}
		fields := strings.Fields(authorizationHeader)
		if len(fields) < 2 {
			return echo.NewHTTPError(http.StatusUnauthorized, "Invalid Authorization format provided")
		}
		authtype := fields[0]
		if authtype != authorisationType {
			return echo.NewHTTPError(http.StatusUnauthorized, "That type of Authorization is not allowed here!")
		}
		accessToken := fields[1]
		maker, _ := pasetos.NewPasetoMaker()
		payload, err := maker.VerifyToken(accessToken)
		if err != nil {
			return echo.NewHTTPError(http.StatusUnauthorized, "That token is invalid!")
		}
		if payload.Admin != "admin" {
			return echo.NewHTTPError(http.StatusUnauthorized, "You are not authorized")
		}
		return next(c)
	}
}
