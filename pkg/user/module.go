package user

import (
	"github.com/namhq1989/bapbi-server/internal/monolith"
	"github.com/namhq1989/bapbi-server/internal/utils/appcontext"
	"github.com/namhq1989/bapbi-server/pkg/user/application"
	"github.com/namhq1989/bapbi-server/pkg/user/grpc"
	"github.com/namhq1989/bapbi-server/pkg/user/infrastructure"
	"github.com/namhq1989/bapbi-server/pkg/user/rest"
	"github.com/namhq1989/bapbi-server/pkg/user/workers"
)

type Module struct{}

func (Module) Name() string {
	return "USER"
}

func (Module) Startup(ctx *appcontext.AppContext, mono monolith.Monolith) error {
	var (
		// infrastructure
		_               = infrastructure.NewUserRepository(mono.Mongo())
		queueRepository = infrastructure.NewQueueRepository(mono.Queue())
		userHub         = infrastructure.NewUserHub(mono.Mongo())

		// application
		app = application.New(queueRepository, userHub)
	)

	// rest server
	if err := rest.RegisterServer(ctx, app, mono.Rest(), mono.JWT()); err != nil {
		return err
	}

	// grpc server
	if err := grpc.RegisterServer(ctx, app, mono.RPC()); err != nil {
		return err
	}

	// workers
	w := workers.New(mono.Queue(), queueRepository)
	w.Start()

	return nil
}
