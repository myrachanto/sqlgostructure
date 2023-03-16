package routes

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

// HealthCheck godoc
// @Summary Show the Health status of server.
// @Description get the Health status of server.
// @Tags Health Status
// @Accept */*
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /health [get]
func HealthCheck(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]interface{}{
		"data": "server is up and running",
	})
}
func HttpErrorHandler(err error, c echo.Context) {
	fmt.Println(err)
}
