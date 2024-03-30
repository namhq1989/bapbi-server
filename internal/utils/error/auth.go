package apperrors

import "errors"

var Auth = struct {
	InvalidAuthToken    error
	InvalidGoogleToken  error
	InvalidRefreshToken error
	NotAllowed          error
}{
	InvalidAuthToken:    errors.New("auth_invalid_auth_token"),
	InvalidGoogleToken:  errors.New("auth_invalid_google_token"),
	InvalidRefreshToken: errors.New("auth_invalid_refresh_token"),
	NotAllowed:          errors.New("auth_not_allowed"),
}
