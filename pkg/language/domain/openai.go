package domain

import "github.com/namhq1989/bapbi-server/internal/utils/appcontext"

type OpenAIRepository interface {
	SearchTerm(ctx *appcontext.AppContext, term, fromLanguage, toLanguage string) (*OpenAISearchTermResult, error)
	SearchSemanticRelations(ctx *appcontext.AppContext, term, language string) (*OpenAISearchSemanticRelationsResult, error)
	GenerateFeaturedWord(ctx *appcontext.AppContext, language string) (*OpenAIFeaturedWordResult, error)
	GenerateWritingExercise(ctx *appcontext.AppContext, language, exType, level string) (*OpenAIWritingExerciseResult, error)
	AssessWritingExercise(ctx *appcontext.AppContext, language, topic, level, content string) (*OpenAIAssessWritingExerciseResult, error)
	AssessVocabularyExercise(ctx *appcontext.AppContext, language, term, tense, content string) (*OpenAIAssessVocabularyExerciseResult, error)
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

type OpenAIAssessWritingExerciseResult struct {
	IsTopicRelevance bool
	Score            int
	Improvement      []string
	Comment          string
}

type OpenAIAssessVocabularyExerciseResult struct {
	IsVocabularyCorrect    bool
	VocabularyIssue        string
	IsTenseCorrect         bool
	TenseIssue             string
	GrammarIssues          []UserVocabularyExerciseAssessmentGrammarIssue
	ImprovementSuggestions []UserVocabularyExerciseAssessmentImprovementSuggestion
}
