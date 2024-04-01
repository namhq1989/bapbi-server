package health

import (
	"github.com/namhq1989/bapbi-server/internal/monolith"
	"github.com/namhq1989/bapbi-server/internal/utils/appcontext"
	"github.com/namhq1989/bapbi-server/pkg/health/application"
	"github.com/namhq1989/bapbi-server/pkg/health/infrastructure"
	"github.com/namhq1989/bapbi-server/pkg/health/rest"
)

type Module struct{}

func (Module) Name() string {
	return "HEALTH"
}

func (Module) Startup(ctx *appcontext.AppContext, mono monolith.Monolith) error {
	var (
		// infrastructure
		healthProfileRepository     = infrastructure.NewHealthProfileRepository(mono.Mongo())
		drinkWaterProfileRepository = infrastructure.NewDrinkWaterProfileRepository(mono.Mongo())

		// application
		app = application.New(healthProfileRepository, drinkWaterProfileRepository)
	)

	// rest server
	if err := rest.RegisterServer(ctx, app, mono.Rest(), mono.JWT()); err != nil {
		return err
	}

	return nil
}
