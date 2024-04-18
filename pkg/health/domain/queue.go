package domain

import (
	"github.com/namhq1989/bapbi-server/internal/utils/appcontext"
)

type QueueRepository interface {
	EnqueueNewWaterIntakeLog(ctx *appcontext.AppContext, log WaterIntakeLog) error
}
