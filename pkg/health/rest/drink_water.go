package rest

import (
	"github.com/labstack/echo/v4"
	"github.com/namhq1989/bapbi-server/internal/utils/appcontext"
	"github.com/namhq1989/bapbi-server/internal/utils/httprespond"
	"github.com/namhq1989/bapbi-server/internal/utils/validation"
	"github.com/namhq1989/bapbi-server/pkg/health/dto"
)

func (s server) registerDrinkWaterRoutes() {
	g := s.echo.Group("/api/health/drink-water")

	g.PATCH("/enable-profile", func(c echo.Context) error {
		var (
			ctx         = c.Get("ctx").(*appcontext.AppContext)
			req         = c.Get("req").(dto.EnableDrinkWaterProfileRequest)
			performerID = ctx.GetUserID()
		)

		resp, err := s.app.EnableDrinkWaterProfile(ctx, performerID, req)
		if err != nil {
			return httprespond.R400(c, err, nil)
		}

		return httprespond.R200(c, resp)
	}, s.jwt.RequireLoggedIn, func(next echo.HandlerFunc) echo.HandlerFunc {
		return validation.ValidateHTTPBody[dto.EnableDrinkWaterProfileRequest](next)
	})

	g.PATCH("/disable-profile", func(c echo.Context) error {
		var (
			ctx         = c.Get("ctx").(*appcontext.AppContext)
			req         = c.Get("req").(dto.DisableDrinkWaterProfileRequest)
			performerID = ctx.GetUserID()
		)

		resp, err := s.app.DisableDrinkWaterProfile(ctx, performerID, req)
		if err != nil {
			return httprespond.R400(c, err, nil)
		}

		return httprespond.R200(c, resp)
	}, s.jwt.RequireLoggedIn, func(next echo.HandlerFunc) echo.HandlerFunc {
		return validation.ValidateHTTPBody[dto.DisableDrinkWaterProfileRequest](next)
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
