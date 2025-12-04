package resource

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

func NewRedis(cfg *Config) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr: cfg.RedisAddress,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to redis: %w", err)
	}

	return client, nil
}
