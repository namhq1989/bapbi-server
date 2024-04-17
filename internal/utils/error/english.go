package apperrors

import "errors"

var English = struct {
	InvalidTerm     error
	InvalidLanguage error
}{
	InvalidTerm:     errors.New("english_invalid_term"),
	InvalidLanguage: errors.New("english_invalid_language"),
}
