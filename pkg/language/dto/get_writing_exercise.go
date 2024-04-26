package dto

import "github.com/namhq1989/bapbi-server/pkg/language/domain"

type GetWritingExerciseRequest struct {
	PageToken string `query:"pageToken"`
	Language  string `query:"language" validate:"required" message:"language_invalid_language"`
	Status    string `query:"status"`
	Level     string `query:"level"`
}

type GetWritingExerciseResponse struct {
	Exercises     []WritingExercise `json:"exercises"`
	NextPageToken string            `json:"nextPageToken"`
}

type WritingExercise struct {
	ID         string   `json:"id"`
	Type       string   `json:"type"`
	Level      string   `json:"level"`
	Topic      string   `json:"topic"`
	Question   string   `json:"question"`
	Vocabulary []string `json:"vocabulary"`
	Data       string   `json:"data"`
	Status     string   `json:"status"`
}

func (d WritingExercise) FromDomain(exercise domain.WritingExerciseDatabaseQuery) WritingExercise {
	return WritingExercise{
		ID:         exercise.ID,
		Type:       exercise.Type.String(),
		Level:      exercise.Level.String(),
		Topic:      exercise.Topic,
		Question:   exercise.Question,
		Vocabulary: exercise.Vocabulary,
		Data:       exercise.Data,
		Status:     exercise.Status.String(),
	}
}
