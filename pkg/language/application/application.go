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
		AddUserTerm(ctx *appcontext.AppContext, performerID, termID string, _ dto.AddUserTermRequest) (*dto.AddUserTermResponse, error)
		ChangeTermFavourite(ctx *appcontext.AppContext, performerID, userTermID string, _ dto.ChangeTermFavouriteRequest) (*dto.ChangeTermFavouriteResponse, error)

		StartUserWritingExercise(ctx *appcontext.AppContext, performerID string, req dto.StartUserWritingExerciseRequest) (*dto.StartUserWritingExerciseResponse, error)
		SubmitUserWritingExercise(ctx *appcontext.AppContext, performerID string, req dto.SubmitUserWritingExerciseRequest) (*dto.SubmitUserWritingExerciseResponse, error)
		ModifyUserWritingExercise(ctx *appcontext.AppContext, performerID string, req dto.ModifyUserWritingExerciseRequest) (*dto.ModifyUserWritingExerciseResponse, error)
	}
	Queries interface {
		GetData(_ *appcontext.AppContext, _ dto.GetDataRequest) (*dto.GetDataResponse, error)

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
		command.AddUserTermHandler
		command.ChangeTermFavouriteHandler

		command.StartUserWritingExerciseHandler
		command.SubmitUserWritingExerciseHandler
		command.ModifyUserWritingExerciseHandler
	}
	appQueryHandler struct {
		query.GetDataHandler

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
	userActionHistoryRepository domain.UserActionHistoryRepository,
	writingExerciseRepository domain.WritingExerciseRepository,
	userWritingExerciseRepository domain.UserWritingExerciseRepository,
	userVocabularyExerciseRepository domain.UserVocabularyExerciseRepository,
	openaiRepository domain.OpenAIRepository,
	scraperRepository domain.ScraperRepository,
	queueRepository domain.QueueRepository,
	userHub domain.UserHub,
) *Application {
	return &Application{
		appCommandHandlers: appCommandHandlers{
			AddUserTermHandler:         command.NewAddUserTermHandler(termRepository, userTermRepository, queueRepository, userHub),
			ChangeTermFavouriteHandler: command.NewChangeTermFavouriteHandler(userTermRepository),

			StartUserWritingExerciseHandler:  command.NewStartUserWritingExerciseHandler(writingExerciseRepository, userWritingExerciseRepository),
			SubmitUserWritingExerciseHandler: command.NewSubmitUserWritingExerciseHandler(writingExerciseRepository, userWritingExerciseRepository, userActionHistoryRepository, openaiRepository, userHub),
			ModifyUserWritingExerciseHandler: command.NewModifyUserWritingExerciseHandler(writingExerciseRepository, userWritingExerciseRepository),
		},
		appQueryHandler: appQueryHandler{
			GetDataHandler: query.NewGetDataHandler(),

			SearchTermHandler:      query.NewSearchTermHandler(termRepository, userActionHistoryRepository, openaiRepository, scraperRepository, userHub),
			GetUserTermsHandler:    query.NewGetUserTermsHandler(termRepository, userTermRepository),
			GetFeaturedTermHandler: query.NewGetFeaturedTermHandler(termRepository),

			GetWritingExercisesHandler:     query.NewGetWritingExercisesHandler(writingExerciseRepository),
			GetUserWritingExercisesHandler: query.NewGetUserWritingExercisesHandler(userWritingExerciseRepository),
		},
		appHubHandler: appHubHandler{},
	}
}
