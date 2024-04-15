package domain

import (
	"github.com/namhq1989/bapbi-server/internal/queue"
	"github.com/namhq1989/bapbi-server/internal/utils/appcontext"
)

type QueueRepository interface {
	EnqueueUserCreated(ctx *appcontext.AppContext, user User) error
	EnqueueUserUpdated(ctx *appcontext.AppContext, user User) error
	EnqueueUserCreatedForHealth(ctx *appcontext.AppContext, user queue.User) error
}
