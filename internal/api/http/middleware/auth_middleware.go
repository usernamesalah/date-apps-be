package middleware

import (
	"date-apps-be/infrastructure/config"
	"date-apps-be/internal/model"
	"fmt"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

func Authorized(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		authHeader := c.Request().Header.Get("Authorization")
		if authHeader == "" {
			return c.JSON(http.StatusUnauthorized, map[string]string{
				"message": "Missing Authorization header",
			})
		}

		tokenString := strings.Replace(authHeader, "Bearer ", "", -1)

		token, err := jwt.ParseWithClaims(tokenString, &model.JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
			if method, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
				return nil, fmt.Errorf("signing method invalid")
			} else if method != jwt.SigningMethodRS256 {
				return nil, fmt.Errorf("signing method invalid")
			}
			return config.Get().JWTRS256PubKey, nil
		})

		if err != nil {
			return echo.NewHTTPError(http.StatusUnauthorized, err)
		}

		claims, ok := token.Claims.(*model.JWTClaims)
		if !ok || !token.Valid {
			return echo.NewHTTPError(http.StatusUnauthorized, err)
		}

		c.Set("userInfo", claims)
		c.Set("authHeader", authHeader)
		next(c)

		return nil
	}
}
