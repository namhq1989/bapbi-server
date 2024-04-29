package query

import (
	"github.com/namhq1989/bapbi-server/internal/utils/appcontext"
	"github.com/namhq1989/bapbi-server/pkg/language/domain"
	"github.com/namhq1989/bapbi-server/pkg/language/dto"
)

type GetDataHandler struct{}

func NewGetDataHandler() GetDataHandler {
	return GetDataHandler{}
}

func (h GetDataHandler) GetData(_ *appcontext.AppContext, _ dto.GetDataRequest) (*dto.GetDataResponse, error) {
	return &dto.GetDataResponse{
		GrammarTenses: dto.GetDataResponse{}.FromDomain(domain.EnglishGrammarTenses),
	}, nil
}
