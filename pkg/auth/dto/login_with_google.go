package dto

type LoginWithGoogleRequest struct {
	Token string `json:"token" validate:"required" message:"auth_invalid_google_token"`
}

type LoginWithGoogleResponse struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}
