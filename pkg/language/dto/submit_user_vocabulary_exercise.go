package dto

import (
	"github.com/namhq1989/bapbi-server/internal/utils/httprespond"
	"github.com/namhq1989/bapbi-server/pkg/language/domain"
)

type SubmitUserVocabularyExerciseRequest struct {
	ExerciseID string `json:"exerciseId" validate:"required" message:"language_invalid_exercise_id"`
	Content    string `json:"content" validate:"required" message:"language_invalid_vocabulary_exercise_data"`
}

type SubmitUserVocabularyExerciseResponse struct {
	Status      string                           `json:"status"`
	CompletedAt *httprespond.TimeResponse        `json:"completedAt"`
	Assessment  UserVocabularyExerciseAssessment `json:"assessment"`
}

type UserVocabularyExerciseAssessmentGrammarIssue struct {
	Issue      string `json:"issue"`
	Correction string `json:"correction"`
}

type UserVocabularyExerciseAssessmentImprovementSuggestion struct {
	Instruction string `json:"instruction"`
	Example     string `json:"example"`
}

type UserVocabularyExerciseAssessment struct {
	IsVocabularyCorrect    bool                                                    `json:"isVocabularyCorrect"`
	VocabularyIssue        string                                                  `json:"vocabularyIssue"`
	IsTenseCorrect         bool                                                    `json:"isTenseCorrect"`
	TenseIssue             string                                                  `json:"tenseIssue"`
	GrammarIssues          []UserVocabularyExerciseAssessmentGrammarIssue          `json:"grammarIssues"`
	ImprovementSuggestions []UserVocabularyExerciseAssessmentImprovementSuggestion `json:"improvementSuggestions"`
}

func (u UserVocabularyExerciseAssessment) FromDomain(assessment domain.UserVocabularyExerciseAssessment) UserVocabularyExerciseAssessment {
	grammarIssues := make([]UserVocabularyExerciseAssessmentGrammarIssue, len(assessment.GrammarIssues))
	for i, issue := range assessment.GrammarIssues {
		grammarIssues[i] = UserVocabularyExerciseAssessmentGrammarIssue{
			Issue:      issue.Issue,
			Correction: issue.Correction,
		}
	}

	improvementSuggestions := make([]UserVocabularyExerciseAssessmentImprovementSuggestion, len(assessment.ImprovementSuggestions))
	for i, suggestion := range assessment.ImprovementSuggestions {
		improvementSuggestions[i] = UserVocabularyExerciseAssessmentImprovementSuggestion{
			Instruction: suggestion.Instruction,
			Example:     suggestion.Example,
		}
	}

	return UserVocabularyExerciseAssessment{
		IsVocabularyCorrect:    assessment.IsVocabularyCorrect,
		VocabularyIssue:        assessment.VocabularyIssue,
		IsTenseCorrect:         assessment.IsTenseCorrect,
		TenseIssue:             assessment.TenseIssue,
		GrammarIssues:          grammarIssues,
		ImprovementSuggestions: improvementSuggestions,
	}
}
