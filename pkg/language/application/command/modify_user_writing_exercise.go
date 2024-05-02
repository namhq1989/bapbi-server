package command

import (
	"github.com/namhq1989/bapbi-server/internal/utils/appcontext"
	apperrors "github.com/namhq1989/bapbi-server/internal/utils/error"
	"github.com/namhq1989/bapbi-server/pkg/language/domain"
	"github.com/namhq1989/bapbi-server/pkg/language/dto"
)

type ModifyUserWritingExerciseHandler struct {
	writingExerciseRepository     domain.WritingExerciseRepository
	userWritingExerciseRepository domain.UserWritingExerciseRepository
}

func NewModifyUserWritingExerciseHandler(
	writingExerciseRepository domain.WritingExerciseRepository,
	userWritingExerciseRepository domain.UserWritingExerciseRepository,
) ModifyUserWritingExerciseHandler {
	return ModifyUserWritingExerciseHandler{
		writingExerciseRepository:     writingExerciseRepository,
		userWritingExerciseRepository: userWritingExerciseRepository,
	}
}

func (h ModifyUserWritingExerciseHandler) ModifyUserWritingExercise(ctx *appcontext.AppContext, performerID string, req dto.ModifyUserWritingExerciseRequest) (*dto.ModifyUserWritingExerciseResponse, error) {
	ctx.Logger().Info("new modify user writing exercise request", appcontext.Fields{"performer": performerID, "exerciseId": req.ExerciseId})

	ctx.Logger().Text("find exercise in db")
	exercise, err := h.writingExerciseRepository.FindByID(ctx, req.ExerciseId)
	if err != nil {
		ctx.Logger().Error("failed to find exercise in db", err, appcontext.Fields{})
		return nil, err
	}
	if exercise == nil {
		ctx.Logger().Error("exercise not found, respond error", nil, appcontext.Fields{})
		return nil, apperrors.Language.ExerciseNotFound
	}

	ctx.Logger().Text("find user exercise in db")
	userExercise, err := h.userWritingExerciseRepository.FindByUserIDAndExerciseID(ctx, performerID, req.ExerciseId)
	if err != nil {
		ctx.Logger().Error("failed to find user exercise in db", err, appcontext.Fields{})
		return nil, err
	}
	if userExercise == nil {
		ctx.Logger().Error("user exercise not found, respond error", nil, appcontext.Fields{})
		return nil, apperrors.Language.ExerciseNotFound
	}
	if userExercise.IsProgressing() {
		ctx.Logger().Text("user exercise is already progressing, respond")
		return &dto.ModifyUserWritingExerciseResponse{}, nil
	}

	ctx.Logger().Text("set user exercise status to progressing")
	userExercise.SetStatus(domain.ExerciseStatusProgressing)

	ctx.Logger().Text("update user exercise in db")
	if err = h.userWritingExerciseRepository.UpdateUserWritingExercise(ctx, *userExercise); err != nil {
		ctx.Logger().Error("failed to update user exercise in db", err, appcontext.Fields{})
		return nil, err
	}

	ctx.Logger().Text("done modify user writing exercise request")
	return &dto.ModifyUserWritingExerciseResponse{}, nil
}
