package config

import (
	"time"

	"github.com/mohammaderm/rootext/repository/postgres"
	"github.com/mohammaderm/rootext/repository/redis"
	"github.com/mohammaderm/rootext/service/authService"
)

type Config struct {
	Postgres   postgres.Config    `koanf:"postgres"`
	HTTPServer HTTPServer         `koanf:"http_server"`
	Auth       authService.Config `koanf:"auth"`
	Redis      redis.Config       `koanf:"redis"`
}

type HTTPServer struct {
	Port int `koanf:"port"`
}

const (
	JwtSignKey                 = "jwt_secret"
	AccessTokenSubject         = "ac"
	RefreshTokenSubject        = "rt"
	AccessTokenExpireDuration  = time.Hour * 24
	RefreshTokenExpireDuration = time.Hour * 24 * 7
	AuthMiddlewareContextKey   = "claims"
)
