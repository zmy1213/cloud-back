package redis

import (
	"context"
	"time"

	redisv9 "github.com/redis/go-redis/v9"
)

func Ping(addr, password string, db int, timeout time.Duration) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	client := redisv9.NewClient(&redisv9.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})
	defer client.Close()

	return client.Ping(ctx).Err()
}
