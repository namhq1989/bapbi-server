package domain

import "github.com/namhq1989/bapbi-server/internal/utils/appcontext"

type UserHub interface {
	GetOneByID(ctx *appcontext.AppContext, id string) (*User, error)
	GetOneByEmail(ctx *appcontext.AppContext, email string) (*User, error)
	CreateUser(ctx *appcontext.AppContext, user User) (string, error)
}

type User struct {
	ID    string
	Name  string
	Email string
}
