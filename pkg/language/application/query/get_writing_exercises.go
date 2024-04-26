package query

import (
	"github.com/namhq1989/bapbi-server/internal/utils/appcontext"
	"github.com/namhq1989/bapbi-server/internal/utils/pagetoken"
	"github.com/namhq1989/bapbi-server/pkg/language/domain"
	"github.com/namhq1989/bapbi-server/pkg/language/dto"
)

type GetWritingExercisesHandler struct {
	writingExerciseRepository domain.WritingExerciseRepository
}

func NewGetWritingExercisesHandler(writingExerciseRepository domain.WritingExerciseRepository) GetWritingExercisesHandler {
	return GetWritingExercisesHandler{
		writingExerciseRepository: writingExerciseRepository,
	}
}

func (h GetWritingExercisesHandler) GetWritingExercises(ctx *appcontext.AppContext, performerID string, req dto.GetWritingExerciseRequest) (*dto.GetWritingExerciseResponse, error) {
	ctx.Logger().Info("new get writing exercises request", appcontext.Fields{"performer": performerID, "language": req.Language, "level": req.Level, "status": req.Status, "pageToken": req.PageToken})

	ctx.Logger().Text("new filter model")
	filter := domain.NewWritingExerciseFilter(performerID, req.Language, req.Level, req.Status, req.PageToken)

	ctx.Logger().Text("find writing exercises in db")
	exercises, err := h.writingExerciseRepository.FindWritingExercises(ctx, filter)
	if err != nil {
		ctx.Logger().Error("failed to find writing exercises in db", err, appcontext.Fields{})
		return nil, err
	}

	var (
		totalExercises = int64(len(exercises))
		result         = &dto.GetWritingExerciseResponse{
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
