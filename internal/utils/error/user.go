package apperrors

import "errors"

var User = struct {
	InvalidUserID                error
	UserNotFound                 error
	GoogleEmailAlreadyRegistered error
}{
	InvalidUserID:                errors.New("user_invalid_id"),
	UserNotFound:                 errors.New("user_not_found"),
	GoogleEmailAlreadyRegistered: errors.New("user_google_email_already_registered"),
}
