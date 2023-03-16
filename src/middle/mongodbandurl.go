package middle

import (

	// "fmt"

	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
)

// IsAdmin middleware evalutes if the user is admin - super admin
func DoingJwt(next echo.HandlerFunc) echo.HandlerFunc {
	fmt.Println("")
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file in routes")
	}
	key := os.Getenv("EncryptionKey")
	// key := "Myrachanto"

	return func(c echo.Context) error {
		//  put := c.Request().Method
		headertoken := c.Request().Header.Get("Authorization")
		token := strings.Split(headertoken, " ")[1]
		claims := jwt.MapClaims{}
		_, err := jwt.ParseWithClaims(token, claims, func(*jwt.Token) (interface{}, error) {
			return []byte(key), nil
		})
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "unable to parse token")
		}
		// fmt.Println("fc", fc)
		// level := claims["role"].(string)
		// mongourl := claims["mongourl"].(string)
		// mongodb := claims["mongodb"].(string)
		bizname := claims["bizname"].(string)
		if bizname == "" {
			return echo.NewHTTPError(http.StatusForbidden, "unable to parse token business")
		}
		// fmt.Println(">>>>>>>>>>>>>>>>>>>", bizname)

		//ensuring the context has the db variable to all routes
		// c.Set("mongourl", mongourl)
		c.Set("bizname", bizname)
		return next(c)
	}
}

// IsAdmin middleware evalutes if the user is admin - super admin
func DoingJwtTest(next echo.HandlerFunc) echo.HandlerFunc {
	fmt.Println("llllllllllllll")
	// err := godotenv.Load()
	// if err != nil {
	// 	log.Fatal("Error loading .env file in routes")
	// }
	// key := os.Getenv("EncryptionKey")
	key := "Myrachanto"
	return func(c echo.Context) error {
		headertoken := c.Request().Header.Get("Authorization")
		token := strings.Split(headertoken, " ")[1]
		claims := jwt.MapClaims{}
		_, err := jwt.ParseWithClaims(token, claims, func(*jwt.Token) (interface{}, error) {
			return []byte(key), nil
		})
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "unable to parse token")
		}
		// fmt.Println("fc", fc)
		// level := claims["role"].(string)
		// mongourl := claims["mongourl"].(string)
		// mongodb := claims["mongodb"].(string)
		bizname := claims["bizname"].(string)
		if bizname == "" {
			return echo.NewHTTPError(http.StatusForbidden, "unable to parse token business")
		}
		// fmt.Println(">>>>>>>>>>>>>>>>>>>", bizname)

		//ensuring the context has the db variable to all routes
		// c.Set("mongourl", mongourl)
		c.Set("bizname", bizname)
		return next(c)
	}
}
