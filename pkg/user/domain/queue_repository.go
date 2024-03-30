package domain

import (
	"github.com/namhq1989/bapbi-server/internal/utils/appcontext"
)

var QueueTypeNames = struct {
	UserCreated string
	UserUpdated string
}{
	UserCreated: "user:user.created",
	UserUpdated: "user:user.updated",
}

type QueueRepository interface {
	EnqueueUserCreated(ctx *appcontext.AppContext, user User) error
	EnqueueUserUpdated(ctx *appcontext.AppContext, user User) error
}
