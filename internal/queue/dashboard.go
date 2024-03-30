package queue

import "github.com/hibiken/asynqmon"

const dashboardPath = "/q"

func dashboard(redisURL string) *asynqmon.HTTPHandler {
	redisConn := getRedisConnFromURL(redisURL)

	return asynqmon.New(asynqmon.Options{
		RootPath:     dashboardPath,
		RedisConnOpt: redisConn,
	})
}
