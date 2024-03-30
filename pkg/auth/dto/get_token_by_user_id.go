package dto

type GetAccessTokenByUserIDRequest struct {
	UserID string `query:"userId" validate:"required" message:"user_invalid_id"`
}

type GetAccessTokenByUserIDResponse struct {
	AccessToken string `json:"accessToken"`
}
