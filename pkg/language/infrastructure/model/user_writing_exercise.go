package model

import (
	"time"

	apperrors "github.com/namhq1989/bapbi-server/internal/utils/error"
	"github.com/namhq1989/bapbi-server/pkg/language/domain"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserWritingExercise struct {
	ID          primitive.ObjectID `bson:"_id"`
	UserID      primitive.ObjectID `bson:"userId"`
	ExerciseID  primitive.ObjectID `bson:"exerciseId"`
	Language    string             `bson:"language"`
	Status      string             `bson:"status"`
	CreatedAt   time.Time          `bson:"createdAt"`
	UpdatedAt   time.Time          `bson:"updatedAt"`
	CompletedAt time.Time          `bson:"completedAt"`
}

func (u UserWritingExercise) ToDomain() domain.UserWritingExercise {
	return domain.UserWritingExercise{
		ID:          u.ID.Hex(),
		UserID:      u.UserID.Hex(),
		ExerciseID:  u.ExerciseID.Hex(),
		Status:      domain.ToExerciseStatus(u.Status),
		Language:    domain.ToLanguage(u.Language),
		CreatedAt:   u.CreatedAt,
		UpdatedAt:   u.UpdatedAt,
		CompletedAt: u.CompletedAt,
	}
}

func (u UserWritingExercise) FromDomain(exercise domain.UserWritingExercise) (*UserWritingExercise, error) {
	id, err := primitive.ObjectIDFromHex(exercise.ID)
	if err != nil {
		return nil, apperrors.Common.InvalidID
	}

	uid, err := primitive.ObjectIDFromHex(exercise.UserID)
	if err != nil {
		return nil, apperrors.User.InvalidUserID
	}

	eid, err := primitive.ObjectIDFromHex(exercise.ExerciseID)
	if err != nil {
		return nil, apperrors.Common.InvalidID
	}

	return &UserWritingExercise{
		ID:          id,
		UserID:      uid,
		ExerciseID:  eid,
		Status:      exercise.Status.String(),
		Language:    exercise.Language.String(),
		CreatedAt:   exercise.CreatedAt,
		UpdatedAt:   exercise.UpdatedAt,
		CompletedAt: exercise.CompletedAt,
	}, nil
}
