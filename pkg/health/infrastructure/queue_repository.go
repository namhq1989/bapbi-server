package infrastructure

import (
	"github.com/namhq1989/bapbi-server/internal/queue"
	"github.com/namhq1989/bapbi-server/internal/utils/appcontext"
	"github.com/namhq1989/bapbi-server/pkg/health/domain"
)

type QueueRepository struct {
	queue *queue.Queue
}

func NewQueueRepository(queue *queue.Queue) QueueRepository {
	return QueueRepository{
		queue: queue,
	}
}

func (r QueueRepository) EnqueueNewWaterIntakeLog(ctx *appcontext.AppContext, log domain.WaterIntakeLog) error {
	typename := r.queue.GenerateTypename(queue.TypeNames.Health.NewWaterIntakeLog)
	t, err := r.queue.RunTask(typename, log, -1)
	if err != nil {
		ctx.Logger().Error("failed to enqueue task", err, appcontext.Fields{"typename": typename, "log": log})
		return err
	}

	ctx.Logger().Info("enqueued task", appcontext.Fields{"taskId": t.ID, "typename": typename})
	return nil
}
