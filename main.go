package main

import (
	"os"
	"os/signal"
	"syscall"

	postrepositoryredis "github.com/mohammaderm/rootext/repository/redis/postRepositoryRedis"
	"github.com/mohammaderm/rootext/service/authService"
	"github.com/mohammaderm/rootext/service/postService"
	"github.com/mohammaderm/rootext/service/userService"

	"github.com/mohammaderm/rootext/config"
	httpserver "github.com/mohammaderm/rootext/presentation/httpServer"
	"github.com/mohammaderm/rootext/presentation/httpServer/postHandler"
	"github.com/mohammaderm/rootext/presentation/httpServer/userHandler"
	"github.com/mohammaderm/rootext/repository/migrator"
	"github.com/mohammaderm/rootext/repository/postgres"
	postrepository "github.com/mohammaderm/rootext/repository/postgres/postRepository"
	"github.com/mohammaderm/rootext/repository/postgres/userRepository"
	"github.com/mohammaderm/rootext/repository/redis"
)

func main() {
	cfg := config.Load("config.yml")

	// db

	migrator := migrator.New(cfg.Postgres)
	// migrator.Down()
	migrator.Up()

	postgresDB := postgres.New(cfg.Postgres)
	redisDB := redis.New(cfg.Redis)

	// repository
	postCache := postrepositoryredis.New(redisDB)
	userRepo := userRepository.New(postgresDB)
	postRepo := postrepository.New(postgresDB, postCache)

	// service
	authSvc := authService.New(cfg.Auth)
	userSvc := userService.New(userRepo, authSvc)
	postSvc := postService.New(postRepo, postCache)

	// handler
	userHandler := userHandler.New(userSvc, cfg.Auth, authSvc)
	postHnadler := postHandler.New(postSvc, cfg.Auth, authSvc)

	server := httpserver.New(cfg, userHandler, postHnadler)

	go func() {
		server.Serve()
	}()

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)
	<-shutdown
}
