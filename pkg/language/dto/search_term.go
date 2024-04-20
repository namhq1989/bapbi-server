package dto

type SearchTermRequest struct {
	Term string `query:"term" validate:"required" message:"english_invalid_term"`
	From string `query:"from" validate:"required" message:"english_invalid_language"`
	To   string `query:"to" validate:"required" message:"english_invalid_language"`
}

type SearchTermResponse struct {
	Term Term `json:"term"`
}
