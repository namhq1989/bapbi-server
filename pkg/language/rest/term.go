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
		return validation.ValidateHTTPBody[dto.SearchTermRequest](next)
	})
}
