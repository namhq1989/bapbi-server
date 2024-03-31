package domain

import (
	"time"

	apperrors "github.com/namhq1989/bapbi-server/internal/utils/error"
	"github.com/namhq1989/bapbi-server/internal/utils/validation"
)

type UserRepository interface{}

type User struct {
	ID        string
	Name      string
	Email     string
	Status    UserStatus
	CreatedAt time.Time
	UpdatedAt time.Time
}

func CreateUser(name, email string) (*User, error) {
	if !validation.IsValidUserName(name) {
		return nil, apperrors.Common.InvalidName
	}

	if !validation.IsValidEmail(email) {
		return nil, apperrors.Common.InvalidEmail
	}

	return &User{
		ID:        "",
		Name:      name,
		Email:     email,
		Status:    UserStatusActive,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}, nil
}
