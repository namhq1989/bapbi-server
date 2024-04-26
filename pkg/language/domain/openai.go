package domain

import "github.com/namhq1989/bapbi-server/internal/utils/appcontext"

type OpenAIRepository interface {
	SearchTerm(ctx *appcontext.AppContext, term, fromLanguage, toLanguage string) (*OpenAISearchTermResult, error)
	SearchSemanticRelations(ctx *appcontext.AppContext, term, language string) (*OpenAISearchSemanticRelationsResult, error)
	FeaturedWord(ctx *appcontext.AppContext, language string) (*OpenAIFeaturedWordResult, error)
	WritingExercise(ctx *appcontext.AppContext, language, exType, level string) (*OpenAIWritingExerciseResult, error)
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

type OpenAIFeaturedWordResult struct {
	Word string `json:"word"`
}

type OpenAIWritingExerciseResult struct {
	Topic      string
	Question   string
	Vocabulary []string
	Data       string
}
