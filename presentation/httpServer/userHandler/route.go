package userHandler

import (
	"github.com/labstack/echo/v4"
)

func (h Handler) SetRoutes(e *echo.Echo) {
	userGroup := e.Group("/user")

	userGroup.POST("/auth/register", h.Register)
	userGroup.POST("/auth/login", h.Login)
	userGroup.POST("/auth/tokenRenew", h.TokenReNew)
}
