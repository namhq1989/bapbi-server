package workers

import (
	"github.com/namhq1989/bapbi-server/internal/queue"
	"github.com/namhq1989/bapbi-server/pkg/health/domain"
)

type Workers struct {
	queue                          *queue.Queue
	healthProfileRepository        domain.HealthProfileRepository
	hydrationProfileRepository     domain.HydrationProfileRepository
	hydrationDailyReportRepository domain.HydrationDailyReportRepository
}

func New(
	queue *queue.Queue,
	healthProfileRepository domain.HealthProfileRepository,
	hydrationProfileRepository domain.HydrationProfileRepository,
	hydrationDailyReportRepository domain.HydrationDailyReportRepository,
) Workers {
	return Workers{
		queue:                          queue,
		healthProfileRepository:        healthProfileRepository,
		hydrationProfileRepository:     hydrationProfileRepository,
		hydrationDailyReportRepository: hydrationDailyReportRepository,
	}
}

func (w Workers) Start() {
	w.queue.Server.HandleFunc(w.queue.GenerateTypename(queue.TypeNames.Health.UserCreated), w.UserCreated)
	w.queue.Server.HandleFunc(w.queue.GenerateTypename(queue.TypeNames.Health.NewWaterIntakeLog), w.NewWaterIntakeLog)
}
