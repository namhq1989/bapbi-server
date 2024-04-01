package domain

import (
	"github.com/namhq1989/bapbi-server/internal/utils/appcontext"
)

var QueueTypeNames = struct {
	NewWaterIntakeLog string
}{
	NewWaterIntakeLog: "health:drinkWater.new",
}

type QueueRepository interface {
	EnqueueNewWaterIntakeLog(ctx *appcontext.AppContext, log WaterIntakeLog) error
}
