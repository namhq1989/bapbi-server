package domain

import (
	"time"

	"github.com/namhq1989/bapbi-server/internal/database"
	"github.com/namhq1989/bapbi-server/internal/utils/appcontext"
	apperrors "github.com/namhq1989/bapbi-server/internal/utils/error"
)

type UserTermRepository interface {
	IsUserTermAdded(ctx *appcontext.AppContext, userID, termID string) (bool, error)
	CountTotalTermAddedByTimeRange(ctx *appcontext.AppContext, userID string, start, end time.Time) (int64, error)
	FindUserTermByID(ctx *appcontext.AppContext, termID string) (*UserTerm, error)
	AddUserTerm(ctx *appcontext.AppContext, term UserTerm) error
	UpdateUserTerm(ctx *appcontext.AppContext, term UserTerm) error
}

type UserTerm struct {
	ID          string
	UserID      string
	Term        string
	IsFavourite bool
	CreatedAt   time.Time
}

func NewUserTerm(userID, term string) (*UserTerm, error) {
	if term == "" {
		return nil, apperrors.Language.InvalidTerm
	}

	return &UserTerm{
		ID:          database.NewStringID(),
		UserID:      userID,
		Term:        term,
		IsFavourite: false,
		CreatedAt:   time.Now(),
	}, nil
}

func (d *UserTerm) SetIsFavourite(isFavourite bool) {
	d.IsFavourite = isFavourite
}