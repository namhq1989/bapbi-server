package dto

type GetUserTermsRequest struct {
	PageToken  string `query:"pageToken"`
	Keyword    string `query:"keyword"`
	IsFavorite string `query:"isFavorite"`
}

type GetUserTermsResponse struct {
	Terms         []Term `json:"terms"`
	NextPageToken string `json:"nextPageToken"`
}

func (d GetUserTermsResponse) DefaultValue() GetUserTermsResponse {
	return GetUserTermsResponse{
		Terms:         make([]Term, 0),
		NextPageToken: "",
	}
}
