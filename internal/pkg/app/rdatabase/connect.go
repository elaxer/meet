package rdatabase

import (
	"context"
	"fmt"
	"meet/internal/config"

	"github.com/redis/go-redis/v9"
)

func Connect(cfg *config.RedisConfig) (*redis.Client, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		Password: cfg.Password,
		DB:       cfg.Databases,
	})
	if err := rdb.Ping(context.Background()).Err(); err != nil {
		return nil, err
	}

	return rdb, nil
}

func MustConnect(cfg *config.RedisConfig) *redis.Client {
	rdb, err := Connect(cfg)
	if err != nil {
		panic(err)
	}

	return rdb
}
