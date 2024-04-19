package health

import (
	"github.com/namhq1989/bapbi-server/internal/monolith"
	"github.com/namhq1989/bapbi-server/internal/utils/appcontext"
	"github.com/namhq1989/bapbi-server/pkg/health/application"
	"github.com/namhq1989/bapbi-server/pkg/health/infrastructure"
	"github.com/namhq1989/bapbi-server/pkg/health/rest"
	"github.com/namhq1989/bapbi-server/pkg/health/worker"
)

type Module struct{}

func (Module) Name() string {
	return "HEALTH"
}

func (Module) Startup(ctx *appcontext.AppContext, mono monolith.Monolith) error {
	var (
		// infrastructure
		healthProfileRepository        = infrastructure.NewHealthProfileRepository(mono.Mongo())
		hydrationProfileRepository     = infrastructure.NewHydrationProfileRepository(mono.Mongo())
		waterIntakeLogRepository       = infrastructure.NewWaterIntakeLogRepository(mono.Mongo())
		hydrationDailyReportRepository = infrastructure.NewHydrationDailyReportRepository(mono.Mongo())
		queueRepository                = infrastructure.NewQueueRepository(mono.Queue())

		// application
		app = application.New(healthProfileRepository, hydrationProfileRepository, hydrationDailyReportRepository, waterIntakeLogRepository, queueRepository)
	)

	// rest server
	if err := rest.RegisterServer(ctx, app, mono.Rest(), mono.JWT()); err != nil {
		return err
	}

	// workers
	w := worker.New(mono.Queue(), healthProfileRepository, hydrationProfileRepository, hydrationDailyReportRepository)
	w.Start()

	return nil
}
