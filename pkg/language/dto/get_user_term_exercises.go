package dto

import (
	"time"

	"github.com/namhq1989/bapbi-server/pkg/language/domain"
)

type GetUserTermExerciseRequest struct {
	PageToken string `query:"pageToken"`
	Language  string `query:"language" validate:"required" message:"language_invalid_language"`
	Status    string `query:"status"`
}

type GetUserTermExerciseResponse struct {
	Exercises     []UserTermExercise `json:"exercises"`
	NextPageToken string             `json:"nextPageToken"`
}

type UserTermExercise struct {
	ID          string                      `json:"id"`
	TermID      string                      `json:"termId"`
	Term        string                      `json:"term"`
	Language    string                      `json:"language"`
	Tense       string                      `json:"tense"`
	Content     string                      `json:"content"`
	Status      string                      `json:"status"`
	Assessment  *UserTermExerciseAssessment `json:"assessment"`
	CreatedAt   time.Time                   `json:"createdAt"`
	UpdatedAt   time.Time                   `json:"updatedAt"`
	CompletedAt time.Time                   `json:"completedAt"`
}

func (d UserTermExercise) FromDomain(exercise domain.UserTermExercise) UserTermExercise {
	return UserTermExercise{
		ID:          exercise.ID,
		TermID:      exercise.TermID,
		Term:        exercise.Term,
		Language:    exercise.Language.String(),
		Tense:       exercise.Tense.String(),
		Content:     exercise.Content,
		Status:      exercise.Status.String(),
		Assessment:  UserTermExerciseAssessment{}.FromDomain(exercise.Assessment),
		CreatedAt:   exercise.CreatedAt,
		UpdatedAt:   exercise.UpdatedAt,
		CompletedAt: exercise.CompletedAt,
	}
}
