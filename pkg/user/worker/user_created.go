package worker

import (
	"context"

	"github.com/namhq1989/bapbi-server/internal/queue"

	"github.com/goccy/go-json"
	"github.com/hibiken/asynq"
	"github.com/namhq1989/bapbi-server/internal/utils/appcontext"
	"github.com/namhq1989/bapbi-server/pkg/user/domain"
)

func (w Workers) UserCreated(bgCtx context.Context, t *asynq.Task) error {
	var (
		ctx  = appcontext.New(bgCtx)
		user domain.User
	)

	ctx.Logger().Info("[worker] process new task", appcontext.Fields{"type": t.Type(), "payload": string(t.Payload())})

	ctx.Logger().Text("unmarshal task payload")
	if err := json.Unmarshal(t.Payload(), &user); err != nil {
		ctx.Logger().Error("failed to unmarshal task payload", err, appcontext.Fields{})
		return err
	}

	// add queue
	ctx.Logger().Text("add health queue task")
	if err := w.queueRepository.EnqueueUserCreatedForHealth(ctx, queue.User{ID: user.ID}); err != nil {
		return err
	}

	ctx.Logger().Info("[worker] done task", appcontext.Fields{"type": t.Type()})
	return nil
}
