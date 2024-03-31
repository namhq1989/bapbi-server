package dto

type SignUpWithGoogleRequest struct {
	Token string `json:"token" validate:"required" message:"auth_invalid_google_token"`
}

type SignUpWithGoogleResponse struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}
