package dto

type SubmitUserWritingExerciseRequest struct {
	ExerciseID string `json:"exerciseId" validate:"required" message:"language_invalid_exercise_id"`
	Content    string `json:"content" validate:"required" message:"language_invalid_writing_exercise_data"`
}

type SubmitUserWritingExerciseResponse struct {
	ID string `json:"id"`
}
