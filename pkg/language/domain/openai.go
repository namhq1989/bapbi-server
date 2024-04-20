package domain

import "github.com/namhq1989/bapbi-server/internal/utils/appcontext"

type OpenAIRepository interface {
	SearchTerm(ctx *appcontext.AppContext, term, fromLanguage, toLanguage string) (*OpenAISearchTermResult, error)
	SearchSemanticRelations(ctx *appcontext.AppContext, term, language string) (*OpenAISearchSemanticRelationsResult, error)
}

type OpenAISearchTermResult struct {
	From     TermByLanguage
	To       TermByLanguage
	Examples []TermExample
}

type OpenAISearchSemanticRelationsResult struct {
	Synonyms []string
	Antonyms []string
}
