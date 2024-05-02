package domain

import (
	"strings"
	"time"

	"github.com/goccy/go-json"

	"github.com/namhq1989/bapbi-server/internal/utils/appcontext"

	"github.com/namhq1989/bapbi-server/internal/database"
	apperrors "github.com/namhq1989/bapbi-server/internal/utils/error"
)

type UserActionHistoryRepository interface {
	CreateUserActionHistory(ctx *appcontext.AppContext, history UserActionHistory) error
	CountTotalActionsByTimeRange(ctx *appcontext.AppContext, userID string, start, end time.Time) (int64, error)
}

//
// USER ACTION TYPE
//

type UserActionType string

const (
	UserActionTypeUnknown                  UserActionType = ""
	UserActionTypeSearchTerm               UserActionType = "search_term"
	UserActionTypeSubmitWritingExercise    UserActionType = "submit_writing_exercise"
	UserActionTypeSubmitVocabularyExercise UserActionType = "submit_vocabulary_exercise"
)

func (s UserActionType) String() string {
	switch s {
	case UserActionTypeSearchTerm, UserActionTypeSubmitWritingExercise, UserActionTypeSubmitVocabularyExercise:
		return string(s)
	default:
		return ""
	}
}

func (s UserActionType) IsValid() bool {
	return s != UserActionTypeUnknown
}

func ToUserActionType(value string) UserActionType {
	switch strings.ToLower(value) {
	case UserActionTypeSearchTerm.String():
		return UserActionTypeSearchTerm
	case UserActionTypeSubmitWritingExercise.String():
		return UserActionTypeSubmitWritingExercise
	case UserActionTypeSubmitVocabularyExercise.String():
		return UserActionTypeSubmitVocabularyExercise
	default:
		return UserActionTypeUnknown
	}
}

//
// USER ACTION HISTORY
//

type UserActionHistory struct {
	ID        string
	UserID    string
	Action    UserActionType
	Data      string
	CreatedAt time.Time
}

func NewUserActionHistory(userID string, act string) (*UserActionHistory, error) {
	action := ToUserActionType(act)
	if !action.IsValid() {
		return nil, apperrors.Language.InvalidUserAction
	}

	return &UserActionHistory{
		ID:        database.NewStringID(),
		UserID:    userID,
		Action:    action,
		Data:      "",
		CreatedAt: time.Now(),
	}, nil
}

func (d *UserActionHistory) SetData(data UserActionHistoryData) {
	b, _ := json.Marshal(data)
	d.Data = string(b)
}

//
// USER ACTION HISTORY DATA
//

type UserActionHistoryData struct {
	Term       string `json:"term,omitempty"`
	IsValid    bool   `json:"isValid,omitempty"`
	ExerciseID string `json:"exerciseId,omitempty"`
}
