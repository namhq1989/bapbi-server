package domain

import "github.com/namhq1989/bapbi-server/internal/utils/appcontext"

type AuthHub interface {
	IsAdmin(ctx *appcontext.AppContext, userID string) (bool, error)
}
