package rest

import (
	"github.com/labstack/echo/v4"
	"github.com/namhq1989/bapbi-server/internal/utils/appcontext"
	"github.com/namhq1989/bapbi-server/internal/utils/httprespond"
	"github.com/namhq1989/bapbi-server/internal/utils/validation"
	"github.com/namhq1989/bapbi-server/pkg/language/dto"
)

func (s server) registerExerciseRoutes() {
	g := s.echo.Group("/api/language/exercise")

	g.GET("/writing", func(c echo.Context) error {
		var (
			ctx         = c.Get("ctx").(*appcontext.AppContext)
			req         = c.Get("req").(dto.GetWritingExerciseRequest)
			performerID = ctx.GetUserID()
		)

		resp, err := s.app.GetWritingExercises(ctx, performerID, req)
		if err != nil {
			return httprespond.R400(c, err, nil)
		}

		return httprespond.R200(c, resp)
	}, s.jwt.RequireLoggedIn, func(next echo.HandlerFunc) echo.HandlerFunc {
		return validation.ValidateHTTPPayload[dto.GetWritingExerciseRequest](next)
	})

	g.POST("/writing", func(c echo.Context) error {
		var (
			ctx         = c.Get("ctx").(*appcontext.AppContext)
			req         = c.Get("req").(dto.StartUserWritingExerciseRequest)
			performerID = ctx.GetUserID()
		)

		resp, err := s.app.StartUserWritingExercise(ctx, performerID, req)
		if err != nil {
			return httprespond.R400(c, err, nil)
		}

		return httprespond.R200(c, resp)
	}, s.jwt.RequireLoggedIn, func(next echo.HandlerFunc) echo.HandlerFunc {
		return validation.ValidateHTTPPayload[dto.StartUserWritingExerciseRequest](next)
	})

	g.PUT("/submit-writing", func(c echo.Context) error {
		var (
			ctx         = c.Get("ctx").(*appcontext.AppContext)
			req         = c.Get("req").(dto.SubmitUserWritingExerciseRequest)
			performerID = ctx.GetUserID()
		)

		resp, err := s.app.SubmitUserWritingExercise(ctx, performerID, req)
		if err != nil {
			return httprespond.R400(c, err, nil)
		}

		return httprespond.R200(c, resp)
	}, s.jwt.RequireLoggedIn, func(next echo.HandlerFunc) echo.HandlerFunc {
		return validation.ValidateHTTPPayload[dto.SubmitUserWritingExerciseRequest](next)
	})

	g.PUT("/modify-writing", func(c echo.Context) error {
		var (
			ctx         = c.Get("ctx").(*appcontext.AppContext)
			req         = c.Get("req").(dto.ModifyUserWritingExerciseRequest)
			performerID = ctx.GetUserID()
		)

		resp, err := s.app.ModifyUserWritingExercise(ctx, performerID, req)
		if err != nil {
			return httprespond.R400(c, err, nil)
		}

		return httprespond.R200(c, resp)
	}, s.jwt.RequireLoggedIn, func(next echo.HandlerFunc) echo.HandlerFunc {
		return validation.ValidateHTTPPayload[dto.ModifyUserWritingExerciseRequest](next)
	})

	g.GET("/user-writing", func(c echo.Context) error {
		var (
			ctx         = c.Get("ctx").(*appcontext.AppContext)
			req         = c.Get("req").(dto.GetUserWritingExerciseRequest)
			performerID = ctx.GetUserID()
		)

		resp, err := s.app.GetUserWritingExercises(ctx, performerID, req)
		if err != nil {
			return httprespond.R400(c, err, nil)
		}

		return httprespond.R200(c, resp)
	}, s.jwt.RequireLoggedIn, func(next echo.HandlerFunc) echo.HandlerFunc {
		return validation.ValidateHTTPPayload[dto.GetUserWritingExerciseRequest](next)
	})

	g.PUT("/submit-term", func(c echo.Context) error {
		var (
			ctx         = c.Get("ctx").(*appcontext.AppContext)
			req         = c.Get("req").(dto.SubmitUserTermExerciseRequest)
			performerID = ctx.GetUserID()
		)

		resp, err := s.app.SubmitUserTermExercise(ctx, performerID, req)
		if err != nil {
			return httprespond.R400(c, err, nil)
		}

		return httprespond.R200(c, resp)
	}, s.jwt.RequireLoggedIn, func(next echo.HandlerFunc) echo.HandlerFunc {
		return validation.ValidateHTTPPayload[dto.SubmitUserTermExerciseRequest](next)
	})

	g.PUT("/modify-term", func(c echo.Context) error {
		var (
			ctx         = c.Get("ctx").(*appcontext.AppContext)
			req         = c.Get("req").(dto.ModifyUserTermExerciseRequest)
			performerID = ctx.GetUserID()
		)

		resp, err := s.app.ModifyUserTermExercise(ctx, performerID, req)
		if err != nil {
			return httprespond.R400(c, err, nil)
		}

		return httprespond.R200(c, resp)
	}, s.jwt.RequireLoggedIn, func(next echo.HandlerFunc) echo.HandlerFunc {
		return validation.ValidateHTTPPayload[dto.ModifyUserTermExerciseRequest](next)
	})

	g.GET("/user-term", func(c echo.Context) error {
		var (
			ctx         = c.Get("ctx").(*appcontext.AppContext)
			req         = c.Get("req").(dto.GetUserTermExerciseRequest)
			performerID = ctx.GetUserID()
		)

		resp, err := s.app.GetUserTermExercises(ctx, performerID, req)
		if err != nil {
			return httprespond.R400(c, err, nil)
		}

		return httprespond.R200(c, resp)
	}, s.jwt.RequireLoggedIn, func(next echo.HandlerFunc) echo.HandlerFunc {
		return validation.ValidateHTTPPayload[dto.GetUserTermExerciseRequest](next)
	})
}
