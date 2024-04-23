package dto

type GetFeaturedTermRequest struct {
	Language string `query:"language" validate:"required" message:"language_invalid_language"`
}

type GetFeaturedTermResponse struct {
	Term *Term `json:"term"`
}
