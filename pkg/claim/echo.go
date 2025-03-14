package claim

import (
	"github.com/mohammaderm/rootext/config"
	"github.com/mohammaderm/rootext/service/authService"

	"github.com/labstack/echo/v4"
)

func GetClaimsFromEchoContext(c echo.Context) *authService.Claims {
	return c.Get(config.AuthMiddlewareContextKey).(*authService.Claims)
}
