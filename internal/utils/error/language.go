package apperrors

import "errors"

var Language = struct {
	InvalidTerm         error
	InvalidLanguage     error
	InvalidLanguageData error
}{
	InvalidTerm:         errors.New("language_invalid_term"),
	InvalidLanguage:     errors.New("language_invalid_language"),
	InvalidLanguageData: errors.New("language_invalid_language_data"),
}
