package apperrors

import "errors"

var User = struct {
	InvalidUserID                error
	UserNotFound                 error
	GoogleEmailAlreadyRegistered error
	InvalidSubscriptionPlan      error
	ExceededPlanLimitation       error
}{
	InvalidUserID:                errors.New("user_invalid_id"),
	UserNotFound:                 errors.New("user_not_found"),
	GoogleEmailAlreadyRegistered: errors.New("user_google_email_already_registered"),
	InvalidSubscriptionPlan:      errors.New("user_invalid_subscription_plan"),
	ExceededPlanLimitation:       errors.New("user_exceeded_plan_limitation"),
}
