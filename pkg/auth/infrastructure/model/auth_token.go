package model

import (
	"time"

	"github.com/namhq1989/bapbi-server/internal/database"
	apperrors "github.com/namhq1989/bapbi-server/internal/utils/error"
	"github.com/namhq1989/bapbi-server/pkg/auth/domain"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AuthToken struct {
	ID           primitive.ObjectID `bson:"_id"`
	User         primitive.ObjectID `bson:"user"`
	RefreshToken string             `bson:"refreshToken"`
	Expiry       time.Time          `bson:"expiry"`
}

func (m AuthToken) ToDomain() domain.RefreshToken {
	return domain.RefreshToken{
		ID:     m.ID.Hex(),
		UserID: m.User.Hex(),
		Token:  m.RefreshToken,
		Expiry: m.Expiry,
	}
}

func (m AuthToken) FromDomain(token domain.RefreshToken) (*AuthToken, error) {
	var (
		id primitive.ObjectID
	)

	// id
	if token.ID == "" {
		id = database.NewObjectID()
	} else {
		oid, err := database.ObjectIDFromString(token.ID)
		if err != nil {
			return nil, apperrors.Common.InvalidID
		}

		id = oid
	}

	// user id
	userID, err := database.ObjectIDFromString(token.UserID)
	if err != nil {
		return nil, apperrors.User.UserNotFound
	}

	return &AuthToken{
		ID:           id,
		User:         userID,
		RefreshToken: token.Token,
		Expiry:       token.Expiry,
	}, nil
}
