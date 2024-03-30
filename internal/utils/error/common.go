package apperrors

import "errors"

var Common = struct {
	Success             error
	BadRequest          error
	NotFound            error
	Unauthorized        error
	Forbidden           error
	EmailAlreadyExisted error
	InvalidID           error
	InvalidName         error
	InvalidEmail        error
	InvalidStatus       error
	InvalidHeight       error
	InvalidWeight       error
	InvalidWakingHours  error
}{
	Success:             errors.New("success"),
	BadRequest:          errors.New("bad_request"),
	NotFound:            errors.New("not_found"),
	Unauthorized:        errors.New("unauthorized"),
	Forbidden:           errors.New("forbidden"),
	EmailAlreadyExisted: errors.New("email_already_existed"),
	InvalidID:           errors.New("invalid_id"),
	InvalidName:         errors.New("invalid_name"),
	InvalidEmail:        errors.New("invalid_email"),
	InvalidStatus:       errors.New("invalid_status"),
	InvalidHeight:       errors.New("invalid_height"),
	InvalidWeight:       errors.New("invalid_weight"),
	InvalidWakingHours:  errors.New("invalid_waking_hours"),
}
