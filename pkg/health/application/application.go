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
) *Application {
	return &Application{
		appCommandHandlers: appCommandHandlers{
			CreateHealthProfileHandler: command.NewCreateHealthProfileHandler(healthProfileRepository),
		},
		appQueryHandler: appQueryHandler{},
		appHubHandler:   appHubHandler{},
	}
}
