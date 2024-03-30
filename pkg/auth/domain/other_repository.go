package domain

import "github.com/namhq1989/bapbi-server/internal/utils/appcontext"

type JwtRepository interface {
	GenerateAccessToken(ctx *appcontext.AppContext, userID string) (string, error)
	GenerateTokens(ctx *appcontext.AppContext, userID string) (*Tokens, error)
}
