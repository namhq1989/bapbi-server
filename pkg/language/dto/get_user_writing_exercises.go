package dto

type GetUserWritingExerciseRequest struct {
	PageToken string `query:"pageToken"`
	Language  string `query:"language" validate:"required" message:"language_invalid_language"`
	Status    string `query:"status"`
}

type GetUserWritingExerciseResponse struct {
	Exercises     []WritingExercise `json:"exercises"`
	NextPageToken string            `json:"nextPageToken"`
}
