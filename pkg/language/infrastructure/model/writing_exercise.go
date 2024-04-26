package model

import (
	"time"

	apperrors "github.com/namhq1989/bapbi-server/internal/utils/error"
	"github.com/namhq1989/bapbi-server/pkg/language/domain"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type WritingExercise struct {
	ID         primitive.ObjectID `bson:"_id"`
	Language   string             `bson:"language"`
	Type       string             `bson:"type"`
	Level      string             `bson:"level"`
	Topic      string             `bson:"topic"`
	Question   string             `bson:"question"`
	Vocabulary []string           `bson:"vocabulary"`
	MinWords   int                `bson:"minWords"`
	Data       string             `bson:"data"`
	CreatedAt  time.Time          `bson:"createdAt"`
}

func (m WritingExercise) ToDomain() domain.WritingExercise {
	return domain.WritingExercise{
		ID:         m.ID.Hex(),
		Language:   domain.ToLanguage(m.Language),
		Type:       domain.ToWritingExerciseType(m.Type),
		Level:      domain.ToLevel(m.Level),
		Topic:      m.Topic,
		Question:   m.Question,
		Vocabulary: m.Vocabulary,
		MinWords:   m.MinWords,
		Data:       m.Data,
		CreatedAt:  m.CreatedAt,
	}
}

func (WritingExercise) FromDomain(exercise domain.WritingExercise) (*WritingExercise, error) {
	id, err := primitive.ObjectIDFromHex(exercise.ID)
	if err != nil {
		return nil, apperrors.Common.InvalidID
	}

	return &WritingExercise{
		ID:         id,
		Language:   exercise.Language.String(),
		Type:       exercise.Type.String(),
		Level:      exercise.Level.String(),
		Topic:      exercise.Topic,
		Question:   exercise.Question,
		Vocabulary: exercise.Vocabulary,
		MinWords:   exercise.MinWords,
		Data:       exercise.Data,
		CreatedAt:  exercise.CreatedAt,
	}, nil
}
