package router

import (
	"date-apps-be/internal/api/http/handler"
	"date-apps-be/internal/api/http/middleware"
	"date-apps-be/internal/container"

	"github.com/labstack/echo/v4"
)

func publicRouter(e *echo.Echo, hc *container.HandlerComponent) {

	// User
	userHandler := handler.NewUserHandler(hc)
	userRoute := e.Group("/users")
	{
		userRoute.Use(middleware.Authorized)
		userRoute.GET("/profile", userHandler.GetUserProfile)
	}

	e.POST("/login", userHandler.Login)
	e.POST("/register", userHandler.Register)
}
