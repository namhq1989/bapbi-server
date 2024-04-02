package rest

import (
	"github.com/labstack/echo/v4"
	"github.com/namhq1989/bapbi-server/internal/utils/appcontext"
	"github.com/namhq1989/bapbi-server/internal/utils/httprespond"
	"github.com/namhq1989/bapbi-server/internal/utils/validation"
	"github.com/namhq1989/bapbi-server/pkg/health/dto"
)

func (s server) registerHydrationRoutes() {
	g := s.echo.Group("/api/health/hydration")

	g.PATCH("/enable-profile", func(c echo.Context) error {
		var (
			ctx         = c.Get("ctx").(*appcontext.AppContext)
			req         = c.Get("req").(dto.EnableHydrationProfileRequest)
			performerID = ctx.GetUserID()
		)

		resp, err := s.app.EnableHydrationProfile(ctx, performerID, req)
		if err != nil {
			return httprespond.R400(c, err, nil)
		}

		return httprespond.R200(c, resp)
	}, s.jwt.RequireLoggedIn, func(next echo.HandlerFunc) echo.HandlerFunc {
		return validation.ValidateHTTPBody[dto.EnableHydrationProfileRequest](next)
	})

	g.PATCH("/disable-profile", func(c echo.Context) error {
		var (
			ctx         = c.Get("ctx").(*appcontext.AppContext)
			req         = c.Get("req").(dto.DisableHydrationProfileRequest)
			performerID = ctx.GetUserID()
		)

		resp, err := s.app.DisableHydrationProfile(ctx, performerID, req)
		if err != nil {
			return httprespond.R400(c, err, nil)
		}

		return httprespond.R200(c, resp)
	}, s.jwt.RequireLoggedIn, func(next echo.HandlerFunc) echo.HandlerFunc {
		return validation.ValidateHTTPBody[dto.DisableHydrationProfileRequest](next)
	})

	g.POST("/water-intake", func(c echo.Context) error {
		var (
			ctx         = c.Get("ctx").(*appcontext.AppContext)
			req         = c.Get("req").(dto.WaterIntakeRequest)
			performerID = ctx.GetUserID()
		)

		resp, err := s.app.WaterIntake(ctx, performerID, req)
		if err != nil {
			return httprespond.R400(c, err, nil)
		}

		return httprespond.R200(c, resp)
	}, s.jwt.RequireLoggedIn, func(next echo.HandlerFunc) echo.HandlerFunc {
		return validation.ValidateHTTPBody[dto.WaterIntakeRequest](next)
	})
}
