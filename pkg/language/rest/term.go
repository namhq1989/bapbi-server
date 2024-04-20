package rest

import (
	"github.com/labstack/echo/v4"
	"github.com/namhq1989/bapbi-server/internal/utils/appcontext"
	"github.com/namhq1989/bapbi-server/internal/utils/httprespond"
	"github.com/namhq1989/bapbi-server/internal/utils/validation"
	"github.com/namhq1989/bapbi-server/pkg/language/dto"
)

func (s server) registerTermRoutes() {
	g := s.echo.Group("/api/language/term")

	g.GET("/search", func(c echo.Context) error {
		var (
			ctx         = c.Get("ctx").(*appcontext.AppContext)
			req         = c.Get("req").(dto.SearchTermRequest)
			performerID = ctx.GetUserID()
		)

		resp, err := s.app.SearchTerm(ctx, performerID, req)
		if err != nil {
			return httprespond.R400(c, err, nil)
		}

		return httprespond.R200(c, resp)
	}, s.jwt.RequireLoggedIn, func(next echo.HandlerFunc) echo.HandlerFunc {
		return validation.ValidateHTTPPayload[dto.SearchTermRequest](next)
	})

	g.POST("/:id/add", func(c echo.Context) error {
		var (
			ctx         = c.Get("ctx").(*appcontext.AppContext)
			req         = c.Get("req").(dto.AddTermRequest)
			termID      = c.Param("id")
			performerID = ctx.GetUserID()
		)

		resp, err := s.app.AddTerm(ctx, performerID, termID, req)
		if err != nil {
			return httprespond.R400(c, err, nil)
		}

		return httprespond.R200(c, resp)
	}, s.jwt.RequireLoggedIn, func(next echo.HandlerFunc) echo.HandlerFunc {
		return validation.ValidateHTTPPayload[dto.AddTermRequest](next)
	})

	g.PATCH("/:userTermId/favourite", func(c echo.Context) error {
		var (
			ctx         = c.Get("ctx").(*appcontext.AppContext)
			req         = c.Get("req").(dto.ChangeTermFavouriteRequest)
			userTermID  = c.Param("userTermId")
			performerID = ctx.GetUserID()
		)

		resp, err := s.app.ChangeTermFavourite(ctx, performerID, userTermID, req)
		if err != nil {
			return httprespond.R400(c, err, nil)
		}

		return httprespond.R200(c, resp)
	}, s.jwt.RequireLoggedIn, func(next echo.HandlerFunc) echo.HandlerFunc {
		return validation.ValidateHTTPPayload[dto.ChangeTermFavouriteRequest](next)
	})
}
