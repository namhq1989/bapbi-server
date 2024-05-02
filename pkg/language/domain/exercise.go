package domain

import "strings"

type ExerciseStatus string

const (
	ExerciseStatusUnknown            ExerciseStatus = ""
	ExerciseStatusCompleted          ExerciseStatus = "completed"
	ExerciseStatusProgressing        ExerciseStatus = "progressing"
	ExerciseStatusCorrectionRequired ExerciseStatus = "correction_required"
)

func (s ExerciseStatus) String() string {
	switch s {
	case ExerciseStatusCompleted, ExerciseStatusProgressing, ExerciseStatusCorrectionRequired:
		return string(s)
	default:
		return ""
	}
}

func (s ExerciseStatus) IsValid() bool {
	return s != ExerciseStatusUnknown
}

func ToExerciseStatus(value string) ExerciseStatus {
	switch strings.ToLower(value) {
	case ExerciseStatusCompleted.String():
		return ExerciseStatusCompleted
	case ExerciseStatusProgressing.String():
		return ExerciseStatusProgressing
	case ExerciseStatusCorrectionRequired.String():
		return ExerciseStatusCorrectionRequired
	default:
		return ExerciseStatusUnknown
	}
}
