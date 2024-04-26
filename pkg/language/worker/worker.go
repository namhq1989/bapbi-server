package worker

import (
	"context"
	"fmt"

	"github.com/namhq1989/bapbi-server/internal/queue"
	"github.com/namhq1989/bapbi-server/internal/utils/appcontext"
	"github.com/namhq1989/bapbi-server/pkg/language/domain"
)

type Workers struct {
	queue             *queue.Queue
	termRepository    domain.TermRepository
	openaiRepository  domain.OpenAIRepository
	scraperRepository domain.ScraperRepository
}

func New(
	queue *queue.Queue,
	termRepository domain.TermRepository,
	openaiRepository domain.OpenAIRepository,
	scraperRepository domain.ScraperRepository,
) Workers {
	return Workers{
		queue:             queue,
		termRepository:    termRepository,
		openaiRepository:  openaiRepository,
		scraperRepository: scraperRepository,
	}
}

type cronjobData struct {
	Task       string      `json:"task"`
	CronSpec   string      `json:"cronSpec"`
	Payload    interface{} `json:"payload"`
	RetryTimes int         `json:"retryTimes"`
}

func (w Workers) Start() {
	// cron jobs
	w.addCronjob()

	w.queue.Server.HandleFunc(w.queue.GenerateTypename(queue.TypeNames.Language.FeaturedWord), w.FeaturedWord)
}

func (w Workers) addCronjob() {
	var (
		ctx  = appcontext.New(context.Background())
		jobs = []cronjobData{
			{
				Task:       w.queue.GenerateTypename(queue.TypeNames.Language.FeaturedWord),
				CronSpec:   "@every 8h",
				Payload:    nil,
				RetryTimes: 3,
			},
		}
	)

	for _, job := range jobs {
		entryID, err := w.queue.ScheduleTask(job.Task, job.Payload, job.CronSpec, job.RetryTimes)
		if err != nil {
			ctx.Logger().Error("error when initializing cronjob", err, appcontext.Fields{})
			panic(err)
		}

		ctx.Logger().Info(fmt.Sprintf("[cronjob] cronjob '%s' initialize successfully with cronSpec '%s' and retryTimes '%d'", job.Task, job.CronSpec, job.RetryTimes), appcontext.Fields{
			"entryId": entryID,
		})
	}
}