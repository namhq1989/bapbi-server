package caching

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

type Caching struct {
	redis *redis.Client
}

func NewCachingClient(redisURL string) *Caching {
	var (
		ctx    = context.Background()
		opt, _ = redis.ParseURL(redisURL)
	)

	// new client
	client := redis.NewClient(opt)

	// ping
	if _, err := client.Ping(ctx).Result(); err != nil {
		panic(err)
	}

	fmt.Printf("⚡️ [caching]: connected \n")

	return &Caching{redis: client}
}
