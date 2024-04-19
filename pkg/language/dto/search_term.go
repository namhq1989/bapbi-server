package dto

import (
	"github.com/namhq1989/bapbi-server/pkg/language/domain"
)

type SearchTermRequest struct {
	Term string `query:"term" validate:"required" message:"english_invalid_term"`
	From string `query:"from" validate:"required" message:"english_invalid_language"`
	To   string `query:"to" validate:"required" message:"english_invalid_language"`
}

type SearchTermResponse struct {
	ID           string         `json:"id"`
	Term         string         `json:"term"`
	From         TermByLanguage `json:"from"`
	To           TermByLanguage `json:"to"`
	Level        string         `json:"level"`
	PartOfSpeech string         `json:"pos"`
	Phonetic     string         `json:"phonetic"`
	ReferenceURL string         `json:"referenceUrl"`
	AudioURL     string         `json:"audioUrl"`
	Synonyms     []string       `json:"synonyms"`
	Antonyms     []string       `json:"antonyms"`
}

type TermByLanguage struct {
	Language   string `json:"language"`
	Definition string `json:"definition"`
	Example    string `json:"example"`
}

func (d SearchTermResponse) FromDomain(term domain.Term) SearchTermResponse {
	return SearchTermResponse{
		ID:   term.ID,
		Term: term.Term,
		From: TermByLanguage{
			Language:   term.From.Language.String(),
			Definition: term.From.Definition,
			Example:    term.From.Example,
		},
		To: TermByLanguage{
			Language:   term.To.Language.String(),
			Definition: term.To.Definition,
			Example:    term.To.Example,
		},
		Level:        term.Level,
		PartOfSpeech: term.PartOfSpeech,
		Phonetic:     term.Phonetic,
		ReferenceURL: term.ReferenceURL,
		AudioURL:     term.AudioURL,
		Synonyms:     term.Synonyms,
		Antonyms:     term.Antonyms,
	}
}
