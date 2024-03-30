package application

import (
	"github.com/namhq1989/bapbi-server/internal/genproto/userpb"
	"github.com/namhq1989/bapbi-server/internal/utils/appcontext"
	"github.com/namhq1989/bapbi-server/pkg/user/application/hub"
	"github.com/namhq1989/bapbi-server/pkg/user/domain"
)

type (
	Commands interface {
	}
	Queries interface {
	}
	Hubs interface {
		CreateUser(ctx *appcontext.AppContext, req *userpb.CreateUserRequest) (*userpb.CreateUserResponse, error)
		GetUserByID(ctx *appcontext.AppContext, req *userpb.GetUserByIDRequest) (*userpb.GetUserByIDResponse, error)
		GetUserByEmail(ctx *appcontext.AppContext, req *userpb.GetUserByEmailRequest) (*userpb.GetUserByEmailResponse, error)
	}
	App interface {
		Commands
		Queries
		Hubs
	}

	appCommandHandlers struct {
	}
	appQueryHandler struct {
	}
	appHubHandler struct {
		hub.CreateUserHandler
		hub.GetUserByIDHandler
		hub.GetUserByEmailHandler
	}
	Application struct {
		appCommandHandlers
		appQueryHandler
		appHubHandler
	}
)

var _ App = (*Application)(nil)

func New(
	queueRepository domain.QueueRepository,
	userHub domain.UserHub,
) *Application {
	return &Application{
		appCommandHandlers: appCommandHandlers{},
		appQueryHandler:    appQueryHandler{},
		appHubHandler: appHubHandler{
			CreateUserHandler:     hub.NewCreateUserHandler(queueRepository, userHub),
			GetUserByIDHandler:    hub.NewGetUserByIDHandler(userHub),
			GetUserByEmailHandler: hub.NewGetUserByEmailHandler(userHub),
		},
	}
}
