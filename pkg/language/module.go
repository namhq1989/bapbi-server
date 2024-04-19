package language

import (
	"github.com/namhq1989/bapbi-server/internal/monolith"
	"github.com/namhq1989/bapbi-server/internal/utils/appcontext"
	"github.com/namhq1989/bapbi-server/pkg/language/application"
	"github.com/namhq1989/bapbi-server/pkg/language/infrastructure"
	"github.com/namhq1989/bapbi-server/pkg/language/rest"
)

type Module struct{}

func (Module) Name() string {
	return "LANGUAGE"
}

func (Module) Startup(ctx *appcontext.AppContext, mono monolith.Monolith) error {
	var (
		// infrastructure
		termRepository    = infrastructure.NewTermRepository(mono.Mongo())
		userSearchHistory = infrastructure.NewUserSearchRepository(mono.Mongo())
		openaiRepository  = infrastructure.NewOpenAIRepository(mono.OpenAI())
		scraperRepository = infrastructure.NewScraperRepository(mono.Scraper())

		// application
		app = application.New(termRepository, userSearchHistory, openaiRepository, scraperRepository)
	)

	// rest server
	if err := rest.RegisterServer(ctx, app, mono.Rest(), mono.JWT()); err != nil {
		return err
	}

	return nil
}
