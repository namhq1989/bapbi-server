package workers

import (
	"github.com/namhq1989/bapbi-server/internal/queue"
	"github.com/namhq1989/bapbi-server/pkg/user/domain"
)

type Workers struct {
	queue           *queue.Queue
	queueRepository domain.QueueRepository
}

func New(queue *queue.Queue, queueRepository domain.QueueRepository) Workers {
	return Workers{
		queue:           queue,
		queueRepository: queueRepository,
	}
}

func (w Workers) Start() {
	w.queue.Server.HandleFunc(w.queue.GenerateTypename(queue.TypeNames.User.UserCreated), w.UserCreated)
	w.queue.Server.HandleFunc(w.queue.GenerateTypename(queue.TypeNames.User.UserUpdated), w.UserUpdated)
}
