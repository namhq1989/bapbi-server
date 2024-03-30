package domain

import "github.com/namhq1989/bapbi-server/internal/utils/appcontext"

type UserHub interface {
	FindOneByID(ctx *appcontext.AppContext, id string) (*User, error)
	FindOneByEmail(ctx *appcontext.AppContext, email string) (*User, error)
	CreateUser(ctx *appcontext.AppContext, user User) (*User, error)
}
