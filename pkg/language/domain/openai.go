package domain

import "github.com/namhq1989/bapbi-server/internal/utils/appcontext"

type OpenAIRepository interface {
	SearchTerm(ctx *appcontext.AppContext, term, fromLang, toLang string) (*OpenAISearchResult, error)
}

type OpenAISearchResult struct {
	IsValid  bool
	Term     string
	From     OpenAITermByLanguage
	To       OpenAITermByLanguage
	Synonyms []string
	Antonyms []string
}

type OpenAITermByLanguage struct {
	Language   string
	Definition string
	Example    string
}
