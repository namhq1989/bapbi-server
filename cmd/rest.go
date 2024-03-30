package main

import (
	"net/http"
	"time"

	"github.com/namhq1989/bapbi-server/internal/utils/appcontext"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/namhq1989/bapbi-server/internal/config"
	"golang.org/x/text/language"
)

func initRest(cfg config.Server) *echo.Echo {
	// echo instance
	e := echo.New()

	// middlewares
	setMiddleware(e, cfg)

	return e
}

func setMiddleware(e *echo.Echo, cfg config.Server) {
	addCorsMiddleware(e)
	addContext(e)
	addIp(e)
	addLanguageMiddleware(e)
	addRateLimiter(e)
	e.Use(middleware.GzipWithConfig(middleware.GzipConfig{
		Level: 5,
	}))
	e.Use(middleware.Secure())

	if cfg.IsEnvRelease {
		e.Use(middleware.Recover())
	} else {
		e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
			Format: "${time_rfc3339} | ${remote_ip} | ${method} ${uri} - ${status} - ${latency_human}\n",
		}))
	}
}

func addContext(e *echo.Echo) {
	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			ctx := appcontext.New(c.Request().Context())
			c.Set("ctx", ctx)

			return next(c)
		}
	})
}

func addIp(e *echo.Echo) {
	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			var (
				ctx = c.Get("ctx").(*appcontext.AppContext)
				ip  = c.RealIP()
			)

			ctx.SetIP(ip)
			return next(c)
		}
	})
}

func addLanguageMiddleware(e *echo.Echo) {
	supportedLanguages := language.NewMatcher([]language.Tag{
		language.English,
		language.Vietnamese,
	})

	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// parse the Accept-Language header
			accept := c.Request().Header.Get("Accept-Language")
			tag, _, _ := language.ParseAcceptLanguage(accept)

			// match the best supported language
			matched, _, _ := supportedLanguages.Match(tag...)

			// Use "en" as default if no match
			lang := language.English.String()
			if matched == language.Vietnamese {
				lang = language.Vietnamese.String()
			}

			// set the language in the context
			c.Set("lang", lang)

			// Call the next handler in the chain
			return next(c)
		}
	})
}

func addCorsMiddleware(e *echo.Echo) {
	allowedOrigins := []string{
		"http://localhost:5173",
		"http://127.0.0.1:5173",
	}
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     allowedOrigins,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS"},
		AllowHeaders:     []string{"Accept", "Authorization", "Content-Type"},
		AllowCredentials: true,
		MaxAge:           300,
	}))
}

func addRateLimiter(e *echo.Echo) {
	cfg := middleware.RateLimiterConfig{
		Skipper: middleware.DefaultSkipper,
		Store: middleware.NewRateLimiterMemoryStoreWithConfig(
			middleware.RateLimiterMemoryStoreConfig{Rate: 60, Burst: 60, ExpiresIn: 5 * time.Minute},
		),
		IdentifierExtractor: func(ctx echo.Context) (string, error) {
			id := ctx.RealIP()
			return id, nil
		},
		ErrorHandler: func(context echo.Context, err error) error {
			return context.JSON(http.StatusForbidden, nil)
		},
		DenyHandler: func(context echo.Context, identifier string, err error) error {
			return context.JSON(http.StatusTooManyRequests, nil)
		},
	}
	e.Use(middleware.RateLimiterWithConfig(cfg))
}
