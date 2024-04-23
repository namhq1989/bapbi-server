package query

import (
	"github.com/namhq1989/bapbi-server/internal/utils/appcontext"
	"github.com/namhq1989/bapbi-server/pkg/language/domain"
	"github.com/namhq1989/bapbi-server/pkg/language/dto"
)

type GetFeaturedTermHandler struct {
	termRepository domain.TermRepository
}

func NewGetFeaturedTermHandler(termRepository domain.TermRepository) GetFeaturedTermHandler {
	return GetFeaturedTermHandler{
		termRepository: termRepository,
	}
}

func (h GetFeaturedTermHandler) GetFeaturedTerm(ctx *appcontext.AppContext, req dto.GetFeaturedTermRequest) (*dto.GetFeaturedTermResponse, error) {
	ctx.Logger().Info("new get featured term request", appcontext.Fields{"language": req.Language})

	ctx.Logger().Text("find in db")
	term, err := h.termRepository.FindFeaturedWord(ctx, req.Language)
	if err != nil {
		ctx.Logger().Error("failed to find in db", err, appcontext.Fields{})
		return nil, err
	}
	if term == nil {
		ctx.Logger().Error("term not found, respond", nil, appcontext.Fields{})
		return &dto.GetFeaturedTermResponse{Term: nil}, nil
	}

	ctx.Logger().Text("done get featured word")
	result := dto.Term{}.FromDomain(*term, false)
	return &dto.GetFeaturedTermResponse{Term: &result}, nil
}
