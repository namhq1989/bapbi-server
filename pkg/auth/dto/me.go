package dto

type MeRequest struct{}

type MeResponse struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}
