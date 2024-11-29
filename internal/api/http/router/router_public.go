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
	userMatchHandler := handler.NewUserMatchHandler(hc)
	premiumConfigHandler := handler.NewPremiumConfigHandler(hc)

	//route
	e.POST("/login", userHandler.Login)
	e.POST("/register", userHandler.Register)

	userRoute := e.Group("/users")
	{
		userRoute.Use(middleware.Authorized)
		userRoute.GET("/profile", userHandler.GetUserProfile)
		userRoute.GET("/package", userHandler.GetMyPackage)
	}

	userMatchRoute := e.Group("/matches")
	{
		userMatchRoute.Use(middleware.Authorized)
		userMatchRoute.POST("", userMatchHandler.CreateMatch)
		userMatchRoute.GET("", userMatchHandler.GetUserMatches)
	}

	premiumConfigRoute := e.Group("/packages")
	{
		premiumConfigRoute.GET("", premiumConfigHandler.GetPackages)
		premiumConfigRoute.GET("/:uid", premiumConfigHandler.GetPackageByUID)
		premiumConfigRoute.POST("/purchase", premiumConfigHandler.PurchasePackage, middleware.Authorized)
	}

}
