package apperrors

import "errors"

var Language = struct {
	TermNotFound        error
	InvalidTerm         error
	InvalidLanguage     error
	InvalidLanguageData error
}{
	TermNotFound:        errors.New("language_term_not_found"),
	InvalidTerm:         errors.New("language_invalid_term"),
	InvalidLanguage:     errors.New("language_invalid_language"),
	InvalidLanguageData: errors.New("language_invalid_language_data"),
}
