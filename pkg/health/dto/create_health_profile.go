package dto

type CreateHealthProfileRequest struct {
	Height      int `json:"height" validate:"required|int|min:1" message:"health_invalid_height"`
	Weight      int `json:"weight" validate:"required|int|min:1" message:"health_invalid_weight"`
	WakeUpHour  int `json:"wakeupHour" validate:"required|int|min:0|max:23" message:"health_invalid_waking_hours"`
	BedtimeHour int `json:"bedtimeHour" validate:"required|int|min:0|max:23" message:"health_invalid_waking_hours"`
}

type CreateHealthProfileResponse struct {
	ID string `json:"id"`
}
