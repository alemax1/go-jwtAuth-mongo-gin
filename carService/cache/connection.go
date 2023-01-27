package cache

import (
	"context"

	"github.com/go-redis/redis/v8"
	"github.com/spf13/viper"
)

type Redis struct {
	client *redis.Client
}

var cache Redis

func CreateRedisConnection() error {
	cache.client = redis.NewClient(&redis.Options{
		Addr:     viper.GetString("redis.url"),
		Password: viper.GetString("redis.password"),
	})

	if _, err := cache.client.Ping(context.Background()).Result(); err != nil {
		return err
	}

	return nil
}
