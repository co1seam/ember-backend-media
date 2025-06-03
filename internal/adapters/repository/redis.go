package repository

import (
	"github.com/go-redis/redis/v8"
)

type Redis struct {
	Redis *redis.Client
}

func NewRedis(host, port string) *Redis {
	client := redis.NewClient(&redis.Options{
		Addr:     host + ":" + port,
		Password: "",
		DB:       0,
	})
	return &Redis{Redis: client}
}
