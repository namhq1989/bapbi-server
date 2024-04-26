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
		ChangeTermFavourite(ctx *appcontext.AppContext, performerID, userTermID string, _ dto.ChangeTermFavouriteRequest) (*dto.ChangeTermFavouriteResponse, error)

		CreateUserWritingExercise(ctx *appcontext.AppContext, performerID string, req dto.CreateUserWritingExerciseRequest) (*dto.CreateUserWritingExerciseResponse, error)
	}
	Queries interface {
		SearchTerm(ctx *appcontext.AppContext, performerID string, req dto.SearchTermRequest) (*dto.SearchTermResponse, error)
		GetUserTerms(ctx *appcontext.AppContext, performerID string, req dto.GetUserTermsRequest) (*dto.GetUserTermsResponse, error)
		GetFeaturedTerm(ctx *appcontext.AppContext, req dto.GetFeaturedTermRequest) (*dto.GetFeaturedTermResponse, error)

		GetWritingExercises(ctx *appcontext.AppContext, performerID string, req dto.GetWritingExerciseRequest) (*dto.GetWritingExerciseResponse, error)
		GetUserWritingExercises(ctx *appcontext.AppContext, performerID string, req dto.GetUserWritingExerciseRequest) (*dto.GetUserWritingExerciseResponse, error)
	}
	Hubs interface{}
	App  interface {
		Commands
		Queries
		Hubs
	}

	appCommandHandlers struct {
		command.AddTermHandler
		command.ChangeTermFavouriteHandler

		command.CreateUserWritingExerciseHandler
	}
	appQueryHandler struct {
		query.SearchTermHandler
		query.GetUserTermsHandler
		query.GetFeaturedTermHandler

		query.GetWritingExercisesHandler
		query.GetUserWritingExercisesHandler
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
	writingExerciseRepository domain.WritingExerciseRepository,
	userWritingExerciseRepository domain.UserWritingExerciseRepository,
	openaiRepository domain.OpenAIRepository,
	scraperRepository domain.ScraperRepository,
	userHub domain.UserHub,
	authHub domain.AuthHub,
) *Application {
	return &Application{
		appCommandHandlers: appCommandHandlers{
			AddTermHandler:             command.NewAddTermHandler(termRepository, userTermRepository, userHub),
			ChangeTermFavouriteHandler: command.NewChangeTermFavouriteHandler(userTermRepository),

			CreateUserWritingExerciseHandler: command.NewCreateUserWritingExerciseHandler(writingExerciseRepository, userWritingExerciseRepository),
		},
		appQueryHandler: appQueryHandler{
			SearchTermHandler:      query.NewSearchTermHandler(termRepository, userSearchHistoryRepository, openaiRepository, scraperRepository, userHub),
			GetUserTermsHandler:    query.NewGetUserTermsHandler(termRepository, userTermRepository),
			GetFeaturedTermHandler: query.NewGetFeaturedTermHandler(termRepository),

			GetWritingExercisesHandler:     query.NewGetWritingExercisesHandler(writingExerciseRepository),
			GetUserWritingExercisesHandler: query.NewGetUserWritingExercisesHandler(userWritingExerciseRepository),
		},
		appHubHandler: appHubHandler{},
	}
}
