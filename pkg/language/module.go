package language

import (
	"github.com/namhq1989/bapbi-server/internal/grpcclient"
	"github.com/namhq1989/bapbi-server/internal/monolith"
	"github.com/namhq1989/bapbi-server/internal/utils/appcontext"
	"github.com/namhq1989/bapbi-server/pkg/language/application"
	"github.com/namhq1989/bapbi-server/pkg/language/infrastructure"
	"github.com/namhq1989/bapbi-server/pkg/language/rest"
	"github.com/namhq1989/bapbi-server/pkg/language/service"
	"github.com/namhq1989/bapbi-server/pkg/language/worker"
)

type Module struct{}

func (Module) Name() string {
	return "LANGUAGE"
}

func (Module) Startup(ctx *appcontext.AppContext, mono monolith.Monolith) error {
	userGRPCClient, err := grpcclient.NewUserClient(ctx, mono.Config().GRPCPort)
	if err != nil {
		return err
	}

	var (
		// infrastructure
		termRepository                   = infrastructure.NewTermRepository(mono.Mongo())
		userTermRepository               = infrastructure.NewUserTermRepository(mono.Mongo())
		userActionHistoryRepository      = infrastructure.NewUserActionHistoryRepository(mono.Mongo())
		writingExerciseRepository        = infrastructure.NewWritingExerciseRepository(mono.Mongo())
		userWritingExerciseRepository    = infrastructure.NewUserWritingExerciseRepository(mono.Mongo())
		userVocabularyExerciseRepository = infrastructure.NewUserVocabularyExerciseRepository(mono.Mongo())

		// 3rd
		queueRepository   = infrastructure.NewQueueRepository(mono.Queue())
		openaiRepository  = infrastructure.NewOpenAIRepository(mono.OpenAI())
		scraperRepository = infrastructure.NewScraperRepository(mono.Scraper())

		// hub
		userHub = infrastructure.NewUserHub(userGRPCClient)

		// service
		languageService = service.NewLanguageService(userTermRepository, userActionHistoryRepository, userHub)

		// application
		app = application.New(
			termRepository,
			userTermRepository,
			writingExerciseRepository,
			userWritingExerciseRepository,
			userVocabularyExerciseRepository,
			openaiRepository,
			scraperRepository,
			queueRepository,
			userHub,
			languageService,
		)
	)

	// rest server
	if err = rest.RegisterServer(ctx, app, mono.Rest(), mono.JWT()); err != nil {
		return err
	}

	// worker
	w := worker.New(
		mono.Queue(),
		termRepository,
		writingExerciseRepository,
		userVocabularyExerciseRepository,
		openaiRepository,
		scraperRepository,
	)
	w.Start()

	return nil
}
