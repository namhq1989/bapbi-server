package query

import (
	"github.com/namhq1989/bapbi-server/internal/utils/appcontext"
	"github.com/namhq1989/bapbi-server/internal/utils/pagetoken"
	"github.com/namhq1989/bapbi-server/pkg/language/domain"
	"github.com/namhq1989/bapbi-server/pkg/language/dto"
)

type GetUserTermExercisesHandler struct {
	userTermExerciseRepository domain.UserTermExerciseRepository
}

func NewGetUserTermExercisesHandler(userTermExerciseRepository domain.UserTermExerciseRepository) GetUserTermExercisesHandler {
	return GetUserTermExercisesHandler{
		userTermExerciseRepository: userTermExerciseRepository,
	}
}

func (h GetUserTermExercisesHandler) GetUserTermExercises(ctx *appcontext.AppContext, performerID string, req dto.GetUserTermExerciseRequest) (*dto.GetUserTermExerciseResponse, error) {
	ctx.Logger().Info("new get user term exercises request", appcontext.Fields{"performer": performerID, "language": req.Language, "status": req.Status, "pageToken": req.PageToken})

	ctx.Logger().Text("new filter model")
	filter := domain.NewUserTermExerciseFilter(performerID, req.Language, req.Status, req.PageToken)

	ctx.Logger().Text("find user term exercises in db")
	exercises, err := h.userTermExerciseRepository.FindUserTermExercises(ctx, filter)
	if err != nil {
		ctx.Logger().Error("failed to find user term exercises in db", err, appcontext.Fields{})
		return nil, err
	}

	var (
		totalExercises = int64(len(exercises))
		result         = &dto.GetUserTermExerciseResponse{
			Exercises:     make([]dto.UserTermExercise, totalExercises),
			NextPageToken: "",
		}
	)
	if totalExercises == 0 {
		return result, nil
	}

	ctx.Logger().Info("found data, mapping to dto", appcontext.Fields{"total": totalExercises})
	for i, exercise := range exercises {
		result.Exercises[i] = dto.UserTermExercise{}.FromDomain(exercise)
	}
	ctx.Logger().Text("generate page token")
	if totalExercises >= filter.Limit {
		result.NextPageToken = pagetoken.NewWithTimestamp(exercises[totalExercises-1].CreatedAt)
	}

	return result, nil
}
