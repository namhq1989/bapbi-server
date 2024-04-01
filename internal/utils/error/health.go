package apperrors

import "errors"

var Health = struct {
	HealthProfileNotFound     error
	DrinkWaterProfileNotFound error
	InvalidHeight             error
	InvalidWeight             error
	InvalidBMI                error
	InvalidWakingHours        error
	InvalidIntakeAmount       error
	InvalidDailyIntakeAmount  error
	InvalidHourlyIntakeAmount error
	InvalidStreak             error
}{
	HealthProfileNotFound:     errors.New("health_profile_not_found"),
	DrinkWaterProfileNotFound: errors.New("health_drink_water_profile_not_found"),
	InvalidHeight:             errors.New("health_invalid_height"),
	InvalidWeight:             errors.New("health_invalid_weight"),
	InvalidBMI:                errors.New("health_invalid_bmi"),
	InvalidWakingHours:        errors.New("health_invalid_waking_hours"),
	InvalidIntakeAmount:       errors.New("health_invalid_intake_amount"),
	InvalidDailyIntakeAmount:  errors.New("health_invalid_daily_intake_amount"),
	InvalidHourlyIntakeAmount: errors.New("health_invalid_hourly_intake_amount"),
	InvalidStreak:             errors.New("health_invalid_streak"),
}
