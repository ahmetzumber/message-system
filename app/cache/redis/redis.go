package redis

import (
	"context"
	"message-system/config"
	"time"

	"github.com/redis/go-redis/v9"
)

type Cache struct {
	Client *redis.Client
	URI    string
}

func NewCache(config *config.RedisConfig) *Cache {
	rdb := redis.NewClient(&redis.Options{
		Addr: config.URI,
	})

	return &Cache{
		Client: rdb,
		URI:    config.URI,
	}
}

func (c *Cache) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	status := c.Client.Set(ctx, key, value, expiration)
	return status.Err()
}
