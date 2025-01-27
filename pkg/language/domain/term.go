package domain

import (
	"strings"
	"time"

	"github.com/namhq1989/bapbi-server/internal/utils/appcontext"

	"github.com/namhq1989/bapbi-server/internal/database"

	apperrors "github.com/namhq1989/bapbi-server/internal/utils/error"
)

type TermRepository interface {
	FindByID(ctx *appcontext.AppContext, termID string) (*Term, error)
	FindByTerm(ctx *appcontext.AppContext, term, fromLanguage string) (*Term, error)
	CreateTerm(ctx *appcontext.AppContext, term Term) error
	UpdateTerm(ctx *appcontext.AppContext, term Term) error
	FindFeaturedWord(ctx *appcontext.AppContext, language string) (*Term, error)
}

type Term struct {
	ID           string
	Term         string
	From         TermByLanguage
	To           TermByLanguage
	Level        string
	PartOfSpeech string
	Phonetic     string
	AudioURL     string
	ReferenceURL string
	Synonyms     []string
	Antonyms     []string
	Examples     []TermExample
	IsFeatured   bool
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type TermByLanguage struct {
	Language   Language
	Definition string
	Example    string
}

type TermExample struct {
	PartOfSpeech string
	From         string
	To           string
}

func NewTerm(term, fromLanguage, toLanguage string) (*Term, error) {
	if term == "" || len(term) > 32 {
		return nil, apperrors.Language.InvalidTerm
	}

	if fromLanguage == toLanguage {
		return nil, apperrors.Language.InvalidLanguage
	}

	fromLang := ToLanguage(fromLanguage)
	toLang := ToLanguage(toLanguage)
	if !fromLang.IsValid() || !toLang.IsValid() {
		return nil, apperrors.Language.InvalidLanguage
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
		Examples:     make([]TermExample, 0),
		IsFeatured:   false,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}, nil
}

func (d *Term) SetLanguage(language Language, definition, example string) error {
	if !language.IsValid() {
		return apperrors.Language.InvalidLanguage
	}

	if definition == "" || example == "" {
		return apperrors.Language.InvalidLanguageData
	}

	if d.From.Language == language {
		d.From.Definition = definition
		d.From.Example = example
	} else if d.To.Language == language {
		d.To.Definition = definition
		d.To.Example = example
	} else {
		return apperrors.Language.InvalidLanguage
	}

	return nil
}

func (d *Term) SetLevel(level string) {
	d.Level = strings.ToLower(level)
}

func (d *Term) SetPartOfSpeech(pos string) {
	d.PartOfSpeech = strings.ToLower(pos)
}

func (d *Term) SetPhonetic(phonetic string) {
	d.Phonetic = phonetic
}

func (d *Term) SetAudioURL(url string) {
	d.AudioURL = url
}

func (d *Term) SetReferenceURL(url string) {
	d.ReferenceURL = url
}

func (d *Term) SetSynonyms(synonyms []string) {
	d.Synonyms = synonyms
}

func (d *Term) SetExamples(examples []TermExample) {
	d.Examples = examples
}

func (d *Term) SetAntonyms(antonyms []string) {
	d.Antonyms = antonyms
}

func (d *Term) SetIsFeatured(isFeatured bool) {
	d.IsFeatured = isFeatured
}

func (d *Term) SetUpdatedAt() {
	d.UpdatedAt = time.Now()
}
