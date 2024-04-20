package application

import (
	"github.com/namhq1989/bapbi-server/internal/utils/appcontext"
	"github.com/namhq1989/bapbi-server/pkg/language/application/command"
	"github.com/namhq1989/bapbi-server/pkg/language/application/query"
	"github.com/namhq1989/bapbi-server/pkg/language/domain"
	"github.com/namhq1989/bapbi-server/pkg/language/dto"
)

type (
	Commands interface {
		AddTerm(ctx *appcontext.AppContext, performerID, termID string, _ dto.AddTermRequest) (*dto.AddTermResponse, error)
	}
	Queries interface {
		SearchTerm(ctx *appcontext.AppContext, performerID string, req dto.SearchTermRequest) (*dto.SearchTermResponse, error)
	}
	Hubs interface{}
	App  interface {
		Commands
		Queries
		Hubs
	}

	appCommandHandlers struct {
		command.AddTermHandler
	}
	appQueryHandler struct {
		query.SearchTermHandler
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
	termRepository domain.TermRepository,
	userTermRepository domain.UserTermRepository,
	userSearchHistoryRepository domain.UserSearchHistoryRepository,
	openaiRepository domain.OpenAIRepository,
	scraperRepository domain.ScraperRepository,
	userHub domain.UserHub,
) *Application {
	return &Application{
		appCommandHandlers: appCommandHandlers{
			AddTermHandler: command.NewAddTermHandler(termRepository, userTermRepository, userHub),
		},
		appQueryHandler: appQueryHandler{
			SearchTermHandler: query.NewSearchTermHandler(termRepository, userSearchHistoryRepository, openaiRepository, scraperRepository),
		},
		appHubHandler: appHubHandler{},
	}
}
