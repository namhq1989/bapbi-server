package domain

import (
	"time"

	"github.com/namhq1989/bapbi-server/internal/database"
	"github.com/namhq1989/bapbi-server/internal/utils/appcontext"
	apperrors "github.com/namhq1989/bapbi-server/internal/utils/error"
	"github.com/namhq1989/bapbi-server/internal/utils/manipulation"
	"github.com/namhq1989/bapbi-server/internal/utils/pagetoken"
)

type UserTermExerciseRepository interface {
	CreateUserTermExercise(ctx *appcontext.AppContext, exercise UserTermExercise) error
	UpdateUserTermExercise(ctx *appcontext.AppContext, exercise UserTermExercise) error
	FindByExerciseID(ctx *appcontext.AppContext, exerciseID string) (*UserTermExercise, error)
	FindByUserIDAndTermID(ctx *appcontext.AppContext, userID, termID string) (*UserTermExercise, error)
	IsExerciseCreated(ctx *appcontext.AppContext, userID, termID string) (bool, error)
	FindUserTermExercises(ctx *appcontext.AppContext, filter UserTermExerciseFilter) ([]UserTermExercise, error)
}

//
// FILTER
//

type UserTermExerciseFilter struct {
	UserID   string
	Language Language
	Status   string
	Time     time.Time
	Limit    int64
}

func NewUserTermExerciseFilter(uId, lang, stt, pageToken string) UserTermExerciseFilter {
	language := ToLanguage(lang)
	if !language.IsValid() {
		language = LanguageEnglish
	}

	pt := pagetoken.Decode(pageToken)
	return UserTermExerciseFilter{
		UserID:   uId,
		Language: language,
		Status:   stt,
		Time:     pt.Timestamp,
		Limit:    10,
	}
}

//
// USER TERM EXERCISE
//

type UserTermExerciseAssessmentGrammarIssue struct {
	Issue      string
	Correction string
}

type UserTermExerciseAssessmentImprovementSuggestion struct {
	Instruction string
	Example     string
}

type UserTermExerciseAssessment struct {
	IsVocabularyCorrect    bool
	VocabularyIssue        string
	IsTenseCorrect         bool
	TenseIssue             string
	GrammarIssues          []UserTermExerciseAssessmentGrammarIssue
	ImprovementSuggestions []UserTermExerciseAssessmentImprovementSuggestion
}

type UserTermExercise struct {
	ID          string
	UserID      string
	TermID      string
	Term        string
	Language    Language
	Tense       GrammarTenseCode
	Content     string
	Status      ExerciseStatus
	Assessment  *UserTermExerciseAssessment
	CreatedAt   time.Time
	UpdatedAt   time.Time
	CompletedAt time.Time
}

func NewUserTermExercise(userID, termID, term, lang, tense string) (*UserTermExercise, error) {
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

	return &UserTermExercise{
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

func (d *UserTermExercise) GetExerciseWordsRange() (int, int) {
	return 5, 30
}

func (d *UserTermExercise) SetContent(content string) error {
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

func (d *UserTermExercise) IsProgressing() bool {
	return d.Status == ExerciseStatusProgressing
}

func (d *UserTermExercise) IsCompleted() bool {
	return d.Status == ExerciseStatusCompleted
}

func (d *UserTermExercise) IsOwner(userID string) bool {
	return d.UserID == userID
}

func (d *UserTermExercise) SetAssessment(
	isVocabularyCorrect bool,
	vocabularyIssue string,
	isTenseCorrect bool,
	tenseIssue string,
	grammarIssues []UserTermExerciseAssessmentGrammarIssue,
	improvementSuggestions []UserTermExerciseAssessmentImprovementSuggestion,
) {
	d.Assessment = &UserTermExerciseAssessment{
		IsVocabularyCorrect:    isVocabularyCorrect,
		VocabularyIssue:        vocabularyIssue,
		IsTenseCorrect:         isTenseCorrect,
		TenseIssue:             tenseIssue,
		GrammarIssues:          grammarIssues,
		ImprovementSuggestions: improvementSuggestions,
	}
}

func (d *UserTermExercise) GetStatusBasedOnAssessment() ExerciseStatus {
	if d.Assessment.IsVocabularyCorrect && d.Assessment.IsTenseCorrect {
		return ExerciseStatusCompleted
	}

	return ExerciseStatusCorrectionRequired
}

func (d *UserTermExercise) SetStatus(status ExerciseStatus) {
	d.Status = status
	d.UpdatedAt = time.Now()

	if status == ExerciseStatusCompleted {
		d.CompletedAt = time.Now()
	}
}
