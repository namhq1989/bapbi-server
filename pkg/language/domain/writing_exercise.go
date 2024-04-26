package domain

import (
	"strings"
	"time"

	"github.com/namhq1989/bapbi-server/internal/database"
	apperrors "github.com/namhq1989/bapbi-server/internal/utils/error"

	"github.com/namhq1989/bapbi-server/internal/utils/pagetoken"
)

type WritingExerciseRepository interface {
}

//
// FILTER
//

type WritingExerciseFilter struct {
	Language Language
	Status   ExerciseStatus
	Level    Level
	Time     time.Time
	Limit    int64
}

func NewWritingExerciseFilter(lang, lvl, stt, pageToken string) WritingExerciseFilter {
	language := ToLanguage(lang)
	if !language.IsValid() {
		language = LanguageEnglish
	}

	pt := pagetoken.Decode(pageToken)
	return WritingExerciseFilter{
		Language: language,
		Level:    ToLevel(lvl),
		Status:   ToExerciseStatus(stt),
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
	Vocabulary []string
	Data       string // JSON string
	CreatedAt  time.Time
}

func NewWritingExercise(lang, exType, lvl, topic, data string) (*WritingExercise, error) {
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

	if topic == "" {
		return nil, apperrors.Language.InvalidWritingExerciseData
	}

	return &WritingExercise{
		ID:         database.NewStringID(),
		Language:   language,
		Type:       exerciseType,
		Level:      level,
		Topic:      topic,
		Vocabulary: []string{},
		Data:       data,
		CreatedAt:  time.Now(),
	}, nil
}
