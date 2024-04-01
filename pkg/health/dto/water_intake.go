package dto

import "time"

type WaterIntakeRequest struct {
	Amount   int       `json:"amount" validate:"required|int|min:1" message:"health_invalid_intake_amount"`
	IntakeAt time.Time `json:"intakeAt"`
}

type WaterIntakeResponse struct{}
