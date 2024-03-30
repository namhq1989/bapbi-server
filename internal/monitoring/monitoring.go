package monitoring

import (
	"fmt"
	"github.com/getsentry/sentry-go"
	sentryecho "github.com/getsentry/sentry-go/echo"
	"github.com/labstack/echo/v4"
)

type Monitoring struct {
	sentry *sentry.Client
}

func Init(e *echo.Echo, dsn, machine, environment string) *Monitoring {
	// skip if the "machine" is not set
	if machine == "" {
		fmt.Printf("⚡️ [monitoring]: machine is not set \n")
		return nil
	}

	opts := sentry.ClientOptions{
		Dsn:                dsn,
		Environment:        fmt.Sprintf("%s-%s", environment, machine),
		EnableTracing:      true,
		TracesSampleRate:   1.0,
		ProfilesSampleRate: 1.0,
		IgnoreTransactions: []string{
			"/q/*",
		},
		// Debug: true,
	}

	if err := sentry.Init(opts); err != nil {
		panic(err)
	}

	client, err := sentry.NewClient(opts)
	if err != nil {
		panic(err)
	}

	// use as middleware
	e.Use(sentryecho.New(sentryecho.Options{}))
	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			ctx := c.Request().Context()
			hub := sentry.GetHubFromContext(ctx)
			if hub == nil {
				hub = sentry.CurrentHub().Clone()
				ctx = sentry.SetHubOnContext(ctx, hub)
			}

			options := []sentry.SpanOption{
				sentry.WithOpName("http.server"),
				sentry.ContinueFromRequest(c.Request()),
				sentry.WithTransactionSource(sentry.SourceURL),
			}

			transaction := sentry.StartTransaction(ctx,
				fmt.Sprintf("%s %s", c.Request().Method, c.Request().URL.Path),
				options...,
			)
			defer transaction.Finish()

			return next(c)
		}
	})

	// recover
	defer sentry.Recover()

	fmt.Printf("⚡️ [monitoring]: connected \n")

	return &Monitoring{sentry: client}
}
