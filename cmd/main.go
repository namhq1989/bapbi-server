package main

import (
	"fmt"
	"time"

	"github.com/namhq1989/bapbi-server/internal/caching"
	"github.com/namhq1989/bapbi-server/internal/config"
	"github.com/namhq1989/bapbi-server/internal/database"
	"github.com/namhq1989/bapbi-server/internal/monitoring"
	"github.com/namhq1989/bapbi-server/internal/monolith"
	"github.com/namhq1989/bapbi-server/internal/queue"
	apperrors "github.com/namhq1989/bapbi-server/internal/utils/error"
	appjwt "github.com/namhq1989/bapbi-server/internal/utils/jwt"
	"github.com/namhq1989/bapbi-server/internal/utils/logger"
	"github.com/namhq1989/bapbi-server/internal/utils/waiter"
	"github.com/namhq1989/bapbi-server/pkg/auth"
	"github.com/namhq1989/bapbi-server/pkg/user"
)

func main() {
	var err error

	// config
	cfg := config.Init()

	// logger
	logger.Init(cfg.Environment)

	// app error
	apperrors.Init()

	// server
	a := app{}
	a.cfg = cfg

	// rest
	a.rest = initRest(cfg)

	// grpc
	a.rpc = initRPC()

	// jwt
	a.jwt, err = appjwt.Init(cfg.AccessTokenSecret, cfg.RefreshTokenSecret, time.Second*time.Duration(cfg.AccessTokenTTL), time.Second*time.Duration(cfg.RefreshTokenTTL))
	if err != nil {
		panic(err)
	}

	// mongodb
	mgClient := database.NewMongoClient(cfg.MongoURL)
	a.mongo = mgClient.Database(cfg.MongoDBName)

	// queue
	a.queue = queue.Init(cfg.RedisURL, cfg.QueueUsername, cfg.QueuePassword, cfg.QueueConcurrency, a.Rest(), cfg.IsEnvRelease)

	// caching
	a.caching = caching.NewCachingClient(cfg.RedisURL)

	// monitoring
	a.monitoring = monitoring.Init(a.Rest(), cfg.SentryDSN, cfg.SentryMachine, cfg.Environment)

	// waiter
	a.waiter = waiter.New(waiter.CatchSignals())

	// modules
	a.modules = []monolith.Module{
		&auth.Module{},
		&user.Module{},
	}

	// start
	if err = a.startupModules(); err != nil {
		panic(err)
	}

	fmt.Println("--- started server-bapbi application")
	defer fmt.Println("--- stopped server-bapbi application")

	// wait for other service starts
	a.waiter.Add(
		a.waitForRest,
		a.waitForRPC,
	)
	if err = a.waiter.Wait(); err != nil {
		panic(err)
	}
}
