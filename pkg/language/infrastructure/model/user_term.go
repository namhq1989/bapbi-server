package model

import (
	"time"

	"github.com/namhq1989/bapbi-server/internal/database"
	apperrors "github.com/namhq1989/bapbi-server/internal/utils/error"
	"github.com/namhq1989/bapbi-server/pkg/language/domain"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserTerm struct {
	ID          primitive.ObjectID `bson:"_id"`
	UserID      primitive.ObjectID `bson:"userId"`
	Term        string             `bson:"term"`
	IsFavourite bool               `bson:"isFavourite"`
	CreatedAt   time.Time          `bson:"createdAt"`
}

func (d UserTerm) ToDomain() domain.UserTerm {
	return domain.UserTerm{
		ID:          d.ID.Hex(),
		UserID:      d.UserID.Hex(),
		Term:        d.Term,
		IsFavourite: d.IsFavourite,
		CreatedAt:   d.CreatedAt,
	}
}

func (d UserTerm) FromDomain(term domain.UserTerm) (*UserTerm, error) {
	id, err := database.ObjectIDFromString(term.ID)
	if err != nil {
		return nil, apperrors.Common.InvalidID
	}

	uid, err := database.ObjectIDFromString(term.UserID)
	if err != nil {
		return nil, apperrors.User.InvalidUserID
	}

	return &UserTerm{
		ID:          id,
		UserID:      uid,
		Term:        term.Term,
		IsFavourite: term.IsFavourite,
		CreatedAt:   time.Now(),
	}, nil
}
