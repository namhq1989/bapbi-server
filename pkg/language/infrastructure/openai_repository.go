package infrastructure

import (
	"github.com/namhq1989/bapbi-server/internal/openai"
	"github.com/namhq1989/bapbi-server/internal/utils/appcontext"
	"github.com/namhq1989/bapbi-server/pkg/language/domain"
)

type OpenAIRepository struct {
	openai *openai.OpenAI
}

func NewOpenAIRepository(openai *openai.OpenAI) OpenAIRepository {
	return OpenAIRepository{
		openai: openai,
	}
}

func (r OpenAIRepository) SearchTerm(ctx *appcontext.AppContext, term, fromLanguage, toLanguage string) (*domain.OpenAISearchResult, error) {
	result, err := r.openai.SearchTerm(ctx, term, fromLanguage, toLanguage)
	if err != nil {
		return nil, err
	}
	if result == nil {
		return nil, nil
	}

	return &domain.OpenAISearchResult{
		IsValid: result.IsValid,
		Term:    term,
		From: domain.OpenAITermByLanguage{
			Language:   result.From.Language,
			Definition: result.From.Definition,
			Example:    result.From.Example,
		},
		To: domain.OpenAITermByLanguage{
			Language:   result.To.Language,
			Definition: result.To.Definition,
			Example:    result.To.Example,
		},
		Synonyms: result.Synonyms,
		Antonyms: result.Antonyms,
	}, nil
}
