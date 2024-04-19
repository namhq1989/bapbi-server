package domain

import "github.com/namhq1989/bapbi-server/internal/utils/appcontext"

type OpenAIRepository interface {
	SearchTerm(ctx *appcontext.AppContext, term, fromLanguage, toLanguage string) (*OpenAISearchTermResult, error)
	SearchSemanticRelations(ctx *appcontext.AppContext, term, language string) (*OpenAISearchSemanticRelationsResult, error)
	SearchPossibleDefinitions(ctx *appcontext.AppContext, term, fromLanguage, toLanguage string) (*OpenAISearchPossibleDefinitionsResult, error)
}

type OpenAISearchTermResult struct {
	IsValid bool
	Term    string
	From    OpenAITermByLanguage
	To      OpenAITermByLanguage
}

type OpenAITermByLanguage struct {
	Language   string
	Definition string
	Example    string
}

type OpenAISearchSemanticRelationsResult struct {
	Synonyms []string
	Antonyms []string
}

type OpenAISearchPossibleDefinitionsResult struct {
	List []TermPossibleDefinition `json:"list"`
}
