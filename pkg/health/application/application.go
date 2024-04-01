package application

import (
	"github.com/namhq1989/bapbi-server/internal/utils/appcontext"
	"github.com/namhq1989/bapbi-server/pkg/health/application/command"
	"github.com/namhq1989/bapbi-server/pkg/health/domain"
	"github.com/namhq1989/bapbi-server/pkg/health/dto"
)

type (
	Commands interface {
		CreateHealthProfile(ctx *appcontext.AppContext, performerID string, req dto.CreateHealthProfileRequest) (*dto.CreateHealthProfileResponse, error)

		EnableDrinkWaterProfile(ctx *appcontext.AppContext, performerID string, req dto.EnableDrinkWaterProfileRequest) (*dto.EnableDrinkWaterProfileResponse, error)
		DisableDrinkWaterProfile(ctx *appcontext.AppContext, performerID string, _ dto.DisableDrinkWaterProfileRequest) (*dto.DisableDrinkWaterProfileResponse, error)
		WaterIntake(ctx *appcontext.AppContext, performerID string, req dto.WaterIntakeRequest) (*dto.WaterIntakeResponse, error)
	}
	Queries interface {
	}
	Hubs interface{}
	App  interface {
		Commands
		Queries
		Hubs
	}

	appCommandHandlers struct {
		command.CreateHealthProfileHandler

		command.EnableDrinkWaterProfileHandler
		command.DisableDrinkWaterProfileHandler
		command.WaterIntakeHandler
	}
	appQueryHandler struct {
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
	drinkWaterProfileRepository domain.DrinkWaterProfileRepository,
	waterIntakeLogRepository domain.WaterIntakeLogRepository,
	queueRepository domain.QueueRepository,
) *Application {
	return &Application{
		appCommandHandlers: appCommandHandlers{
			CreateHealthProfileHandler: command.NewCreateHealthProfileHandler(healthProfileRepository),

			EnableDrinkWaterProfileHandler:  command.NewEnableDrinkWaterProfileHandler(healthProfileRepository, drinkWaterProfileRepository),
			DisableDrinkWaterProfileHandler: command.NewDisableDrinkWaterProfileHandler(healthProfileRepository, drinkWaterProfileRepository),
			WaterIntakeHandler:              command.NewWaterIntakeHandler(drinkWaterProfileRepository, waterIntakeLogRepository, queueRepository),
		},
		appQueryHandler: appQueryHandler{},
		appHubHandler:   appHubHandler{},
	}
}
