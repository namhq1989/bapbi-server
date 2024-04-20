package model

import (
	"time"

	"github.com/namhq1989/bapbi-server/internal/database"
	apperrors "github.com/namhq1989/bapbi-server/internal/utils/error"
	"github.com/namhq1989/bapbi-server/pkg/language/domain"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Term struct {
	ID           primitive.ObjectID `bson:"_id"`
	Term         string             `bson:"term"`
	From         TermByLanguage     `bson:"from"`
	To           TermByLanguage     `bson:"to"`
	Level        string             `bson:"level"`
	PartOfSpeech string             `bson:"partOfSpeech"`
	Phonetic     string             `bson:"phonetic"`
	AudioURL     string             `bson:"audioUrl"`
	ReferenceURL string             `bson:"referenceUrl"`
	Synonyms     []string           `bson:"synonyms"`
	Antonyms     []string           `bson:"antonyms"`
	Examples     []TermExample      `bson:"examples"`
	CreatedAt    time.Time          `bson:"createdAt"`
	UpdatedAt    time.Time          `bson:"updatedAt"`
}

type TermByLanguage struct {
	Language   string `bson:"language"`
	Definition string `bson:"definition"`
	Example    string `bson:"example"`
}

type TermExample struct {
	PartOfSpeech string `bson:"pos"`
	From         string `bson:"from"`
	To           string `bson:"to"`
}

func (m Term) ToDomain() domain.Term {
	examples := make([]domain.TermExample, len(m.Examples))
	for i, example := range m.Examples {
		examples[i] = domain.TermExample{
			PartOfSpeech: example.PartOfSpeech,
			From:         example.From,
			To:           example.To,
		}
	}
	return domain.Term{
		ID:   m.ID.Hex(),
		Term: m.Term,
		From: domain.TermByLanguage{
			Language:   domain.ToLanguage(m.From.Language),
			Definition: m.From.Definition,
			Example:    m.From.Example,
		},
		To: domain.TermByLanguage{
			Language:   domain.ToLanguage(m.To.Language),
			Definition: m.To.Definition,
			Example:    m.To.Example,
		},
		Level:        m.Level,
		PartOfSpeech: m.PartOfSpeech,
		Phonetic:     m.Phonetic,
		AudioURL:     m.AudioURL,
		ReferenceURL: m.ReferenceURL,
		Synonyms:     m.Synonyms,
		Antonyms:     m.Antonyms,
		Examples:     examples,
		CreatedAt:    m.CreatedAt,
		UpdatedAt:    m.UpdatedAt,
	}
}

func (m Term) FromDomain(term domain.Term) (*Term, error) {
	id, err := database.ObjectIDFromString(term.ID)
	if err != nil {
		return nil, apperrors.Common.InvalidID
	}

	examples := make([]TermExample, len(term.Examples))
	for i, example := range term.Examples {
		examples[i] = TermExample{
			PartOfSpeech: example.PartOfSpeech,
			From:         example.From,
			To:           example.To,
		}
	}

	return &Term{
		ID:           id,
		Term:         term.Term,
		From:         TermByLanguage{Language: term.From.Language.String(), Definition: term.From.Definition, Example: term.From.Example},
		To:           TermByLanguage{Language: term.To.Language.String(), Definition: term.To.Definition, Example: term.To.Example},
		Level:        term.Level,
		PartOfSpeech: term.PartOfSpeech,
		Phonetic:     term.Phonetic,
		AudioURL:     term.AudioURL,
		ReferenceURL: term.ReferenceURL,
		Synonyms:     term.Synonyms,
		Antonyms:     term.Antonyms,
		Examples:     examples,
		CreatedAt:    term.CreatedAt,
		UpdatedAt:    term.UpdatedAt,
	}, nil
}
