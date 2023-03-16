package middle

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
	httperors "github.com/myrachanto/erroring"
	"github.com/spf13/viper"
)

type Userkey struct {
	EncryptionKey string `mapstructure:"EncryptionKey"`
}

func Loaduserkey() (userkey Userkey, err error) {
	viper.AddConfigPath(".")
	viper.SetConfigName("app")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}
	err = viper.Unmarshal(&userkey)
	return
}

// I
var next func(c echo.Context) error

func IsAuthorized() echo.MiddlewareFunc {
	// var next echo.HandlerFunc
	return Receipts_Read(next)
}
func Receipts_Read(next echo.HandlerFunc) echo.HandlerFunc {
	userkey, err := Loaduserkey()
	if err != nil {
		log.Fatal("cannot load config:", err)
	}
	return func(c echo.Context) error {
		claims, err := Claimer(userkey, c)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "unable to parse token")
		}
		// fmt.Println("fc", fc)
		authmodules := claims["authmodules"].([]interface{})
		authings := Reading(authmodules)
		reader := Moduling("receipts", authings)
		if !reader {
			return echo.NewHTTPError(http.StatusForbidden, "You are not authorized to be here!!!")
		}
		return next(c)
	}
}

func Reading(authmodules []interface{}) []AuthModule {
	authing := AuthModule{}
	authings := []AuthModule{}
	for _, val := range authmodules {
		auth := val.(map[string]interface{})
		authing.ModuleName = fmt.Sprintf("%s", auth["module_name"])
		authing.Read = auth["read"].(bool)
		authing.Write = auth["write"].(bool)
		authings = append(authings, authing)

	}
	return authings
}

func CheckAuth(module string, authings []AuthModule) bool {
	var reader bool
	for _, l := range authings {
		if l.ModuleName == module && l.Read {
			reader = l.Read
		}
		if l.ModuleName == module && l.Write {
			reader = l.Read
		}
	}
	return reader
}

func Claimer(userkey Userkey, c echo.Context) (jwt.MapClaims, httperors.HttpErr) {
	headertoken := c.Request().Header.Get("Authorization")
	token := strings.Split(headertoken, " ")[1]
	claims := jwt.MapClaims{}
	_, err := jwt.ParseWithClaims(token, claims, func(*jwt.Token) (interface{}, error) {
		return []byte(userkey.EncryptionKey), nil
	})
	if err != nil {
		return nil, httperors.NewNotFoundError("unable to parse token")
	}
	return claims, nil
}

type AuthModule struct {
	ModuleName string
	Read       bool
	Write      bool
}

func Moduling(module string, authings []AuthModule) bool {
	var reader bool
	// auths := model.AuthModule{}
	for _, l := range authings {
		if l.ModuleName == module {
			reader = l.Read
			// auths.ModuleName = l.ModuleName
			// auths.Write = l.Write
			// auths.Read = l.Read
		}

	}
	return reader
}
