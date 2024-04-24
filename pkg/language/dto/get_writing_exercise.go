package dto

type GetWritingExerciseRequest struct {
	PageToken string `query:"pageToken"`
	Language  string `query:"language" validate:"required" message:"language_invalid_language"`
	Status    string `query:"status"`
	Level     string `query:"level"`
}

type GetWritingExerciseResponse struct {
	Exercise WritingExercise `json:"exercise"`
}

type WritingExercise struct {
	ID         string   `json:"id"`
	Type       string   `json:"type"`
	Level      string   `json:"level"`
	Topic      string   `json:"topic"`
	Vocabulary []string `json:"vocabulary"`
	Data       string   `json:"data"`
}
