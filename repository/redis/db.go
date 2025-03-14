package redis

import (
	"fmt"

	"github.com/redis/go-redis/v9"
)

type Config struct {
	Host     string
	Port     int
	Password string
	DB       int
}

type Redis struct {
	config Config
	client *redis.Client
}

func New(config Config) *Redis {
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", config.Host, config.Port),
		Password: config.Password,
		DB:       config.DB,
	})
	return &Redis{
		config: config,
		client: rdb,
	}
}

func (r *Redis) Conn() *redis.Client {
	return r.client
}
