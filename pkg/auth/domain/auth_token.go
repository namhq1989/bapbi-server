package domain

import (
	"time"

	"github.com/namhq1989/bapbi-server/internal/utils/appcontext"
)

type AuthTokenRepository interface {
	CreateAuthToken(ctx *appcontext.AppContext, token RefreshToken) error
	DeleteAuthToken(ctx *appcontext.AppContext, token RefreshToken) error
	FindAuthToken(ctx *appcontext.AppContext, userID, refreshToken string) (*RefreshToken, error)
}

type Tokens struct {
	AccessToken        string
	RefreshToken       string
	AccessTokenExpiry  time.Time
	RefreshTokenExpiry time.Time
}

type RefreshToken struct {
	ID     string
	UserID string
	Token  string
	Expiry time.Time
}
