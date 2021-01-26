package infra

import (
	"github.com/go-redis/redis/v7"

	"github.com/Min-Feng/ratelimiter/pkg/configs"
)

func NewRedis(cfg *configs.Redis) *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:       cfg.Address(),
		Password:   cfg.Password,
		DB:         0,
		MaxRetries: 3,
	})

	_, err := client.Ping().Result()
	if err != nil {
		panic(err)
	}
	return client
}
