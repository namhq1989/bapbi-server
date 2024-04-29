package domain

import "github.com/namhq1989/bapbi-server/internal/utils/appcontext"

type QueueRepository interface {
	EnqueueNewUserTerm(ctx *appcontext.AppContext, userTerm UserTerm) error
}
