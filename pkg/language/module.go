package language

import (
	"github.com/namhq1989/bapbi-server/internal/grpcclient"
	"github.com/namhq1989/bapbi-server/internal/monolith"
	"github.com/namhq1989/bapbi-server/internal/utils/appcontext"
	"github.com/namhq1989/bapbi-server/pkg/language/application"
	"github.com/namhq1989/bapbi-server/pkg/language/infrastructure"
	"github.com/namhq1989/bapbi-server/pkg/language/rest"
	"github.com/namhq1989/bapbi-server/pkg/language/worker"
)

type Module struct{}

func (Module) Name() string {
	return "LANGUAGE"
}

func (Module) Startup(ctx *appcontext.AppContext, mono monolith.Monolith) error {
	authGRPCClient, err := grpcclient.NewAuthClient(ctx, mono.Config().GRPCPort)
	if err != nil {
		return err
	}

	userGRPCClient, err := grpcclient.NewUserClient(ctx, mono.Config().GRPCPort)
	if err != nil {
		return err
	}

	var (
		// infrastructure
		termRepository                = infrastructure.NewTermRepository(mono.Mongo())
		userTermRepository            = infrastructure.NewUserTermRepository(mono.Mongo())
		userSearchHistoryRepository   = infrastructure.NewUserSearchHistoryRepository(mono.Mongo())
		writingExerciseRepository     = infrastructure.NewWritingExerciseRepository(mono.Mongo())
		userWritingExerciseRepository = infrastructure.NewUserWritingExerciseRepository(mono.Mongo())
		openaiRepository              = infrastructure.NewOpenAIRepository(mono.OpenAI())
		scraperRepository             = infrastructure.NewScraperRepository(mono.Scraper())

		// hub
		userHub = infrastructure.NewUserHub(userGRPCClient)
		authHub = infrastructure.NewAuthHub(authGRPCClient)

		// application
		app = application.New(
			termRepository,
			userTermRepository,
			userSearchHistoryRepository,
			writingExerciseRepository,
			userWritingExerciseRepository,
			openaiRepository,
			scraperRepository,
			userHub,
			authHub,
		)
	)

	// rest server
	if err = rest.RegisterServer(ctx, app, mono.Rest(), mono.JWT()); err != nil {
		return err
	}

	// worker
	w := worker.New(mono.Queue(), termRepository, writingExerciseRepository, openaiRepository, scraperRepository)
	w.Start()

	return nil
}
