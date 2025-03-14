package userHandler

import (
	"github.com/mohammaderm/rootext/service/authService"
	"github.com/mohammaderm/rootext/service/userService"
)

type Handler struct {
	userSvc    userService.Service
	authSvc    authService.AuthService
	authConfig authService.Config
}

func New(userSvc userService.Service, authConfig authService.Config, authSvc authService.AuthService) Handler {
	return Handler{
		userSvc:    userSvc,
		authSvc:    authSvc,
		authConfig: authConfig,
	}
}
