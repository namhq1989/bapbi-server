package workers

import (
	"context"

	"github.com/goccy/go-json"
	"github.com/hibiken/asynq"
	"github.com/namhq1989/bapbi-server/internal/utils/appcontext"
	"github.com/namhq1989/bapbi-server/pkg/health/domain"
)

func (w Workers) NewWaterIntakeLog(bgCtx context.Context, t *asynq.Task) error {
	var (
		ctx = appcontext.New(bgCtx)
		log domain.WaterIntakeLog
	)

	ctx.Logger().Info("[worker] process new task", appcontext.Fields{"type": t.Type(), "payload": string(t.Payload())})

	ctx.Logger().Info("unmarshal task payload", appcontext.Fields{})
	if err := json.Unmarshal(t.Payload(), &log); err != nil {
		ctx.Logger().Error("failed to unmarshal task payload", err, appcontext.Fields{})
		return err
	}

	ctx.Logger().Text("*** CURRENTLY DO NOTHING WITH THIS TASK ***")
	// 1. calculate today's intake amount
	// increase the total amount
	// if exceed over daily's amount, set achieved to true
	// note: need a model for daily water intake summart

	ctx.Logger().Info("[worker] done task", appcontext.Fields{"type": t.Type()})
	return nil
}
