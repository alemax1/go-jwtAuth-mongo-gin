package cache

import (
	"context"

	"github.com/go-redis/redis/v8"
	"github.com/spf13/viper"
)

type cache struct {
	Client *redis.Client
}

var Redis cache

func CreateRedisConnection() error {
	Redis.Client = redis.NewClient(&redis.Options{
		Addr:     viper.GetString("redis.url"),
		Password: viper.GetString("redis.password"),
	})

	if _, err := Redis.Client.Ping(context.Background()).Result(); err != nil {
		return err
	}

	return nil
}
