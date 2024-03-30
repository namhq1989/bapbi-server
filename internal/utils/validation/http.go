package validation

import (
	"github.com/gookit/validate"
	"github.com/labstack/echo/v4"
	"github.com/namhq1989/bapbi-server/internal/utils/httprespond"
)

func ValidateHTTPQuery[T any](next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		var (
			query T
		)

		if err := c.Bind(&query); err != nil {
			return httprespond.R400(c, err, echo.Map{})
		}

		if v := validate.Struct(query); !v.Validate() {
			return httprespond.R400(c, v.Errors.OneError(), echo.Map{})
		}

		// assign to context
		c.Set("req", query)
		return next(c)
	}
}

func ValidateHTTPBody[T any](next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		var (
			body T
		)

		if err := c.Bind(&body); err != nil {
			return httprespond.R400(c, err, echo.Map{})
		}

		if v := validate.Struct(body); !v.Validate() {
			return httprespond.R400(c, v.Errors.OneError(), echo.Map{})
		}

		// assign to context
		c.Set("req", body)
		return next(c)
	}
}
