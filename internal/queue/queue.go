package queue

import (
	"crypto/subtle"
	"fmt"
	"time"

	"github.com/labstack/echo/v4/middleware"

	"github.com/labstack/echo/v4"

	"github.com/hibiken/asynq"
)

type Queue struct {
	Client    *asynq.Client
	Server    *asynq.ServeMux
	Scheduler *asynq.Scheduler
}

func Init(redisURL, username, password string, concurrency int, e *echo.Echo, isEnvRelease bool) *Queue {
	// redis connection
	redisConn := getRedisConnFromURL(redisURL)

	// init instance
	q := &Queue{}
	q.Server = initServer(redisConn, concurrency)
	q.Scheduler = initScheduler(redisConn)
	q.Client = initClient(redisConn)

	// init dashboard
	e.Any(fmt.Sprintf("%s/*", dashboardPath), echo.WrapHandler(dashboard(redisURL)), middleware.BasicAuth(func(username, password string, c echo.Context) (bool, error) {
		if !isEnvRelease {
			return true, nil
		}

		return subtle.ConstantTimeCompare([]byte(username), []byte(username)) == 1 &&
			subtle.ConstantTimeCompare([]byte(password), []byte(password)) == 1, nil
	}))

	fmt.Printf("⚡️ [queue]: initialized \n")

	return q
}

func initServer(redisConn asynq.RedisClientOpt, concurrency int) *asynq.ServeMux {
	// set unrecognized for concurrency
	if concurrency == 0 {
		concurrency = 100
	}

	// unrecognized retry delay in 10s
	retryDelayFunc := func(n int, e error, t *asynq.Task) time.Duration {
		return 10 * time.Second
	}

	// init server
	server := asynq.NewServer(redisConn, asynq.Config{
		Concurrency:    concurrency,
		RetryDelayFunc: retryDelayFunc,
		Queues: map[string]int{
			queueCronjob: 10,
			queueDefault: 3,
		},
	})

	// init mux server
	mux := asynq.NewServeMux()

	// run server
	go func() {
		if err := server.Run(mux); err != nil {
			msg := fmt.Sprintf("error when initializing queue SERVER: %s", err.Error())
			panic(msg)
		}
	}()

	return mux
}

func initScheduler(redisConn asynq.RedisClientOpt) *asynq.Scheduler {
	// always run at HCM timezone
	l, _ := time.LoadLocation("Asia/Ho_Chi_Minh")

	// init scheduler
	scheduler := asynq.NewScheduler(redisConn, &asynq.SchedulerOpts{
		Location: l,
	})

	// run scheduler
	go func() {
		if err := scheduler.Run(); err != nil {
			msg := fmt.Sprintf("error when initializing queue SCHEDULER: %s", err.Error())
			panic(msg)
		}
	}()

	return scheduler
}

func initClient(redisConn asynq.RedisClientOpt) *asynq.Client {
	client := asynq.NewClient(redisConn)
	if client == nil {
		panic("error when initializing queue CLIENT")
	}
	return client
}
