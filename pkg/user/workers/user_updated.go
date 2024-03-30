package workers

import (
	"context"

	"github.com/goccy/go-json"
	"github.com/hibiken/asynq"
	"github.com/namhq1989/bapbi-server/internal/utils/appcontext"
	"github.com/namhq1989/bapbi-server/pkg/user/domain"
)

func (w Workers) UserUpdated(bgCtx context.Context, t *asynq.Task) error {
	var (
		ctx  = appcontext.New(bgCtx)
		user domain.User
	)

	ctx.Logger().Info("[worker] process new task", appcontext.Fields{"type": t.Type()})

	ctx.Logger().Info("unmarshal task payload", appcontext.Fields{})
	if err := json.Unmarshal(t.Payload(), &user); err != nil {
		ctx.Logger().Error("failed to unmarshal task payload", err, appcontext.Fields{})
		return err
	}

	ctx.Logger().Text("*** CURRENTLY DO NOTHING WITH THIS TASK ***")

	ctx.Logger().Info("[worker] done task", appcontext.Fields{"type": t.Type()})
	return nil
}
