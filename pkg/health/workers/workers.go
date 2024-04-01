package workers

import (
	"github.com/namhq1989/bapbi-server/internal/queue"
	"github.com/namhq1989/bapbi-server/pkg/health/domain"
)

type Workers struct {
	queue                          *queue.Queue
	drinkWaterProfileRepository    domain.DrinkWaterProfileRepository
	dailyHydrationReportRepository domain.DailyHydrationReportRepository
}

func New(queue *queue.Queue, drinkWaterProfileRepository domain.DrinkWaterProfileRepository, dailyHydrationReportRepository domain.DailyHydrationReportRepository) Workers {
	return Workers{
		queue:                          queue,
		drinkWaterProfileRepository:    drinkWaterProfileRepository,
		dailyHydrationReportRepository: dailyHydrationReportRepository,
	}
}

func (w Workers) Start() {
	w.queue.Server.HandleFunc(w.queue.GenerateTypename(domain.QueueTypeNames.NewWaterIntakeLog), w.NewWaterIntakeLog)
}
