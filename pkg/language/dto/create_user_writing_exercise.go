package dto

type CreateUserWritingExerciseRequest struct {
	ExerciseID string `json:"exerciseId" validate:"required" message:"language_invalid_exercise_id"`
}

type CreateUserWritingExerciseResponse struct {
	ID string `json:"id"`
}
