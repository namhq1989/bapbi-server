package query

import (
	"github.com/namhq1989/bapbi-server/internal/utils/appcontext"
	"github.com/namhq1989/bapbi-server/internal/utils/pagetoken"
	"github.com/namhq1989/bapbi-server/pkg/language/domain"
	"github.com/namhq1989/bapbi-server/pkg/language/dto"
)

type GetUserWritingExercisesHandler struct {
	userWritingExerciseRepository domain.UserWritingExerciseRepository
}

func NewGetUserWritingExercisesHandler(userWritingExerciseRepository domain.UserWritingExerciseRepository) GetUserWritingExercisesHandler {
	return GetUserWritingExercisesHandler{
		userWritingExerciseRepository: userWritingExerciseRepository,
	}
}

func (h GetUserWritingExercisesHandler) GetUserWritingExercises(ctx *appcontext.AppContext, performerID string, req dto.GetUserWritingExerciseRequest) (*dto.GetUserWritingExerciseResponse, error) {
	ctx.Logger().Info("new get user writing exercises request", appcontext.Fields{"performer": performerID, "language": req.Language, "status": req.Status, "pageToken": req.PageToken})

	ctx.Logger().Text("new filter model")
	filter := domain.NewUserWritingExerciseFilter(performerID, req.Language, req.Status, req.PageToken)

	ctx.Logger().Text("find user writing exercises in db")
	exercises, err := h.userWritingExerciseRepository.FindUserWritingExercises(ctx, filter)
	if err != nil {
		ctx.Logger().Error("failed to find user writing exercises in db", err, appcontext.Fields{})
		return nil, err
	}

	var (
		totalExercises = int64(len(exercises))
		result         = &dto.GetUserWritingExerciseResponse{
			Exercises:     make([]dto.WritingExercise, totalExercises),
			NextPageToken: "",
		}
	)
	if totalExercises == 0 {
		return result, nil
	}

	ctx.Logger().Info("found data, mapping to dto", appcontext.Fields{"total": totalExercises})
	for i, exercise := range exercises {
		result.Exercises[i] = dto.WritingExercise{}.FromDomain(exercise)
	}
	ctx.Logger().Text("generate page token")
	if totalExercises >= filter.Limit {
		result.NextPageToken = pagetoken.NewWithTimestamp(exercises[totalExercises-1].CreatedAt)
	}

	return result, nil
}
