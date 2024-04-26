package apperrors

import "errors"

var Language = struct {
	TermNotFound               error
	ExerciseNotFound           error
	InvalidTerm                error
	InvalidLevel               error
	InvalidLanguage            error
	InvalidLanguageData        error
	InvalidWritingExerciseData error
	InvalidExerciseID          error
	InvalidExerciseStatus      error
	UserExerciseExisted        error
	InvalidUserAction          error
	ExerciseAlreadyCompleted   error
}{
	TermNotFound:               errors.New("language_term_not_found"),
	ExerciseNotFound:           errors.New("language_exercise_not_found"),
	InvalidTerm:                errors.New("language_invalid_term"),
	InvalidLevel:               errors.New("language_invalid_level"),
	InvalidLanguage:            errors.New("language_invalid_language"),
	InvalidLanguageData:        errors.New("language_invalid_language_data"),
	InvalidWritingExerciseData: errors.New("language_invalid_writing_exercise_data"),
	InvalidExerciseID:          errors.New("language_invalid_exercise_id"),
	InvalidExerciseStatus:      errors.New("language_invalid_exercise_status"),
	UserExerciseExisted:        errors.New("language_user_exercise_existed"),
	InvalidUserAction:          errors.New("language_invalid_user_action"),
	ExerciseAlreadyCompleted:   errors.New("language_exercise_already_completed"),
}
