package routes

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

func BindInstrumentationHandlers(e *echo.Group) error {
	e.GET("health", HealthCheck)
	return nil
}

func HealthCheck(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]interface{}{
		"data": "Server is up and running",
	})
}
