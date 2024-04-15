package infrastructure

import (
	"github.com/namhq1989/bapbi-server/internal/queue"
	"github.com/namhq1989/bapbi-server/internal/utils/appcontext"
	"github.com/namhq1989/bapbi-server/pkg/user/domain"
)

type QueueRepository struct {
	queue *queue.Queue
}

func NewQueueRepository(queue *queue.Queue) QueueRepository {
	return QueueRepository{
		queue: queue,
	}
}

func (r QueueRepository) EnqueueUserCreated(ctx *appcontext.AppContext, user domain.User) error {
	typename := r.queue.GenerateTypename(queue.TypeNames.User.UserCreated)
	t, err := r.queue.RunTask(typename, user, -1)
	if err != nil {
		ctx.Logger().Error("failed to enqueue task", err, appcontext.Fields{"typename": typename, "user": user})
		return err
	}

	ctx.Logger().Info("enqueued task", appcontext.Fields{"taskId": t.ID, "typename": typename})
	return nil
}

func (r QueueRepository) EnqueueUserUpdated(ctx *appcontext.AppContext, user domain.User) error {
	typename := r.queue.GenerateTypename(queue.TypeNames.User.UserUpdated)
	t, err := r.queue.RunTask(typename, user, -1)
	if err != nil {
		ctx.Logger().Error("failed to enqueue task", err, appcontext.Fields{"typename": typename, "user": user})
		return err
	}

	ctx.Logger().Info("enqueued task", appcontext.Fields{"taskId": t.ID, "typename": typename})
	return nil
}

func (r QueueRepository) EnqueueUserCreatedForHealth(ctx *appcontext.AppContext, user queue.User) error {
	typename := r.queue.GenerateTypename(queue.TypeNames.Health.UserCreated)
	t, err := r.queue.RunTask(typename, user, -1)
	if err != nil {
		ctx.Logger().Error("failed to enqueue task", err, appcontext.Fields{"typename": typename, "user": user})
		return err
	}

	ctx.Logger().Info("enqueued task", appcontext.Fields{"taskId": t.ID, "typename": typename})
	return nil
}
