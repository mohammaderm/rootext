package middleware

import (
	"rootext/config"
	"rootext/service/authService"

	mw "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
)

func Auth(service authService.AuthService, cfg authService.Config) echo.MiddlewareFunc {
	return mw.WithConfig(mw.Config{

		ContextKey: config.AuthMiddlewareContextKey,
		SigningKey: []byte(cfg.SignKey),

		SigningMethod: "HS256",
		ParseTokenFunc: func(c echo.Context, auth string) (interface{}, error) {
			claims, err := service.ParseToken(auth)
			if err != nil {
				return nil, err
			}
			return claims, nil
		},
	})
}
