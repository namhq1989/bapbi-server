package query

import (
	"github.com/namhq1989/bapbi-server/internal/utils/appcontext"
	"github.com/namhq1989/bapbi-server/internal/utils/pagetoken"
	"github.com/namhq1989/bapbi-server/pkg/language/domain"
	"github.com/namhq1989/bapbi-server/pkg/language/dto"
)

type GetUserTermsHandler struct {
	termRepository     domain.TermRepository
	userTermRepository domain.UserTermRepository
}

func NewGetUserTermsHandler(
	termRepository domain.TermRepository,
	userTermRepository domain.UserTermRepository,
) GetUserTermsHandler {
	return GetUserTermsHandler{
		termRepository:     termRepository,
		userTermRepository: userTermRepository,
	}
}

func (h GetUserTermsHandler) GetUserTerms(ctx *appcontext.AppContext, performerID string, req dto.GetUserTermsRequest) (*dto.GetUserTermsResponse, error) {
	ctx.Logger().Info("new get user terms request", appcontext.Fields{"performer": performerID, "keyword": req.Keyword, "pageToken": req.PageToken, "isFavorite": req.IsFavorite})

	ctx.Logger().Text("find user's terms in db")
	filter := domain.NewUserTermFilter(performerID, req.Keyword, req.IsFavorite, req.PageToken)
	userTerms, err := h.userTermRepository.FindUserTerms(ctx, filter)
	if err != nil {
		ctx.Logger().Error("failed to find user's terms in db", err, appcontext.Fields{"performer": performerID, "keyword": req.Keyword, "pageToken": req.PageToken, "isFavorite": req.IsFavorite})
		return nil, err
	}

	totalTerms := int64(len(userTerms))
	if totalTerms == 0 {
		ctx.Logger().Text("no terms found")
		result := dto.GetUserTermsResponse{}.DefaultValue()
		return &result, nil
	}

	ctx.Logger().Text("fetch terms data")

	var result = &dto.GetUserTermsResponse{
		Terms: make([]dto.Term, totalTerms), NextPageToken: "",
	}
	for i, userTerm := range userTerms {
		term, _ := h.termRepository.FindByID(ctx, userTerm.TermID)
		if term == nil {
			result.Terms[i] = dto.Term{}
			continue
		}

		result.Terms[i] = dto.Term{}.FromDomain(*term, userTerm.IsFavourite)
	}

	ctx.Logger().Text("generate page token")
	if totalTerms >= filter.Limit {
		result.NextPageToken = pagetoken.NewWithTimestamp(userTerms[totalTerms-1].CreatedAt)
	}

	ctx.Logger().Text("done get user terms")
	return result, nil
}
