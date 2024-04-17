package domain

import (
	"strings"
	"time"

	"github.com/namhq1989/bapbi-server/internal/database"

	apperrors "github.com/namhq1989/bapbi-server/internal/utils/error"
)

type Term struct {
	ID           string
	Term         string
	From         TermByLanguage
	To           TermByLanguage
	Level        string
	PartOfSpeech string
	Phonetic     string
	Synonyms     []string
	Antonyms     []string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type TermByLanguage struct {
	Language   Language
	Definition string
	Example    string
}

func NewTerm(term, fromLanguage, toLanguage string) (*Term, error) {
	if term == "" || len(term) > 32 {
		return nil, apperrors.English.InvalidTerm
	}

	if fromLanguage == toLanguage {
		return nil, apperrors.English.InvalidLanguage
	}

	fromLang := ToLanguage(fromLanguage)
	toLang := ToLanguage(toLanguage)
	if !fromLang.IsValid() || !toLang.IsValid() {
		return nil, apperrors.English.InvalidLanguage
	}

	return &Term{
		ID:   database.NewStringID(),
		Term: strings.ToLower(term),
		From: TermByLanguage{
			Language:   fromLang,
			Definition: "",
			Example:    "",
		},
		To: TermByLanguage{
			Language:   toLang,
			Definition: "",
			Example:    "",
		},
		Level:        "",
		PartOfSpeech: "",
		Phonetic:     "",
		Synonyms:     make([]string, 0),
		Antonyms:     make([]string, 0),
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}, nil
}
