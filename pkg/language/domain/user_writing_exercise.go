package domain

import (
	"time"

	"github.com/namhq1989/bapbi-server/internal/utils/appcontext"

	"github.com/namhq1989/bapbi-server/internal/database"
)

type UserWritingExerciseRepository interface {
	CreateUserWritingExercise(ctx *appcontext.AppContext, exercise UserWritingExercise) error
	UpdateUserWritingExercise(ctx *appcontext.AppContext, exercise UserWritingExercise) error
	IsExerciseCreated(ctx *appcontext.AppContext, userID, exerciseID string) (bool, error)
}

type UserWritingExercise struct {
	ID          string
	UserID      string
	ExerciseID  string
	Status      ExerciseStatus
	CreatedAt   time.Time
	UpdatedAt   time.Time
	CompletedAt time.Time
}

func NewUserWritingExercise(userID, exerciseID string) (*UserWritingExercise, error) {
	return &UserWritingExercise{
		ID:         database.NewStringID(),
		UserID:     userID,
		ExerciseID: exerciseID,
		Status:     ExerciseStatusProgressing,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}, nil
}

func (d *UserWritingExercise) SetCompleted() {
	d.Status = ExerciseStatusCompleted
	d.CompletedAt = time.Now()
	d.UpdatedAt = time.Now()
}
