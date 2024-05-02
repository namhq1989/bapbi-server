package dto

import (
	"github.com/namhq1989/bapbi-server/internal/utils/httprespond"
	"github.com/namhq1989/bapbi-server/pkg/language/domain"
)

type SubmitUserTermExerciseRequest struct {
	ExerciseID string `json:"exerciseId" validate:"required" message:"language_invalid_exercise_id"`
	Content    string `json:"content" validate:"required" message:"language_invalid_vocabulary_exercise_data"`
}

type SubmitUserTermExerciseResponse struct {
	Status      string                     `json:"status"`
	CompletedAt *httprespond.TimeResponse  `json:"completedAt"`
	Assessment  UserTermExerciseAssessment `json:"assessment"`
}

type UserTermExerciseAssessmentGrammarIssue struct {
	Issue      string `json:"issue"`
	Correction string `json:"correction"`
}

type UserTermExerciseAssessmentImprovementSuggestion struct {
	Instruction string `json:"instruction"`
	Example     string `json:"example"`
}

type UserTermExerciseAssessment struct {
	IsVocabularyCorrect    bool                                              `json:"isVocabularyCorrect"`
	VocabularyIssue        string                                            `json:"vocabularyIssue"`
	IsTenseCorrect         bool                                              `json:"isTenseCorrect"`
	TenseIssue             string                                            `json:"tenseIssue"`
	GrammarIssues          []UserTermExerciseAssessmentGrammarIssue          `json:"grammarIssues"`
	ImprovementSuggestions []UserTermExerciseAssessmentImprovementSuggestion `json:"improvementSuggestions"`
}

func (u UserTermExerciseAssessment) FromDomain(assessment domain.UserTermExerciseAssessment) UserTermExerciseAssessment {
	grammarIssues := make([]UserTermExerciseAssessmentGrammarIssue, len(assessment.GrammarIssues))
	for i, issue := range assessment.GrammarIssues {
		grammarIssues[i] = UserTermExerciseAssessmentGrammarIssue{
			Issue:      issue.Issue,
			Correction: issue.Correction,
		}
	}

	improvementSuggestions := make([]UserTermExerciseAssessmentImprovementSuggestion, len(assessment.ImprovementSuggestions))
	for i, suggestion := range assessment.ImprovementSuggestions {
		improvementSuggestions[i] = UserTermExerciseAssessmentImprovementSuggestion{
			Instruction: suggestion.Instruction,
			Example:     suggestion.Example,
		}
	}

	return UserTermExerciseAssessment{
		IsVocabularyCorrect:    assessment.IsVocabularyCorrect,
		VocabularyIssue:        assessment.VocabularyIssue,
		IsTenseCorrect:         assessment.IsTenseCorrect,
		TenseIssue:             assessment.TenseIssue,
		GrammarIssues:          grammarIssues,
		ImprovementSuggestions: improvementSuggestions,
	}
}
