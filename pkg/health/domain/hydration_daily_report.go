package domain

import (
	"time"

	"github.com/namhq1989/bapbi-server/internal/utils/appcontext"

	"github.com/namhq1989/bapbi-server/internal/database"
)

type HydrationDailyReportRepository interface {
	FindHydrationDailyReportByUserID(ctx *appcontext.AppContext, userID string, date time.Time) (*HydrationDailyReport, error)
	CreateHydrationDailyReport(ctx *appcontext.AppContext, report HydrationDailyReport) error
	UpdateHydrationDailyReport(ctx *appcontext.AppContext, report HydrationDailyReport) error
}

type HydrationDailyReport struct {
	ID            string
	UserID        string
	GoalAmount    int
	IntakeAmount  int
	IntakeTimes   int
	IsAchieved    bool
	CurrentStreak int
	Date          time.Time
}

func NewHydrationDailyReport(userID string, goalAmount, intakeAmount int, date time.Time) (*HydrationDailyReport, error) {
	report := &HydrationDailyReport{
		ID:            database.NewStringID(),
		UserID:        userID,
		GoalAmount:    goalAmount,
		IntakeAmount:  intakeAmount,
		IntakeTimes:   1,
		IsAchieved:    intakeAmount >= goalAmount,
		CurrentStreak: 0,
		Date:          date,
	}

	return report, nil
}

func (d *HydrationDailyReport) AddIntakeAmount(intakeAmount int) error {
	d.IntakeAmount += intakeAmount
	d.IntakeTimes += 1
	d.IsAchieved = d.IntakeAmount >= d.GoalAmount
	if d.IsAchieved {
		d.CurrentStreak += 1
	}
	return nil
}
