package model

import (
	"time"

	"github.com/namhq1989/bapbi-server/internal/database"
	apperrors "github.com/namhq1989/bapbi-server/internal/utils/error"
	"github.com/namhq1989/bapbi-server/pkg/user/domain"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID                   primitive.ObjectID `bson:"_id" json:"id"`
	Name                 string             `bson:"name" json:"name"`
	Email                string             `bson:"email" json:"email"`
	Status               string             `bson:"status" json:"status"`
	SubscriptionPlan     string             `bson:"subscriptionPlan" json:"subscriptionPlan"`
	SubscriptionExpireAt time.Time          `bson:"subscriptionExpireAt" json:"subscriptionExpireAt"`
	CreatedAt            time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedAt            time.Time          `bson:"updatedAt" json:"updatedAt"`
}

func (m User) ToDomain() domain.User {
	return domain.User{
		ID:                   m.ID.Hex(),
		Name:                 m.Name,
		Email:                m.Email,
		Status:               domain.ToUserStatus(m.Status),
		SubscriptionPlan:     domain.ToSubscriptionPlan(m.SubscriptionPlan),
		SubscriptionExpireAt: m.SubscriptionExpireAt,
		CreatedAt:            m.CreatedAt,
		UpdatedAt:            m.UpdatedAt,
	}
}

func (m User) FromDomain(user domain.User) (*User, error) {
	id, err := database.ObjectIDFromString(user.ID)
	if err != nil {
		return nil, apperrors.Common.InvalidID
	}

	return &User{
		ID:                   id,
		Name:                 user.Name,
		Email:                user.Email,
		Status:               user.Status.String(),
		SubscriptionPlan:     user.SubscriptionPlan.String(),
		SubscriptionExpireAt: user.SubscriptionExpireAt,
		CreatedAt:            user.CreatedAt,
		UpdatedAt:            user.UpdatedAt,
	}, nil
}
