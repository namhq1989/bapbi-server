package domain

import (
	"time"

	"github.com/namhq1989/bapbi-server/internal/utils/appcontext"

	"github.com/namhq1989/bapbi-server/internal/database"
	apperrors "github.com/namhq1989/bapbi-server/internal/utils/error"
)

type UserSearchHistoryRepository interface {
	CreateUserSearchHistory(ctx *appcontext.AppContext, history UserSearchHistory) error
	CountTotalSearchedByTimeRange(ctx *appcontext.AppContext, userID string, start, end time.Time) (int64, error)
}

type UserSearchHistory struct {
	ID        string
	UserID    string
	Term      string
	IsValid   bool
	CreatedAt time.Time
}

func NewUserSearchHistory(userID, term string, isValid bool) (*UserSearchHistory, error) {
	if term == "" {
		return nil, apperrors.Language.InvalidTerm
	}

	return &UserSearchHistory{
		ID:        database.NewStringID(),
		UserID:    userID,
		Term:      term,
		IsValid:   isValid,
		CreatedAt: time.Now(),
	}, nil
}
