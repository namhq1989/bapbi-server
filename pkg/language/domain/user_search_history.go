package domain

import (
	"time"

	"github.com/namhq1989/bapbi-server/internal/utils/appcontext"

	"github.com/namhq1989/bapbi-server/internal/database"
	apperrors "github.com/namhq1989/bapbi-server/internal/utils/error"
)

type UserSearchHistoryRepository interface {
	CreateUserSearchHistory(ctx *appcontext.AppContext, history UserSearchHistory) error
}

type UserSearchHistory struct {
	ID          string
	UserID      string
	Term        string
	IsValid     bool
	IsFavourite bool
	CreatedAt   time.Time
}

func NewUserSearchHistory(userID, term string, isValid bool) (*UserSearchHistory, error) {
	if term == "" {
		return nil, apperrors.Language.InvalidTerm
	}

	return &UserSearchHistory{
		ID:          database.NewStringID(),
		UserID:      userID,
		Term:        term,
		IsValid:     isValid,
		IsFavourite: false,
		CreatedAt:   time.Now(),
	}, nil
}
