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

func (r OpenAIRepository) SearchTerm(ctx *appcontext.AppContext, term, fromLanguage, toLanguage string) (*domain.OpenAISearchTermResult, error) {
	result, err := r.openai.SearchTerm(ctx, term, fromLanguage, toLanguage)
	if err != nil {
		return nil, err
	}
	if result == nil || !result.IsValid {
		return nil, nil
	}

	return &domain.OpenAISearchTermResult{
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

	ctx.Logger().Print("SearchSemanticRelations", result)

	return &domain.OpenAISearchSemanticRelationsResult{
		Synonyms: result.Synonyms,
		Antonyms: result.Antonyms,
	}, nil
}

func (r OpenAIRepository) SearchPossibleDefinitions(ctx *appcontext.AppContext, term, fromLanguage, toLanguage string) (*domain.OpenAISearchPossibleDefinitionsResult, error) {
	result, err := r.openai.SearchPossibleDefinitions(ctx, term, fromLanguage, toLanguage)
	if err != nil {
		return nil, err
	}
	if result == nil {
		return nil, nil
	}

	list := make([]domain.TermPossibleDefinition, len(result.List))
	for i, item := range result.List {
		list[i] = domain.TermPossibleDefinition{
			Definition:   item.Definition,
			PartOfSpeech: item.Pos,
		}
	}

	return &domain.OpenAISearchPossibleDefinitionsResult{
		List: list,
	}, nil
}
