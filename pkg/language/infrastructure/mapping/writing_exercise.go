package mapping

import (
	"github.com/namhq1989/bapbi-server/pkg/language/domain"
	"github.com/namhq1989/bapbi-server/pkg/language/infrastructure/model"
)

type WritingExerciseDatabaseQuery struct {
	model.WritingExercise `bson:",inline"`
	Status                string `bson:"status"`
}

func (q WritingExerciseDatabaseQuery) ToDomain() domain.WritingExerciseDatabaseQuery {
	return domain.WritingExerciseDatabaseQuery{
		WritingExercise: q.WritingExercise.ToDomain(),
		Status:          domain.ToExerciseStatus(q.Status),
	}
}
