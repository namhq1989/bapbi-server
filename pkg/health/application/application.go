package application

import (
	"github.com/namhq1989/bapbi-server/internal/utils/appcontext"
	"github.com/namhq1989/bapbi-server/pkg/health/application/command"
	"github.com/namhq1989/bapbi-server/pkg/health/application/query"
	"github.com/namhq1989/bapbi-server/pkg/health/domain"
	"github.com/namhq1989/bapbi-server/pkg/health/dto"
)

type (
	Commands interface {
		CreateHealthProfile(ctx *appcontext.AppContext, performerID string, req dto.CreateHealthProfileRequest) (*dto.CreateHealthProfileResponse, error)

		EnableHydrationProfile(ctx *appcontext.AppContext, performerID string, _ dto.EnableHydrationProfileRequest) (*dto.EnableHydrationProfileResponse, error)
		DisableHydrationProfile(ctx *appcontext.AppContext, performerID string, _ dto.DisableHydrationProfileRequest) (*dto.DisableHydrationProfileResponse, error)
		WaterIntake(ctx *appcontext.AppContext, performerID string, req dto.WaterIntakeRequest) (*dto.WaterIntakeResponse, error)
	}
	Queries interface {
		HydrationStats(ctx *appcontext.AppContext, performerID string, _ dto.HydrationStatsRequest) (*dto.HydrationStatsResponse, error)
	}
	Hubs interface{}
	App  interface {
		Commands
		Queries
		Hubs
	}

	appCommandHandlers struct {
		command.CreateHealthProfileHandler

		command.EnableHydrationProfileHandler
		command.DisableHydrationProfileHandler
		command.WaterIntakeHandler
	}
	appQueryHandler struct {
		query.HydrationStatsHandler
	}
	appHubHandler struct{}
	Application   struct {
		appCommandHandlers
		appQueryHandler
		appHubHandler
	}
)

var _ App = (*Application)(nil)

func New(
	healthProfileRepository domain.HealthProfileRepository,
	hydrationProfileRepository domain.HydrationProfileRepository,
	hydrationDailyReportRepository domain.HydrationDailyReportRepository,
	waterIntakeLogRepository domain.WaterIntakeLogRepository,
	queueRepository domain.QueueRepository,
) *Application {
	return &Application{
		appCommandHandlers: appCommandHandlers{
			CreateHealthProfileHandler: command.NewCreateHealthProfileHandler(healthProfileRepository),

			EnableHydrationProfileHandler:  command.NewEnableHydrationProfileHandler(healthProfileRepository, hydrationProfileRepository),
			DisableHydrationProfileHandler: command.NewDisableHydrationProfileHandler(healthProfileRepository, hydrationProfileRepository),
			WaterIntakeHandler:             command.NewWaterIntakeHandler(hydrationProfileRepository, waterIntakeLogRepository, queueRepository),
		},
		appQueryHandler: appQueryHandler{
			HydrationStatsHandler: query.NewHydrationStatsHandler(hydrationProfileRepository, hydrationDailyReportRepository, waterIntakeLogRepository),
		},
		appHubHandler: appHubHandler{},
	}
}
