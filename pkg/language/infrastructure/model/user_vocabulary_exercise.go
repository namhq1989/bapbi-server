package model

import (
	"time"

	apperrors "github.com/namhq1989/bapbi-server/internal/utils/error"
	"github.com/namhq1989/bapbi-server/pkg/language/domain"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserVocabularyExerciseAssessmentGrammarIssue struct {
	Issue      string `bson:"issue"`
	Correction string `bson:"correction"`
}

type UserVocabularyExerciseAssessmentImprovementSuggestion struct {
	Instruction string `bson:"instruction"`
	Example     string `bson:"example"`
}

type UserVocabularyExerciseAssessment struct {
	IsVocabularyCorrect    bool                                                    `bson:"isVocabularyCorrect"`
	VocabularyIssue        string                                                  `bson:"vocabularyIssue"`
	IsTenseCorrect         bool                                                    `bson:"isTenseCorrect"`
	TenseIssue             string                                                  `bson:"tenseIssue"`
	GrammarIssues          []UserVocabularyExerciseAssessmentGrammarIssue          `bson:"grammarIssues"`
	ImprovementSuggestions []UserVocabularyExerciseAssessmentImprovementSuggestion `bson:"improvementSuggestions"`
}

type UserVocabularyExercise struct {
	ID          primitive.ObjectID                `bson:"_id"`
	UserID      primitive.ObjectID                `bson:"userId"`
	TermID      primitive.ObjectID                `bson:"termId"`
	Term        string                            `bson:"term"`
	Language    string                            `bson:"language"`
	Tense       string                            `bson:"tense"`
	Content     string                            `bson:"content"`
	Status      string                            `bson:"status"`
	Assessment  *UserVocabularyExerciseAssessment `bson:"assessment"`
	CreatedAt   time.Time                         `bson:"createdAt"`
	UpdatedAt   time.Time                         `bson:"updatedAt"`
	CompletedAt time.Time                         `bson:"completedAt"`
}

func (u UserVocabularyExercise) ToDomain() domain.UserVocabularyExercise {
	exercise := domain.UserVocabularyExercise{
		ID:          u.ID.Hex(),
		UserID:      u.UserID.Hex(),
		TermID:      u.TermID.Hex(),
		Term:        u.Term,
		Language:    domain.ToLanguage(u.Language),
		Tense:       domain.ToGrammarTenseCode(u.Tense),
		Content:     u.Content,
		Status:      domain.ToExerciseStatus(u.Status),
		Assessment:  nil,
		CreatedAt:   u.CreatedAt,
		UpdatedAt:   u.UpdatedAt,
		CompletedAt: u.CompletedAt,
	}

	if u.Assessment != nil {
		grammarIssues := make([]domain.UserVocabularyExerciseAssessmentGrammarIssue, len(u.Assessment.GrammarIssues))
		for i, issue := range u.Assessment.GrammarIssues {
			grammarIssues[i] = domain.UserVocabularyExerciseAssessmentGrammarIssue{
				Issue:      issue.Issue,
				Correction: issue.Correction,
			}
		}

		improvementSuggestions := make([]domain.UserVocabularyExerciseAssessmentImprovementSuggestion, len(u.Assessment.ImprovementSuggestions))
		for i, suggestion := range u.Assessment.ImprovementSuggestions {
			improvementSuggestions[i] = domain.UserVocabularyExerciseAssessmentImprovementSuggestion{
				Instruction: suggestion.Instruction,
				Example:     suggestion.Example,
			}
		}

		exercise.Assessment = &domain.UserVocabularyExerciseAssessment{
			IsVocabularyCorrect:    u.Assessment.IsVocabularyCorrect,
			VocabularyIssue:        u.Assessment.VocabularyIssue,
			IsTenseCorrect:         u.Assessment.IsTenseCorrect,
			TenseIssue:             u.Assessment.TenseIssue,
			GrammarIssues:          grammarIssues,
			ImprovementSuggestions: improvementSuggestions,
		}
	}
	return exercise
}

func (u UserVocabularyExercise) FromDomain(exercise domain.UserVocabularyExercise) (*UserVocabularyExercise, error) {
	id, err := primitive.ObjectIDFromHex(exercise.ID)
	if err != nil {
		return nil, apperrors.Common.InvalidID
	}

	uid, err := primitive.ObjectIDFromHex(exercise.UserID)
	if err != nil {
		return nil, apperrors.User.InvalidUserID
	}

	tid, err := primitive.ObjectIDFromHex(exercise.TermID)
	if err != nil {
		return nil, apperrors.Common.InvalidID
	}

	var assessment *UserVocabularyExerciseAssessment
	if exercise.Assessment != nil {
		grammarIssues := make([]UserVocabularyExerciseAssessmentGrammarIssue, len(exercise.Assessment.GrammarIssues))
		for i, issue := range exercise.Assessment.GrammarIssues {
			grammarIssues[i] = UserVocabularyExerciseAssessmentGrammarIssue{
				Issue:      issue.Issue,
				Correction: issue.Correction,
			}
		}

		improvementSuggestions := make([]UserVocabularyExerciseAssessmentImprovementSuggestion, len(exercise.Assessment.ImprovementSuggestions))
		for i, suggestion := range exercise.Assessment.ImprovementSuggestions {
			improvementSuggestions[i] = UserVocabularyExerciseAssessmentImprovementSuggestion{
				Instruction: suggestion.Instruction,
				Example:     suggestion.Example,
			}
		}

		assessment = &UserVocabularyExerciseAssessment{
			IsVocabularyCorrect:    exercise.Assessment.IsVocabularyCorrect,
			VocabularyIssue:        exercise.Assessment.VocabularyIssue,
			IsTenseCorrect:         exercise.Assessment.IsTenseCorrect,
			TenseIssue:             exercise.Assessment.TenseIssue,
			GrammarIssues:          grammarIssues,
			ImprovementSuggestions: improvementSuggestions,
		}
	}

	return &UserVocabularyExercise{
		ID:          id,
		UserID:      uid,
		TermID:      tid,
		Term:        exercise.Term,
		Language:    exercise.Language.String(),
		Tense:       exercise.Tense.String(),
		Content:     exercise.Content,
		Status:      exercise.Status.String(),
		Assessment:  assessment,
		CreatedAt:   exercise.CreatedAt,
		UpdatedAt:   exercise.UpdatedAt,
		CompletedAt: exercise.CompletedAt,
	}, nil
}
