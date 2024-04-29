package domain

import (
	"math"
	"time"

	"github.com/namhq1989/bapbi-server/internal/utils/manipulation"

	apperrors "github.com/namhq1989/bapbi-server/internal/utils/error"

	"github.com/namhq1989/bapbi-server/internal/utils/pagetoken"

	"github.com/namhq1989/bapbi-server/internal/utils/appcontext"

	"github.com/namhq1989/bapbi-server/internal/database"
)

type UserWritingExerciseRepository interface {
	CreateUserWritingExercise(ctx *appcontext.AppContext, exercise UserWritingExercise) error
	UpdateUserWritingExercise(ctx *appcontext.AppContext, exercise UserWritingExercise) error
	FindByUserIDAndExerciseID(ctx *appcontext.AppContext, userID, exerciseID string) (*UserWritingExercise, error)
	IsExerciseCreated(ctx *appcontext.AppContext, userID, exerciseID string) (bool, error)
	FindUserWritingExercises(ctx *appcontext.AppContext, filter UserWritingExerciseFilter) ([]WritingExerciseDatabaseQuery, error)
}

//
// FILTER
//

type UserWritingExerciseFilter struct {
	UserID   string
	Language Language
	Status   string
	Time     time.Time
	Limit    int64
}

func NewUserWritingExerciseFilter(uId, lang, stt, pageToken string) UserWritingExerciseFilter {
	language := ToLanguage(lang)
	if !language.IsValid() {
		language = LanguageEnglish
	}

	pt := pagetoken.Decode(pageToken)
	return UserWritingExerciseFilter{
		UserID:   uId,
		Language: language,
		Status:   stt,
		Time:     pt.Timestamp,
		Limit:    10,
	}
}

//
// USER WRITING EXERCISE
//

type UserWritingExerciseAssessment struct {
	IsTopicRelevance bool
	Score            int
	Improvement      []string
	Comment          string
}

type UserWritingExercise struct {
	ID          string
	UserID      string
	ExerciseID  string
	Status      ExerciseStatus
	Language    Language
	Content     string
	Assessment  *UserWritingExerciseAssessment
	CreatedAt   time.Time
	UpdatedAt   time.Time
	CompletedAt time.Time
}

func NewUserWritingExercise(userID, exerciseID, lang string) (*UserWritingExercise, error) {
	language := ToLanguage(lang)
	if !language.IsValid() {
		return nil, apperrors.Language.InvalidLanguage
	}

	return &UserWritingExercise{
		ID:         database.NewStringID(),
		UserID:     userID,
		ExerciseID: exerciseID,
		Status:     ExerciseStatusProgressing,
		Language:   language,
		Assessment: nil,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}, nil
}

func (d *UserWritingExercise) SetContent(content string, minWords int) error {
	var (
		l        = manipulation.CountTotalWords(content)
		maxWords = math.Round(float64(minWords * 2))
	)
	if l < minWords || l > int(maxWords) {
		return apperrors.Language.InvalidWritingExerciseData
	}

	d.Content = content
	return nil
}

func (d *UserWritingExercise) IsProgressing() bool {
	return d.Status == ExerciseStatusProgressing
}

func (d *UserWritingExercise) IsCompleted() bool {
	return d.Status == ExerciseStatusCompleted
}

func (d *UserWritingExercise) SetAssessment(isTopicRelevance bool, score int, improvement []string, comment string) {
	d.Assessment = &UserWritingExerciseAssessment{
		IsTopicRelevance: isTopicRelevance,
		Score:            score,
		Improvement:      improvement,
		Comment:          comment,
	}
}

func (d *UserWritingExercise) SetProgressing() {
	d.Status = ExerciseStatusProgressing
	d.UpdatedAt = time.Now()
}

func (d *UserWritingExercise) SetCompleted() {
	d.Status = ExerciseStatusCompleted
	d.CompletedAt = time.Now()
	d.UpdatedAt = time.Now()
}
