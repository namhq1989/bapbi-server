package model

import (
	"time"

	apperrors "github.com/namhq1989/bapbi-server/internal/utils/error"
	"github.com/namhq1989/bapbi-server/pkg/language/domain"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserVocabularyExerciseAssessment struct {
	IsRelevance bool     `bson:"isRelevance"`
	Score       int      `bson:"score"`
	Improvement []string `bson:"improvement"`
	Comment     string   `bson:"comment"`
}

type UserVocabularyExercise struct {
	ID          primitive.ObjectID                `bson:"_id"`
	UserID      primitive.ObjectID                `bson:"userId"`
	TermID      primitive.ObjectID                `bson:"termId"`
	Term        string                            `bson:"term"`
	Language    string                            `bson:"language"`
	Tense       string                            `bson:"tense"`
	Content     string                            `bson:"content"`
	Status      string                            `bson:"status"`
	Assessment  *UserVocabularyExerciseAssessment `bson:"assessment"`
	CreatedAt   time.Time                         `bson:"createdAt"`
	UpdatedAt   time.Time                         `bson:"updatedAt"`
	CompletedAt time.Time                         `bson:"completedAt"`
}

func (u UserVocabularyExercise) ToDomain() domain.UserVocabularyExercise {
	exercise := domain.UserVocabularyExercise{
		ID:          u.ID.Hex(),
		UserID:      u.UserID.Hex(),
		TermID:      u.TermID.Hex(),
		Term:        u.Term,
		Language:    domain.ToLanguage(u.Language),
		Tense:       domain.ToGrammarTenseCode(u.Tense),
		Content:     u.Content,
		Status:      domain.ToExerciseStatus(u.Status),
		Assessment:  nil,
		CreatedAt:   u.CreatedAt,
		UpdatedAt:   u.UpdatedAt,
		CompletedAt: u.CompletedAt,
	}

	if u.Assessment != nil {
		exercise.Assessment = &domain.UserVocabularyExerciseAssessment{
			IsRelevance: u.Assessment.IsRelevance,
			Score:       u.Assessment.Score,
			Improvement: u.Assessment.Improvement,
			Comment:     u.Assessment.Comment,
		}
	}
	return exercise
}

func (u UserVocabularyExercise) FromDomain(exercise domain.UserVocabularyExercise) (*UserVocabularyExercise, error) {
	id, err := primitive.ObjectIDFromHex(exercise.ID)
	if err != nil {
		return nil, apperrors.Common.InvalidID
	}

	uid, err := primitive.ObjectIDFromHex(exercise.UserID)
	if err != nil {
		return nil, apperrors.User.InvalidUserID
	}

	tid, err := primitive.ObjectIDFromHex(exercise.TermID)
	if err != nil {
		return nil, apperrors.Common.InvalidID
	}

	var assessment *UserVocabularyExerciseAssessment
	if exercise.Assessment != nil {
		assessment = &UserVocabularyExerciseAssessment{
			IsRelevance: exercise.Assessment.IsRelevance,
			Score:       exercise.Assessment.Score,
			Improvement: exercise.Assessment.Improvement,
			Comment:     exercise.Assessment.Comment,
		}
	}

	return &UserVocabularyExercise{
		ID:          id,
		UserID:      uid,
		TermID:      tid,
		Term:        exercise.Term,
		Language:    exercise.Language.String(),
		Tense:       exercise.Tense.String(),
		Content:     exercise.Content,
		Status:      exercise.Status.String(),
		Assessment:  assessment,
		CreatedAt:   exercise.CreatedAt,
		UpdatedAt:   exercise.UpdatedAt,
		CompletedAt: exercise.CompletedAt,
	}, nil
}
