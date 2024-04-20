package domain

import (
	"time"

	"github.com/namhq1989/bapbi-server/internal/database"

	apperrors "github.com/namhq1989/bapbi-server/internal/utils/error"
	"github.com/namhq1989/bapbi-server/internal/utils/validation"
)

type UserRepository interface{}

type User struct {
	ID                   string
	Name                 string
	Email                string
	Status               UserStatus
	SubscriptionPlan     SubscriptionPlan
	SubscriptionExpireAt time.Time
	CreatedAt            time.Time
	UpdatedAt            time.Time
}

func CreateUser(name, email string) (*User, error) {
	if !validation.IsValidUserName(name) {
		return nil, apperrors.Common.InvalidName
	}

	if !validation.IsValidEmail(email) {
		return nil, apperrors.Common.InvalidEmail
	}

	return &User{
		ID:                   database.NewStringID(),
		Name:                 name,
		Email:                email,
		Status:               UserStatusActive,
		SubscriptionPlan:     SubscriptionPlanFree,
		SubscriptionExpireAt: time.Time{},
		CreatedAt:            time.Now(),
		UpdatedAt:            time.Now(),
	}, nil
}
