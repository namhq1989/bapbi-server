package rest

import (
	"github.com/labstack/echo/v4"
	"github.com/namhq1989/bapbi-server/internal/utils/appcontext"
	appjwt "github.com/namhq1989/bapbi-server/internal/utils/jwt"
	"github.com/namhq1989/bapbi-server/pkg/health/application"
)

type server struct {
	app  application.App
	echo *echo.Echo
	jwt  *appjwt.JWT
}

func RegisterServer(_ *appcontext.AppContext, app application.App, e *echo.Echo, jwt *appjwt.JWT) error {
	var s = server{
		app:  app,
		echo: e,
		jwt:  jwt,
	}

	s.registerHealthProfileRoutes()
	s.registerDrinkWaterRoutes()

	return nil
}
