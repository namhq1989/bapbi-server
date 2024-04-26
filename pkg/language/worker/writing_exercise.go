package worker

import (
	"context"

	"github.com/hibiken/asynq"
	"github.com/namhq1989/bapbi-server/internal/utils/appcontext"
	"github.com/namhq1989/bapbi-server/pkg/language/domain"
)

func (w Workers) GenerateWritingExercises(bgCtx context.Context, t *asynq.Task) error {
	var (
		ctx      = appcontext.New(bgCtx)
		language = domain.LanguageEnglish.String()
	)

	ctx.Logger().Info("[worker] process new task", appcontext.Fields{"type": t.Type(), "payload": string(t.Payload())})

	// call OpenAI to generate writing exercises
	// 1 exercise is a combination of writing exercise type and level
	for _, level := range domain.ListLevels {
		for _, exType := range domain.ListWritingExerciseTypes {
			ctx.Logger().Info("call OpenAI to generate writing exercises", appcontext.Fields{"level": level, "exType": exType})
			result, err := w.openaiRepository.WritingExercise(ctx, language, exType.String(), level.String())
			if err != nil {
				ctx.Logger().Error("failed to call OpenAI to generate writing exercises", err, appcontext.Fields{})
				continue
			}

			if result == nil {
				ctx.Logger().Text("no result from OpenAI")
				continue
			}

			ctx.Logger().Text("new writing exercise domain model")
			domainExercise, err := domain.NewWritingExercise(language, exType.String(), level.String(), result.Topic, result.Question, result.Data, result.Vocabulary)
			if err != nil {
				ctx.Logger().Error("failed to create new writing exercise domain model", err, appcontext.Fields{})
				continue
			}

			ctx.Logger().Text("insert writing exercise to database")
			if err = w.writingExerciseRepository.CreateWritingExercise(ctx, *domainExercise); err != nil {
				ctx.Logger().Error("failed to insert writing exercise to database", err, appcontext.Fields{})
				continue
			}
			ctx.Logger().Text("inserted writing exercise to database")
		}
	}

	ctx.Logger().Info("[worker] done task", appcontext.Fields{"type": t.Type()})
	return nil
}
