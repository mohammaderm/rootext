package postHandler

import (
	"github.com/mohammaderm/rootext/presentation/httpServer/middleware"

	"github.com/labstack/echo/v4"
)

func (h Handler) SetRoutes(e *echo.Echo) {
	postGroup := e.Group("/post")

	// Operations related to user's own posts
	postGroup.POST("", h.Create, middleware.Auth(h.authSvc, h.authConfig))
	postGroup.DELETE("/:id", h.Delete, middleware.Auth(h.authSvc, h.authConfig))
	postGroup.PUT("", h.Update, middleware.Auth(h.authSvc, h.authConfig))
	postGroup.GET("", h.GetAll, middleware.Auth(h.authSvc, h.authConfig))
	postGroup.GET("/:id", h.GetById, middleware.Auth(h.authSvc, h.authConfig))

	// Operations related to all posts
	postGroup.POST("/vote", h.VotePost, middleware.Auth(h.authSvc, h.authConfig))
	postGroup.GET("/getSorted", h.GetSortedPost)

}
