package domain

import (
	"time"

	apperrors "github.com/namhq1989/bapbi-server/internal/utils/error"

	"github.com/namhq1989/bapbi-server/internal/utils/pagetoken"

	"github.com/namhq1989/bapbi-server/internal/utils/appcontext"

	"github.com/namhq1989/bapbi-server/internal/database"
)

type UserWritingExerciseRepository interface {
	CreateUserWritingExercise(ctx *appcontext.AppContext, exercise UserWritingExercise) error
	UpdateUserWritingExercise(ctx *appcontext.AppContext, exercise UserWritingExercise) error
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

type UserWritingExercise struct {
	ID          string
	UserID      string
	ExerciseID  string
	Status      ExerciseStatus
	Language    Language
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
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}, nil
}

func (d *UserWritingExercise) SetCompleted() {
	d.Status = ExerciseStatusCompleted
	d.CompletedAt = time.Now()
	d.UpdatedAt = time.Now()
}
