package workers

import (
	"github.com/namhq1989/bapbi-server/internal/queue"
	"github.com/namhq1989/bapbi-server/pkg/health/domain"
)

type Workers struct {
	queue                          *queue.Queue
	hydrationProfileRepository     domain.HydrationProfileRepository
	hydrationDailyReportRepository domain.HydrationDailyReportRepository
}

func New(
	queue *queue.Queue,
	hydrationProfileRepository domain.HydrationProfileRepository,
	hydrationDailyReportRepository domain.HydrationDailyReportRepository,
) Workers {
	return Workers{
		queue:                          queue,
		hydrationProfileRepository:     hydrationProfileRepository,
		hydrationDailyReportRepository: hydrationDailyReportRepository,
	}
}

func (w Workers) Start() {
	w.queue.Server.HandleFunc(w.queue.GenerateTypename(domain.QueueTypeNames.NewWaterIntakeLog), w.NewWaterIntakeLog)
}
