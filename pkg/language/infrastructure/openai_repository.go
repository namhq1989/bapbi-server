package infrastructure

import (
	"strings"

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

func (r OpenAIRepository) SearchTerm(ctx *appcontext.AppContext, term, fromLanguage, toLanguage string) (*domain.OpenAISearchTermResult, error) {
	result, err := r.openai.SearchTerm(ctx, term, fromLanguage, toLanguage)
	if err != nil {
		return nil, err
	}
	if result == nil {
		return nil, nil
	}

	ctx.Logger().Print("got new term", result)

	examples := make([]domain.TermExample, len(result.Examples))
	for i, example := range result.Examples {
		examples[i] = domain.TermExample{
			PartOfSpeech: example.PartOfSpeech,
			From:         example.From,
			To:           example.To,
		}
	}

	return &domain.OpenAISearchTermResult{
		From: domain.TermByLanguage{
			Language:   domain.ToLanguage(strings.ToLower(result.From.Language)),
			Definition: result.From.Definition,
			Example:    result.From.Example,
		},
		To: domain.TermByLanguage{
			Language:   domain.ToLanguage(strings.ToLower(result.To.Language)),
			Definition: result.To.Definition,
			Example:    result.To.Example,
		},
		Examples: examples,
	}, nil
}

func (r OpenAIRepository) SearchSemanticRelations(ctx *appcontext.AppContext, term, language string) (*domain.OpenAISearchSemanticRelationsResult, error) {
	result, err := r.openai.SearchSemanticRelations(ctx, term, language)
	if err != nil {
		return nil, err
	}
	if result == nil {
		return nil, nil
	}

	return &domain.OpenAISearchSemanticRelationsResult{
		Synonyms: result.Synonyms,
		Antonyms: result.Antonyms,
	}, nil
}
