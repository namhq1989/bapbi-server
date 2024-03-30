package dto

type RefreshAccessTokenRequest struct {
	UserID       string `json:"userId" validate:"required" message:"user_invalid_id"`
	RefreshToken string `json:"refreshToken" validate:"required" message:"auth_invalid_refresh_token"`
}

type RefreshAccessTokenResponse struct {
	AccessToken string `json:"accessToken"`
}
