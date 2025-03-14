package postHandler

import (
	"github.com/mohammaderm/rootext/service/authService"
	"github.com/mohammaderm/rootext/service/postService"
)

type Handler struct {
	postSvc    postService.PostService
	authSvc    authService.AuthService
	authConfig authService.Config
}

func New(postSvc postService.PostService, authConfig authService.Config, authSvc authService.AuthService) Handler {
	return Handler{
		postSvc:    postSvc,
		authSvc:    authSvc,
		authConfig: authConfig,
	}
}
