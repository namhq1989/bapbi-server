package application

import (
	"github.com/namhq1989/bapbi-server/internal/utils/appcontext"
	"github.com/namhq1989/bapbi-server/pkg/auth/application/command"
	"github.com/namhq1989/bapbi-server/pkg/auth/application/query"
	"github.com/namhq1989/bapbi-server/pkg/auth/domain"
	"github.com/namhq1989/bapbi-server/pkg/auth/dto"
)

type (
	Commands interface {
		LoginWithGoogle(ctx *appcontext.AppContext, req dto.LoginWithGoogleRequest) (*dto.LoginWithGoogleResponse, error)
		RefreshAccessToken(ctx *appcontext.AppContext, req dto.RefreshAccessTokenRequest) (*dto.RefreshAccessTokenResponse, error)
	}
	Queries interface {
		GetAccessTokenByUserId(ctx *appcontext.AppContext, req dto.GetAccessTokenByUserIDRequest) (*dto.GetAccessTokenByUserIDResponse, error)
		Me(ctx *appcontext.AppContext, req dto.MeRequest) (*dto.MeResponse, error)
	}
	Hubs interface {
	}
	App interface {
		Commands
		Queries
		Hubs
	}

	appCommandHandlers struct {
		command.LoginWithGoogleHandler
		command.RefreshAccessTokenHandler
	}
	appQueryHandler struct {
		query.GetAccessTokenByUserIdHandler
		query.MeHandler
	}
	appHubHandler struct {
	}
	Application struct {
		appCommandHandlers
		appQueryHandler
		appHubHandler
	}
)

var _ App = (*Application)(nil)

func New(
	ssoRepository domain.SSORepository,
	authTokenRepository domain.AuthTokenRepository,
	userHub domain.UserHub,
	jwtRepository domain.JwtRepository,
) *Application {
	return &Application{
		appCommandHandlers: appCommandHandlers{
			LoginWithGoogleHandler:    command.NewLoginWithGoogleHandler(ssoRepository, authTokenRepository, userHub, jwtRepository),
			RefreshAccessTokenHandler: command.NewRefreshAccessTokenHandler(authTokenRepository, userHub, jwtRepository),
		},
		appQueryHandler: appQueryHandler{
			GetAccessTokenByUserIdHandler: query.NewGetAccessTokenByUserIdHandler(jwtRepository),
			MeHandler:                     query.NewMeHandler(userHub),
		},
		appHubHandler: appHubHandler{},
	}
}
