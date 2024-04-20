package dto

type ChangeTermFavouriteRequest struct{}

type ChangeTermFavouriteResponse struct {
	IsFavourite bool `json:"isFavourite"`
}
