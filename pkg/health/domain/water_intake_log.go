package domain

import (
	"github.com/namhq1989/bapbi-server/internal/utils/appcontext"
	"time"

	"github.com/namhq1989/bapbi-server/internal/database"
	apperrors "github.com/namhq1989/bapbi-server/internal/utils/error"
)

type WaterIntakeLogRepository interface {
	CreateWaterIntakeLog(ctx *appcontext.AppContext, log WaterIntakeLog) error
}

type WaterIntakeLog struct {
	ID        string
	UserID    string
	Amount    int
	IntakeAt  time.Time
	CreatedAt time.Time
}

func NewWaterIntakeLog(userID string, amount int, intakeAt time.Time) (*WaterIntakeLog, error) {
	if amount < 0 {
		return nil, apperrors.Health.InvalidIntakeAmount
	}

	createdAt := time.Now()
	if intakeAt.IsZero() {
		intakeAt = createdAt
	}

	return &WaterIntakeLog{
		ID:        database.NewStringID(),
		UserID:    userID,
		Amount:    amount,
		IntakeAt:  intakeAt,
		CreatedAt: createdAt,
	}, nil
}
