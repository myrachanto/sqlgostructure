package routes

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/myrachanto/sqlgostructure/src/api/controller"
	"github.com/myrachanto/sqlgostructure/src/api/repository"
	"github.com/myrachanto/sqlgostructure/src/api/service"
	"github.com/myrachanto/sqlgostructure/src/middle"
)

// var passer echo.MiddlewareFunc

func ApiServer() {
	u := controller.NewUserController(service.NewUserService(repository.NewUserRepo()))
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())
	e.HTTPErrorHandler = HttpErrorHandler

	api := e.Group("/api")
	// echo.MiddlewareFunc()
	// api.Use(middle.IsAuthorized())

	{
		// api.POST("/register", u.Create, middle.IsAuthorized())
		api.POST("/register", u.Create)
		api.POST("/login", u.Login)
		api.GET("/users", u.GetAll, middle.PasetoAuthMiddleware)
		api.GET("/users/:code", u.GetOne, middle.PasetoAuthMiddleware)
		api.GET("/users/logout", u.Logout, middle.PasetoAuthMiddleware)
		api.POST("/users/renew", u.RenewAccessToken)
		api.PUT("/users/passwordUpdate", u.PasswordUpdate)
		api.PUT("/users/passwordReset", u.PasswordReset)
	}
	e.GET("/healthy", HealthCheck)
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file in routes")
	}
	PORT := os.Getenv("PORT")
	e.Logger.Fatal(e.Start(PORT))
}
