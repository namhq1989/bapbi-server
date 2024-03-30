package apperrors

import "errors"

var User = struct {
	InvalidUserID error
	UserNotFound  error
}{
	InvalidUserID: errors.New("user_invalid_id"),
	UserNotFound:  errors.New("user_not_found"),
}
