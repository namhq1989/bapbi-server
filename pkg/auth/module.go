package auth

import (
	"github.com/namhq1989/bapbi-server/internal/grpcclient"
	"github.com/namhq1989/bapbi-server/internal/monolith"
	"github.com/namhq1989/bapbi-server/internal/utils/appcontext"
	"github.com/namhq1989/bapbi-server/pkg/auth/application"
	"github.com/namhq1989/bapbi-server/pkg/auth/grpc"
	"github.com/namhq1989/bapbi-server/pkg/auth/infrastructure"
	"github.com/namhq1989/bapbi-server/pkg/auth/rest"
)

type Module struct{}

func (Module) Name() string {
	return "Auth"
}

func (Module) Startup(ctx *appcontext.AppContext, mono monolith.Monolith) error {
	userGRPCClient, err := grpcclient.NewUserClient(ctx, mono.Config().GRPCPort)
	if err != nil {
		return err
	}

	var (
		cfg = mono.Config()

		ssoRepository       = infrastructure.NewSSORepository(cfg.SSOGoogleClientID, cfg.SSOGoogleClientSecret)
		authTokenRepository = infrastructure.NewAuthTokenRepository(mono.Mongo())
		userHub             = infrastructure.NewUserHub(userGRPCClient)
		jwtRepository       = infrastructure.NewJwtRepository(mono.JWT())

		// app
		app = application.New(ssoRepository, authTokenRepository, userHub, jwtRepository)
	)

	// rest server
	if err = rest.RegisterServer(ctx, app, mono.Rest(), mono.JWT(), cfg.IsEnvRelease); err != nil {
		return err
	}

	// grpc server
	if err = grpc.RegisterServer(ctx, app, mono.RPC()); err != nil {
		return err
	}

	return nil
}
