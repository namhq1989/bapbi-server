package queue

import (
	"fmt"
	"time"

	"github.com/goccy/go-json"

	"github.com/hibiken/asynq"
)

const (
	typenamePrefix = "app"
	queueDefault   = "default"
	queueCronjob   = "cronjob"

	taskTimeout   time.Duration = 30 * time.Second
	taskRetention               = 24 * 7 * time.Hour
)

func (Queue) GenerateTypename(name string) string {
	return fmt.Sprintf("%s:%s", typenamePrefix, name)
}

// RunTask ...
func (q Queue) RunTask(queue string, payload interface{}, retryTimes int) (*asynq.TaskInfo, error) {
	b, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal payload: %s", err.Error())
	}

	// create task and options
	task := asynq.NewTask(queue, b)
	options := make([]asynq.Option, 0)

	// retry times
	if retryTimes < 0 {
		retryTimes = 0
	}

	// append options
	options = append(options,
		asynq.MaxRetry(retryTimes),
		asynq.Timeout(taskTimeout),
		asynq.Retention(taskRetention),
	)

	// enqueue task
	return q.Client.Enqueue(task, options...)
}

// ScheduleTask create new task and run at specific time
// cronSpec follow cron expression
// https://www.freeformatter.com/cron-expression-generator-quartz.html
func (q Queue) ScheduleTask(typename string, payload interface{}, cronSpec string, retryTimes int) (string, error) {
	b, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("failed to marshal payload: %s", err.Error())
	}

	// create task and options
	task := asynq.NewTask(typename, b)
	options := make([]asynq.Option, 0)

	// retry times
	if retryTimes < 0 {
		retryTimes = 0
	}

	// append options
	options = append(options,
		asynq.Queue(queueCronjob),
		asynq.MaxRetry(retryTimes),
		asynq.Timeout(taskTimeout),
		asynq.Retention(taskRetention),
	)

	return q.Scheduler.Register(cronSpec, task, options...)
}

func (q Queue) RemoveScheduler(id string) error {
	return q.Scheduler.Unregister(id)
}
