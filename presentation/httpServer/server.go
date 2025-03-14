package httpserver

import (
	"fmt"
	"net/http"

	"github.com/mohammaderm/rootext/config"
	"github.com/mohammaderm/rootext/logger"
	"github.com/mohammaderm/rootext/presentation/httpServer/postHandler"
	"github.com/mohammaderm/rootext/presentation/httpServer/userHandler"

	_ "github.com/mohammaderm/rootext/docs"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"
	"go.uber.org/zap"
)

type Server struct {
	Router      *echo.Echo
	config      config.Config
	userHandler userHandler.Handler
	postHandler postHandler.Handler
}

// @title Rootext API
// @version 1.0
// @description This is a sample server for Rootext.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization

func New(config config.Config, userHandler userHandler.Handler, postHandler postHandler.Handler) Server {
	return Server{
		Router:      echo.New(),
		config:      config,
		userHandler: userHandler,
		postHandler: postHandler,
	}
}

func (s Server) Serve() {
	s.Router.Use(middleware.RequestID())

	s.Router.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogURI:           true,
		LogStatus:        true,
		LogHost:          true,
		LogRemoteIP:      true,
		LogRequestID:     true,
		LogMethod:        true,
		LogContentLength: true,
		LogResponseSize:  true,
		LogLatency:       true,
		LogError:         true,
		LogProtocol:      true,
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			errMsg := ""
			if v.Error != nil {
				errMsg = v.Error.Error()
			}

			logger.Logger.Named("http-server").Info("request",
				zap.String("request_id", v.RequestID),
				zap.String("host", v.Host),
				zap.String("content-length", v.ContentLength),
				zap.String("protocol", v.Protocol),
				zap.String("method", v.Method),
				zap.Duration("latency", v.Latency),
				zap.String("error", errMsg),
				zap.String("remote_ip", v.RemoteIP),
				zap.Int64("response_size", v.ResponseSize),
				zap.String("uri", v.URI),
				zap.Int("status", v.Status),
			)

			return nil
		},
	}))

	s.Router.Use(middleware.Recover())

	s.Router.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
	}))

	s.Router.GET("/health-check", s.healthCheck)
	s.Router.GET("/swagger/*", echoSwagger.WrapHandler)

	s.userHandler.SetRoutes(s.Router)
	s.postHandler.SetRoutes(s.Router)

	address := fmt.Sprintf(":%d", s.config.HTTPServer.Port)
	fmt.Printf("start echo server on %s\n", address)

	if err := s.Router.Start(address); err != nil {
		fmt.Println("router start error", err)
	}
}

func (s Server) healthCheck(c echo.Context) error {
	return c.JSON(http.StatusOK, echo.Map{"message": "good"})
}
