package domain

import (
	"strings"
	"time"

	"github.com/namhq1989/bapbi-server/internal/utils/pagetoken"
)

type WritingExerciseRepository interface {
}

//
// FILTER
//

type WritingExerciseFilter struct {
	Language Language
	Status   WritingExerciseStatus
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
		Status:   ToWritingExerciseStatus(stt),
		Time:     pt.Timestamp,
		Limit:    10,
	}
}

//
// STATUS
//

type WritingExerciseStatus string

const (
	WritingExerciseStatusUnknown   WritingExerciseStatus = ""
	WritingExerciseStatusCompleted WritingExerciseStatus = "completed"
	WritingExerciseStatusAvailable WritingExerciseStatus = "available"
)

func (s WritingExerciseStatus) String() string {
	switch s {
	case WritingExerciseStatusCompleted, WritingExerciseStatusAvailable:
		return string(s)
	default:
		return ""
	}
}

func ToWritingExerciseStatus(value string) WritingExerciseStatus {
	switch strings.ToLower(value) {
	case WritingExerciseStatusCompleted.String():
		return WritingExerciseStatusCompleted
	case WritingExerciseStatusAvailable.String():
		return WritingExerciseStatusAvailable
	default:
		return WritingExerciseStatusUnknown
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
