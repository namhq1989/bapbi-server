package command

import (
	"time"

	"github.com/namhq1989/bapbi-server/internal/utils/appcontext"
	apperrors "github.com/namhq1989/bapbi-server/internal/utils/error"
	"github.com/namhq1989/bapbi-server/internal/utils/manipulation"
	"github.com/namhq1989/bapbi-server/pkg/language/domain"
	"github.com/namhq1989/bapbi-server/pkg/language/dto"
)

type SubmitUserWritingExerciseHandler struct {
	writingExerciseRepository     domain.WritingExerciseRepository
	userWritingExerciseRepository domain.UserWritingExerciseRepository
	userActionHistoryRepository   domain.UserActionHistoryRepository
	openaiRepository              domain.OpenAIRepository
	userHub                       domain.UserHub
}

func NewSubmitUserWritingExerciseHandler(
	writingExerciseRepository domain.WritingExerciseRepository,
	userWritingExerciseRepository domain.UserWritingExerciseRepository,
	userActionHistoryRepository domain.UserActionHistoryRepository,
	openaiRepository domain.OpenAIRepository,
	userHub domain.UserHub,
) SubmitUserWritingExerciseHandler {
	return SubmitUserWritingExerciseHandler{
		writingExerciseRepository:     writingExerciseRepository,
		userWritingExerciseRepository: userWritingExerciseRepository,
		userActionHistoryRepository:   userActionHistoryRepository,
		openaiRepository:              openaiRepository,
		userHub:                       userHub,
	}
}

func (h SubmitUserWritingExerciseHandler) SubmitUserWritingExercise(ctx *appcontext.AppContext, performerID string, req dto.SubmitUserWritingExerciseRequest) (*dto.SubmitUserWritingExerciseResponse, error) {
	ctx.Logger().Info("new submit user writing exercise request", appcontext.Fields{"performer": performerID, "exerciseId": req.ExerciseID})

	ctx.Logger().Text("get user's subscription plan")
	plan, err := h.userHub.GetUserPlan(ctx, performerID)
	if err != nil {
		ctx.Logger().Error("failed to get user's subscription plan", err, appcontext.Fields{})
		return nil, err
	}

	ctx.Logger().Text("count total actions today")
	var (
		start = manipulation.StartOfToday()
		end   = time.Now()
	)
	totalActions, err := h.userActionHistoryRepository.CountTotalActionsByTimeRange(ctx, performerID, start, end)
	if err != nil {
		ctx.Logger().Error("failed to count total actions today", err, appcontext.Fields{})
		return nil, err
	}
	if isExceeded := plan.IsExceededActionLimitation(totalActions); isExceeded {
		ctx.Logger().Error("exceeded search term limitation", nil, appcontext.Fields{"plan": plan.String(), "actions": totalActions})
		return nil, apperrors.User.ExceededPlanLimitation
	}

	ctx.Logger().Info("still available to submit writing exercise, find exercise in db", appcontext.Fields{"actions": totalActions})
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
	_ = h.insertUserActionHistory(ctx, performerID, exercise.ID)
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
	return &dto.SubmitUserWritingExerciseResponse{
		ID: exercise.ID,
	}, nil
}

func (h SubmitUserWritingExerciseHandler) insertUserActionHistory(ctx *appcontext.AppContext, performerID, exerciseID string) error {
	ctx.Logger().Text("new user action history")

	action, err := domain.NewUserActionHistory(performerID, domain.UserActionTypeSubmitWritingExercise.String())
	if err != nil {
		ctx.Logger().Error("failed to create new user action history", err, appcontext.Fields{})
		return err
	}

	// set data
	action.SetData(domain.UserActionHistoryData{
		ExerciseID: exerciseID,
	})

	ctx.Logger().Text("insert user action history to database")
	if err = h.userActionHistoryRepository.CreateUserActionHistory(ctx, *action); err != nil {
		ctx.Logger().Error("failed to insert user action history to database", err, appcontext.Fields{})
		return err
	}

	return nil
}
