package command

import (
	"github.com/namhq1989/bapbi-server/internal/utils/appcontext"
	apperrors "github.com/namhq1989/bapbi-server/internal/utils/error"
	"github.com/namhq1989/bapbi-server/pkg/language/domain"
	"github.com/namhq1989/bapbi-server/pkg/language/dto"
)

type ModifyUserTermExerciseHandler struct {
	UserTermExerciseRepository domain.UserTermExerciseRepository
}

func NewModifyUserTermExerciseHandler(
	UserTermExerciseRepository domain.UserTermExerciseRepository,
) ModifyUserTermExerciseHandler {
	return ModifyUserTermExerciseHandler{
		UserTermExerciseRepository: UserTermExerciseRepository,
	}
}

func (h ModifyUserTermExerciseHandler) ModifyUserTermExercise(ctx *appcontext.AppContext, performerID string, req dto.ModifyUserTermExerciseRequest) (*dto.ModifyUserTermExerciseResponse, error) {
	ctx.Logger().Info("new modify user vocabulary exercise request", appcontext.Fields{"performer": performerID, "exerciseId": req.ExerciseId})

	ctx.Logger().Text("find exercise in db")
	exercise, err := h.UserTermExerciseRepository.FindByExerciseID(ctx, req.ExerciseId)
	if err != nil {
		ctx.Logger().Error("failed to find exercise in db", err, appcontext.Fields{})
		return nil, err
	}
	if exercise == nil {
		ctx.Logger().Error("exercise not found, respond error", nil, appcontext.Fields{})
		return nil, apperrors.Language.ExerciseNotFound
	}

	if exercise.IsProgressing() {
		ctx.Logger().Text("user exercise is already progressing, respond")
		return &dto.ModifyUserTermExerciseResponse{}, nil
	}

	ctx.Logger().Text("set user exercise status to progressing")
	exercise.SetProgressing()

	ctx.Logger().Text("update user exercise in db")
	if err = h.UserTermExerciseRepository.UpdateUserTermExercise(ctx, *exercise); err != nil {
		ctx.Logger().Error("failed to update user exercise in db", err, appcontext.Fields{})
		return nil, err
	}

	ctx.Logger().Text("done modify user vocabulary exercise request")
	return &dto.ModifyUserTermExerciseResponse{}, nil
}
