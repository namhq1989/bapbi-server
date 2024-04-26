package mapping

import "github.com/namhq1989/bapbi-server/pkg/language/infrastructure/model"

type QueryUserWritingExercise struct {
	model.WritingExercise
	IsCompleted bool `bson:"isCompleted"`
}
