package apperrors

import "errors"

var Language = struct {
	TermNotFound               error
	InvalidTerm                error
	InvalidLevel               error
	InvalidLanguage            error
	InvalidLanguageData        error
	InvalidWritingExerciseData error
	InvalidExerciseStatus      error
}{
	TermNotFound:               errors.New("language_term_not_found"),
	InvalidTerm:                errors.New("language_invalid_term"),
	InvalidLevel:               errors.New("language_invalid_level"),
	InvalidLanguage:            errors.New("language_invalid_language"),
	InvalidLanguageData:        errors.New("language_invalid_language_data"),
	InvalidWritingExerciseData: errors.New("language_invalid_writing_exercise_data"),
	InvalidExerciseStatus:      errors.New("language_invalid_exercise_status"),
}
