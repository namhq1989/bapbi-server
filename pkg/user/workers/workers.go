package workers

import (
	"github.com/namhq1989/bapbi-server/internal/queue"
	"github.com/namhq1989/bapbi-server/pkg/user/domain"
)

type Workers struct {
	queue *queue.Queue
}

func New(queue *queue.Queue) Workers {
	return Workers{
		queue: queue,
	}
}

func (w Workers) Start() {
	w.queue.Server.HandleFunc(w.queue.GenerateTypename(domain.QueueTypeNames.UserCreated), w.UserCreated)
	w.queue.Server.HandleFunc(w.queue.GenerateTypename(domain.QueueTypeNames.UserUpdated), w.UserUpdated)
}
