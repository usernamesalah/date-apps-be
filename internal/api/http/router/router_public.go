package router

import (
	"date-apps-be/internal/api/http/handler"
	"date-apps-be/internal/api/http/middleware"
	"date-apps-be/internal/container"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	echoMiddleware "github.com/labstack/echo/v4/middleware"
	"golang.org/x/time/rate"
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
		config := echoMiddleware.RateLimiterConfig{
			Skipper: echoMiddleware.DefaultSkipper,
			Store: echoMiddleware.NewRateLimiterMemoryStoreWithConfig(
				echoMiddleware.RateLimiterMemoryStoreConfig{Rate: rate.Limit(10), Burst: 10, ExpiresIn: 1 * time.Second},
			),
			IdentifierExtractor: func(ctx echo.Context) (string, error) {
				id := ctx.RealIP()
				return id, nil
			},
			ErrorHandler: func(context echo.Context, err error) error {
				return context.JSON(http.StatusForbidden, nil)
			},
			DenyHandler: func(context echo.Context, identifier string, err error) error {
				return context.JSON(http.StatusTooManyRequests, nil)
			},
		}

		userMatchRoute.Use(echoMiddleware.RateLimiterWithConfig(config))
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
