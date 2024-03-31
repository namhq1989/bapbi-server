package dto

type SignInWithGoogleRequest struct {
	Token string `json:"token" validate:"required" message:"auth_invalid_google_token"`
}

type SignInWithGoogleResponse struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}
