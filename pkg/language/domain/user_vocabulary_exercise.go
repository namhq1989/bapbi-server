package domain

import (
	"time"

	"github.com/namhq1989/bapbi-server/internal/database"
	"github.com/namhq1989/bapbi-server/internal/utils/appcontext"
	apperrors "github.com/namhq1989/bapbi-server/internal/utils/error"
	"github.com/namhq1989/bapbi-server/internal/utils/manipulation"
)

type UserVocabularyExerciseRepository interface {
	CreateUserVocabularyExercise(ctx *appcontext.AppContext, exercise UserVocabularyExercise) error
	UpdateUserVocabularyExercise(ctx *appcontext.AppContext, exercise UserVocabularyExercise) error
	FindByUserIDAndTermID(ctx *appcontext.AppContext, userID, termID string) (*UserVocabularyExercise, error)
	IsExerciseCreated(ctx *appcontext.AppContext, userID, termID string) (bool, error)
}

//
// USER VOCABULARY EXERCISE
//

type UserVocabularyExerciseAssessment struct {
	IsRelevance bool
	Score       int
	Improvement []string
	Comment     string
}

type UserVocabularyExercise struct {
	ID          string
	UserID      string
	TermID      string
	Term        string
	Language    Language
	Tense       GrammarTenseCode
	Content     string
	Status      ExerciseStatus
	Assessment  *UserVocabularyExerciseAssessment
	CreatedAt   time.Time
	UpdatedAt   time.Time
	CompletedAt time.Time
}

func NewUserVocabularyExercise(userID, termID, term, lang, tense string) (*UserVocabularyExercise, error) {
	if term == "" {
		return nil, apperrors.Language.InvalidTerm
	}

	language := ToLanguage(lang)
	if !language.IsValid() {
		return nil, apperrors.Language.InvalidLanguage
	}

	tenseCode := ToGrammarTenseCode(tense)
	if !tenseCode.IsValid() {
		return nil, apperrors.Language.InvalidVocabularyExerciseData
	}

	return &UserVocabularyExercise{
		ID:         database.NewStringID(),
		UserID:     userID,
		TermID:     termID,
		Term:       term,
		Language:   language,
		Tense:      tenseCode,
		Content:    "",
		Assessment: nil,
		Status:     ExerciseStatusProgressing,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}, nil
}

func (d *UserVocabularyExercise) GetExerciseWordsRange() (int, int) {
	return 5, 30
}

func (d *UserVocabularyExercise) SetContent(content string) error {
	var (
		l                  = manipulation.CountTotalWords(content)
		minWords, maxWords = d.GetExerciseWordsRange()
	)

	if l < minWords || l > maxWords {
		return apperrors.Language.InvalidVocabularyExerciseData
	}

	d.Content = content
	return nil
}

func (d *UserVocabularyExercise) IsProgressing() bool {
	return d.Status == ExerciseStatusProgressing
}

func (d *UserVocabularyExercise) IsCompleted() bool {
	return d.Status == ExerciseStatusCompleted
}

func (d *UserVocabularyExercise) SetAssessment(isRelevance bool, score int, improvement []string, comment string) {
	d.Assessment = &UserVocabularyExerciseAssessment{
		IsRelevance: isRelevance,
		Score:       score,
		Improvement: improvement,
		Comment:     comment,
	}
}

func (d *UserVocabularyExercise) SetProgressing() {
	d.Status = ExerciseStatusProgressing
	d.UpdatedAt = time.Now()
}

func (d *UserVocabularyExercise) SetCompleted() {
	d.Status = ExerciseStatusCompleted
	d.CompletedAt = time.Now()
	d.UpdatedAt = time.Now()
}
