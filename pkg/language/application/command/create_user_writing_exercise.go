package command

import (
	"github.com/namhq1989/bapbi-server/internal/utils/appcontext"
	apperrors "github.com/namhq1989/bapbi-server/internal/utils/error"
	"github.com/namhq1989/bapbi-server/pkg/language/domain"
	"github.com/namhq1989/bapbi-server/pkg/language/dto"
)

type CreateUserWritingExerciseHandler struct {
	writingExerciseRepository     domain.WritingExerciseRepository
	userWritingExerciseRepository domain.UserWritingExerciseRepository
}

func NewCreateUserWritingExerciseHandler(
	writingExerciseRepository domain.WritingExerciseRepository,
	userWritingExerciseRepository domain.UserWritingExerciseRepository,
) CreateUserWritingExerciseHandler {
	return CreateUserWritingExerciseHandler{
		writingExerciseRepository:     writingExerciseRepository,
		userWritingExerciseRepository: userWritingExerciseRepository,
	}
}

func (h CreateUserWritingExerciseHandler) CreateUserWritingExercise(ctx *appcontext.AppContext, performerID string, req dto.CreateUserWritingExerciseRequest) (*dto.CreateUserWritingExerciseResponse, error) {
	ctx.Logger().Info("new create user writing exercise request", appcontext.Fields{"performer": performerID, "exerciseId": req.ExerciseID})

	ctx.Logger().Text("check if exercise already created")
	isCreated, err := h.userWritingExerciseRepository.IsExerciseCreated(ctx, performerID, req.ExerciseID)
	if err != nil {
		ctx.Logger().Error("failed to check if exercise already created", err, appcontext.Fields{})
		return nil, err
	}
	if isCreated {
		ctx.Logger().Error("exercise already created, respond error", nil, appcontext.Fields{})
		return nil, apperrors.Language.UserExerciseExisted
	}

	ctx.Logger().Text("find exercise in db")
	exercise, err := h.writingExerciseRepository.FindByID(ctx, req.ExerciseID)
	if err != nil {
		ctx.Logger().Error("failed to find exercise in db", err, appcontext.Fields{})
		return nil, err
	}
	if exercise == nil {
		ctx.Logger().Error("exercise not found, respond error", nil, appcontext.Fields{})
		return nil, apperrors.Language.ExerciseNotFound
	}

	ctx.Logger().Text("create user writing exercise domain model")
	domainExercise, err := domain.NewUserWritingExercise(performerID, exercise.ID)
	if err != nil {
		ctx.Logger().Error("failed to create user writing exercise domain model", err, appcontext.Fields{})
		return nil, err
	}

	ctx.Logger().Text("insert user writing exercise to database")
	if err = h.userWritingExerciseRepository.CreateUserWritingExercise(ctx, *domainExercise); err != nil {
		ctx.Logger().Error("failed to insert user writing exercise to database", err, appcontext.Fields{})
		return nil, err
	}

	ctx.Logger().Text("done create user writing exercise request")
	return &dto.CreateUserWritingExerciseResponse{
		ID: domainExercise.ID,
	}, nil
}
