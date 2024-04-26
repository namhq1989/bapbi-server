package dto

type StartUserWritingExerciseRequest struct {
	ExerciseID string `json:"exerciseId" validate:"required" message:"language_invalid_exercise_id"`
}

type StartUserWritingExerciseResponse struct {
	ID string `json:"id"`
}
