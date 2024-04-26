package domain

import (
	"strings"
	"time"

	"github.com/namhq1989/bapbi-server/internal/utils/appcontext"

	"github.com/namhq1989/bapbi-server/internal/database"
	apperrors "github.com/namhq1989/bapbi-server/internal/utils/error"

	"github.com/namhq1989/bapbi-server/internal/utils/pagetoken"
)

type WritingExerciseRepository interface {
	CreateWritingExercise(ctx *appcontext.AppContext, exercise WritingExercise) error
	FindByID(ctx *appcontext.AppContext, exerciseID string) (*WritingExercise, error)
	FindWritingExercises(ctx *appcontext.AppContext, filter WritingExerciseFilter) ([]WritingExerciseDatabaseQuery, error)
}

//
// FILTER
//

type WritingExerciseFilter struct {
	UserID   string
	Language Language
	Level    Level
	Time     time.Time
	Limit    int64
}

func NewWritingExerciseFilter(uId, lang, lvl, pageToken string) WritingExerciseFilter {
	language := ToLanguage(lang)
	if !language.IsValid() {
		language = LanguageEnglish
	}

	pt := pagetoken.Decode(pageToken)
	return WritingExerciseFilter{
		UserID:   uId,
		Language: language,
		Level:    ToLevel(lvl),
		Time:     pt.Timestamp,
		Limit:    10,
	}
}

//
// TYPE
//

type WritingExerciseType string

const (
	WritingExerciseTypeUnknown WritingExerciseType = ""
	WritingExerciseTypeBasic   WritingExerciseType = "basic"
	WritingExerciseTypeAnalyze WritingExerciseType = "analyze"
)

var ListWritingExerciseTypes = []WritingExerciseType{WritingExerciseTypeBasic, WritingExerciseTypeAnalyze}

func (s WritingExerciseType) String() string {
	switch s {
	case WritingExerciseTypeBasic, WritingExerciseTypeAnalyze:
		return string(s)
	default:
		return ""
	}
}

func (s WritingExerciseType) IsValid() bool {
	return s != WritingExerciseTypeUnknown
}

func ToWritingExerciseType(value string) WritingExerciseType {
	switch strings.ToLower(value) {
	case WritingExerciseTypeBasic.String():
		return WritingExerciseTypeBasic
	case WritingExerciseTypeAnalyze.String():
		return WritingExerciseTypeAnalyze
	default:
		return WritingExerciseTypeUnknown
	}
}

//
// WRITING EXERCISE
//

type WritingExercise struct {
	ID         string
	Language   Language
	Type       WritingExerciseType
	Level      Level
	Topic      string
	Question   string
	Vocabulary []string
	MinWords   int
	Data       string // JSON string
	CreatedAt  time.Time
}

func NewWritingExercise(lang, exType, lvl, topic, question, data string, vocabulary []string) (*WritingExercise, error) {
	language := ToLanguage(lang)
	if !language.IsValid() {
		return nil, apperrors.Language.InvalidLanguage
	}

	level := ToLevel(lvl)
	if !level.IsValid() {
		return nil, apperrors.Language.InvalidLevel
	}

	exerciseType := ToWritingExerciseType(exType)
	if !exerciseType.IsValid() {
		return nil, apperrors.Language.InvalidWritingExerciseData
	}

	if topic == "" || question == "" {
		return nil, apperrors.Language.InvalidWritingExerciseData
	}

	return &WritingExercise{
		ID:         database.NewStringID(),
		Language:   language,
		Type:       exerciseType,
		Level:      level,
		Topic:      topic,
		Question:   question,
		Vocabulary: vocabulary,
		MinWords:   GetMinWordsBasedOnLevel(level),
		Data:       data,
		CreatedAt:  time.Now(),
	}, nil
}

//
// QUERY
//

type WritingExerciseDatabaseQuery struct {
	WritingExercise
	Status      ExerciseStatus
	CompletedAt time.Time
}
