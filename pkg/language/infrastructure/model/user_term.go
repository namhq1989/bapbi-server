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
	TermID      primitive.ObjectID `bson:"termId"`
	Term        string             `bson:"term"`
	IsFavourite bool               `bson:"isFavourite"`
	CreatedAt   time.Time          `bson:"createdAt"`
}

func (m UserTerm) ToDomain() domain.UserTerm {
	return domain.UserTerm{
		ID:          m.ID.Hex(),
		UserID:      m.UserID.Hex(),
		TermID:      m.TermID.Hex(),
		Term:        m.Term,
		IsFavourite: m.IsFavourite,
		CreatedAt:   m.CreatedAt,
	}
}

func (UserTerm) FromDomain(term domain.UserTerm) (*UserTerm, error) {
	id, err := database.ObjectIDFromString(term.ID)
	if err != nil {
		return nil, apperrors.Common.InvalidID
	}

	uid, err := database.ObjectIDFromString(term.UserID)
	if err != nil {
		return nil, apperrors.User.InvalidUserID
	}

	tid, err := database.ObjectIDFromString(term.TermID)
	if err != nil {
		return nil, apperrors.Common.InvalidID
	}

	return &UserTerm{
		ID:          id,
		UserID:      uid,
		TermID:      tid,
		Term:        term.Term,
		IsFavourite: term.IsFavourite,
		CreatedAt:   term.CreatedAt,
	}, nil
}
