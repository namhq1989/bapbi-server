package domain

import (
	"github.com/namhq1989/bapbi-server/internal/utils/appcontext"
	"time"

	"github.com/namhq1989/bapbi-server/internal/database"
)

type DailyHydrationReportRepository interface {
	FindDailyHydrationReportByUserID(ctx *appcontext.AppContext, userID string, date time.Time) (*DailyHydrationReport, error)
	CreateDailyHydrationReport(ctx *appcontext.AppContext, report DailyHydrationReport) error
	UpdateDailyHydrationReport(ctx *appcontext.AppContext, report DailyHydrationReport) error
}

type DailyHydrationReport struct {
	ID           string
	UserID       string
	GoalAmount   int
	IntakeAmount int
	IntakeTimes  int
	IsAchieved   bool
	Date         time.Time
}

func NewDailyHydrationReport(userID string, goalAmount, intakeAmount int, date time.Time) (*DailyHydrationReport, error) {
	report := &DailyHydrationReport{
		ID:           database.NewStringID(),
		UserID:       userID,
		GoalAmount:   goalAmount,
		IntakeAmount: intakeAmount,
		IntakeTimes:  1,
		IsAchieved:   intakeAmount >= goalAmount,
		Date:         date,
	}

	return report, nil
}

func (d *DailyHydrationReport) AddIntakeAmount(intakeAmount int) error {
	d.IntakeAmount += intakeAmount
	d.IntakeTimes += 1
	d.IsAchieved = d.IntakeAmount >= d.GoalAmount
	return nil
}
