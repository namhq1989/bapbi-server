package domain

import (
	"time"

	"github.com/namhq1989/bapbi-server/internal/database"
)

type DailyHydrationReport struct {
	ID           string
	UserID       string
	GoalAmount   int
	IntakeAmount int
	IntakeTimes  int
	IsAchieved   bool
	Date         time.Time
}

func NewDailyHydrationReport(userID string, goalAmount, intakeAmount, intakeTimes int, date time.Time) (*DailyHydrationReport, error) {
	report := &DailyHydrationReport{
		ID:           database.NewStringID(),
		UserID:       userID,
		GoalAmount:   goalAmount,
		IntakeAmount: intakeAmount,
		IntakeTimes:  intakeTimes,
		IsAchieved:   intakeAmount >= goalAmount,
		Date:         date,
	}

	return report, nil
}
