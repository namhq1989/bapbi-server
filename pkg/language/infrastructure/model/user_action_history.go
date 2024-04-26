package model

import (
	"time"

	"github.com/namhq1989/bapbi-server/internal/database"
	apperrors "github.com/namhq1989/bapbi-server/internal/utils/error"
	"github.com/namhq1989/bapbi-server/pkg/language/domain"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserActionHistory struct {
	ID        primitive.ObjectID `bson:"_id"`
	UserID    primitive.ObjectID `bson:"userId"`
	Action    string             `bson:"action"`
	Data      string             `bson:"data"`
	CreatedAt time.Time          `bson:"createdAt"`
}

func (m UserActionHistory) ToDomain() domain.UserActionHistory {
	return domain.UserActionHistory{
		ID:        m.ID.Hex(),
		UserID:    m.UserID.Hex(),
		Action:    domain.ToUserActionType(m.Action),
		Data:      m.Data,
		CreatedAt: m.CreatedAt,
	}
}

func (m UserActionHistory) FromDomain(history domain.UserActionHistory) (*UserActionHistory, error) {
	id, err := database.ObjectIDFromString(history.ID)
	if err != nil {
		return nil, apperrors.Common.InvalidID
	}

	uid, err := database.ObjectIDFromString(history.UserID)
	if err != nil {
		return nil, apperrors.User.InvalidUserID
	}

	return &UserActionHistory{
		ID:        id,
		UserID:    uid,
		Action:    history.Action.String(),
		Data:      history.Data,
		CreatedAt: history.CreatedAt,
	}, nil
}
