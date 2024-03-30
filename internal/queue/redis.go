package queue

import (
	"github.com/hibiken/asynq"
	"github.com/redis/go-redis/v9"
)

func getRedisConnFromURL(redisURL string) asynq.RedisClientOpt {
	opt, _ := redis.ParseURL(redisURL)

	return asynq.RedisClientOpt{
		Addr:     opt.Addr,
		Username: opt.Username,
		Password: opt.Password,
		DB:       0,
	}
}
