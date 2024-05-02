package command

import (
	"github.com/namhq1989/bapbi-server/internal/utils/httprespond"

	"github.com/namhq1989/bapbi-server/internal/utils/appcontext"
	apperrors "github.com/namhq1989/bapbi-server/internal/utils/error"
	"github.com/namhq1989/bapbi-server/pkg/language/domain"
	"github.com/namhq1989/bapbi-server/pkg/language/dto"
)

type SubmitUserWritingExerciseHandler struct {
	writingExerciseRepository     domain.WritingExerciseRepository
	userWritingExerciseRepository domain.UserWritingExerciseRepository
	openaiRepository              domain.OpenAIRepository
	userHub                       domain.UserHub
	languageService               domain.LanguageService
}

func NewSubmitUserWritingExerciseHandler(
	writingExerciseRepository domain.WritingExerciseRepository,
	userWritingExerciseRepository domain.UserWritingExerciseRepository,
	openaiRepository domain.OpenAIRepository,
	userHub domain.UserHub,
	languageService domain.LanguageService,
) SubmitUserWritingExerciseHandler {
	return SubmitUserWritingExerciseHandler{
		writingExerciseRepository:     writingExerciseRepository,
		userWritingExerciseRepository: userWritingExerciseRepository,
		openaiRepository:              openaiRepository,
		userHub:                       userHub,
		languageService:               languageService,
	}
}

func (h SubmitUserWritingExerciseHandler) SubmitUserWritingExercise(ctx *appcontext.AppContext, performerID string, req dto.SubmitUserWritingExerciseRequest) (*dto.SubmitUserWritingExerciseResponse, error) {
	ctx.Logger().Info("new submit user writing exercise request", appcontext.Fields{"performer": performerID, "exerciseId": req.ExerciseID})

	ctx.Logger().Text("check today actions limitation")
	isExceeded, err := h.languageService.IsExceededActionLimitation(ctx, performerID)
	if err != nil {
		ctx.Logger().Error("failed to check today actions limitation", err, appcontext.Fields{})
		return nil, err
	}
	if isExceeded {
		ctx.Logger().Error("exceeded action limitation", nil, appcontext.Fields{})
		return nil, apperrors.User.ExceededPlanLimitation
	}

	exercise, err := h.writingExerciseRepository.FindByID(ctx, req.ExerciseID)
	if err != nil {
		ctx.Logger().Error("failed to find exercise in db", err, appcontext.Fields{})
		return nil, err
	}
	if exercise == nil {
		ctx.Logger().Error("exercise not found, respond error", nil, appcontext.Fields{})
		return nil, apperrors.Language.ExerciseNotFound
	}

	ctx.Logger().Text("find user writing exercise in db")
	userWritingExercise, err := h.userWritingExerciseRepository.FindByUserIDAndExerciseID(ctx, performerID, exercise.ID)
	if err != nil {
		ctx.Logger().Error("failed to find user writing exercise in db", err, appcontext.Fields{})
		return nil, err
	}
	if userWritingExercise == nil {
		ctx.Logger().Error("user writing exercise not found, respond error", nil, appcontext.Fields{})
		return nil, apperrors.Language.ExerciseNotFound
	}
	if userWritingExercise.IsCompleted() {
		ctx.Logger().Error("user writing exercise already completed, respond error", nil, appcontext.Fields{})
		return nil, apperrors.Language.ExerciseAlreadyCompleted
	}

	ctx.Logger().Text("set user writing exercise content")
	if err = userWritingExercise.SetContent(req.Content, exercise.MinWords); err != nil {
		ctx.Logger().Error("failed to set user writing exercise content", err, appcontext.Fields{"content": req.Content, "minWords": exercise.MinWords})
		return nil, err
	}

	ctx.Logger().Text("call OpenAI's api to assess the exercise's content")
	assessment, err := h.openaiRepository.AssessWritingExercise(ctx, exercise.Language.String(), exercise.Topic, exercise.Level.String(), req.Content)
	if err != nil {
		ctx.Logger().Error("failed to call OpenAI's api to assess the exercise's content", err, appcontext.Fields{})
		return nil, err
	}
	_ = h.languageService.PersistUserActionHistory(ctx, performerID, domain.UserActionTypeSubmitWritingExercise.String(), domain.UserActionHistoryData{ExerciseID: exercise.ID})
	if assessment == nil {
		ctx.Logger().Error("failed to assess the exercise's content, respond error", nil, appcontext.Fields{})
		return nil, apperrors.Common.BadRequest
	}

	ctx.Logger().Text("set assessment")
	userWritingExercise.SetAssessment(assessment.IsTopicRelevance, assessment.Score, assessment.Improvement, assessment.Comment)

	ctx.Logger().Text("set user writing exercise status to completed")
	userWritingExercise.SetCompleted()

	ctx.Logger().Text("update user writing exercise in db")
	err = h.userWritingExerciseRepository.UpdateUserWritingExercise(ctx, *userWritingExercise)
	if err != nil {
		ctx.Logger().Error("failed to update user writing exercise in db", err, appcontext.Fields{})
		return nil, err
	}

	ctx.Logger().Text("done submit user writing exercise request")
	dtoAssessment := dto.UserWritingExerciseAssessment{}.FromDomain(*userWritingExercise.Assessment)
	return &dto.SubmitUserWritingExerciseResponse{
		Status:      userWritingExercise.Status.String(),
		CompletedAt: httprespond.NewTimeResponse(userWritingExercise.CompletedAt),
		Assessment:  dtoAssessment,
	}, nil
}
