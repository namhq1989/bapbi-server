package dto

import "github.com/namhq1989/bapbi-server/internal/utils/httprespond"

type SubmitUserWritingExerciseRequest struct {
	ExerciseID string `json:"exerciseId" validate:"required" message:"language_invalid_exercise_id"`
	Content    string `json:"content" validate:"required" message:"language_invalid_writing_exercise_data"`
}

type SubmitUserWritingExerciseResponse struct {
	Status           string                    `json:"status"`
	CompletedAt      *httprespond.TimeResponse `json:"completedAt"`
	IsTopicRelevance bool                      `json:"isTopicRelevance"`
	Score            int                       `json:"score"`
	Improvement      []string                  `json:"improvement"`
	Comment          string                    `json:"comment"`
}
