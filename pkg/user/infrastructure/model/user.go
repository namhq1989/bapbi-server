package model

import (
	"time"

	"github.com/namhq1989/bapbi-server/internal/database"
	apperrors "github.com/namhq1989/bapbi-server/internal/utils/error"
	"github.com/namhq1989/bapbi-server/pkg/user/domain"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID        primitive.ObjectID `bson:"_id" json:"id"`
	Name      string             `bson:"name" json:"name"`
	Email     string             `bson:"email" json:"email"`
	Status    string             `bson:"status" json:"status"`
	CreatedAt time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedAt time.Time          `bson:"updatedAt" json:"updatedAt"`
}

func (m User) ToDomain() domain.User {
	return domain.User{
		ID:        m.ID.Hex(),
		Name:      m.Name,
		Email:     m.Email,
		Status:    domain.ToUserStatus(m.Status),
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
	}
}

func (m User) FromDomain(user domain.User) (*User, error) {
	var id primitive.ObjectID

	if user.ID == "" {
		id = database.NewObjectID()
	} else {
		oid, err := database.ObjectIDFromString(user.ID)
		if err != nil {
			return nil, apperrors.Common.InvalidID
		}

		id = oid
	}

	return &User{
		ID:        id,
		Name:      user.Name,
		Email:     user.Email,
		Status:    user.Status.String(),
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}, nil
}
