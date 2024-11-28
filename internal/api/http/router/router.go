package router

import (
	"date-apps-be/internal/container"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func Init(e *echo.Echo, hc *container.HandlerComponent, sc *container.SharedComponent) {
	e.Pre(middleware.Rewrite(map[string]string{
		"/v1/*": "/$1",
	}))

	e.GET("/ping", ping)
	publicRouter(e, hc)

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     hc.Config.CORSOrigins,
		AllowMethods:     []string{"GET", "POST", "DELETE", "PUT", "PATCH", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Accept", "Content-Type", "Authorization", "X-User-ID", "x-token", "x-mkt-token", "X-Requested-With", "x-device-id", "x-service-authorization", "x-service-timestamp"},
		AllowCredentials: true,
	}))
}

// ping write pong to http.ResponseWriter.
func ping(c echo.Context) error {
	return c.String(http.StatusOK, "pong")
}
