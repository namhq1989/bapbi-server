package model

import (
	"time"

	"github.com/namhq1989/bapbi-server/internal/database"
	apperrors "github.com/namhq1989/bapbi-server/internal/utils/error"
	"github.com/namhq1989/bapbi-server/pkg/language/domain"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserSearchHistory struct {
	ID          primitive.ObjectID `bson:"_id"`
	UserID      primitive.ObjectID `bson:"userId"`
	Term        string             `bson:"term"`
	IsValid     bool               `bson:"isValid"`
	IsFavourite bool               `bson:"isFavourite"`
	CreatedAt   time.Time          `bson:"createdAt"`
}

func (m UserSearchHistory) ToDomain() domain.UserSearchHistory {
	return domain.UserSearchHistory{
		ID:          m.ID.Hex(),
		UserID:      m.UserID.Hex(),
		Term:        m.Term,
		IsValid:     m.IsValid,
		IsFavourite: m.IsFavourite,
		CreatedAt:   m.CreatedAt,
	}
}

func (m UserSearchHistory) FromDomain(history domain.UserSearchHistory) (*UserSearchHistory, error) {
	id, err := database.ObjectIDFromString(history.ID)
	if err != nil {
		return nil, apperrors.Common.InvalidID
	}

	uid, err := database.ObjectIDFromString(history.UserID)
	if err != nil {
		return nil, apperrors.User.InvalidUserID
	}

	return &UserSearchHistory{
		ID:          id,
		UserID:      uid,
		Term:        history.Term,
		IsValid:     history.IsValid,
		IsFavourite: history.IsFavourite,
		CreatedAt:   history.CreatedAt,
	}, nil
}
