package domain

import (
	"github.com/namhq1989/bapbi-server/internal/utils/appcontext"
)

type SSORepository interface {
	GetUserDataWithGoogleToken(ctx *appcontext.AppContext, token string) (*SSOGoogleUser, error)
}

type SSOGoogle struct {
	Token string
}

type SSOGoogleUser struct {
	ID    string
	Email string
	Name  string
}
