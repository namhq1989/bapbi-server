package domain

import (
	"time"

	"github.com/namhq1989/bapbi-server/internal/utils/manipulation"

	"github.com/namhq1989/bapbi-server/internal/utils/appcontext"

	"github.com/namhq1989/bapbi-server/internal/database"

	apperrors "github.com/namhq1989/bapbi-server/internal/utils/error"
)

type HydrationProfileRepository interface {
	FindHydrationProfileByUserID(ctx *appcontext.AppContext, userID string) (*HydrationProfile, error)
	CreateHydrationProfile(ctx *appcontext.AppContext, profile HydrationProfile) error
	UpdateHydrationProfile(ctx *appcontext.AppContext, profile HydrationProfile) error
}

type HydrationProfile struct {
	ID                        string
	UserID                    string
	IsEnabled                 bool
	DailyIntakeAmount         int
	HourlyIntakeAmount        int
	CurrentStreakValue        int
	CurrentStreakDate         time.Time
	LongestSuccessStreakValue int
	LongestSuccessStreakAt    time.Time
	HighestIntakeAmountValue  int
	HighestIntakeAmountAt     time.Time
	EnabledAt                 time.Time
	DisabledAt                time.Time
}

func NewHydrationProfile(userID string, dailyIntakeAmount int, hourlyIntakeAmount int) (*HydrationProfile, error) {
	if dailyIntakeAmount == 0 {
		return nil, apperrors.Health.InvalidDailyIntakeAmount
	}

	if hourlyIntakeAmount == 0 {
		return nil, apperrors.Health.InvalidHourlyIntakeAmount
	}

	return &HydrationProfile{
		ID:                        database.NewStringID(),
		UserID:                    userID,
		IsEnabled:                 true,
		DailyIntakeAmount:         dailyIntakeAmount,
		HourlyIntakeAmount:        hourlyIntakeAmount,
		CurrentStreakValue:        0,
		CurrentStreakDate:         time.Time{},
		LongestSuccessStreakValue: 0,
		LongestSuccessStreakAt:    time.Time{},
		HighestIntakeAmountValue:  0,
		HighestIntakeAmountAt:     time.Time{},
		EnabledAt:                 time.Now(),
		DisabledAt:                time.Time{},
	}, nil
}

func (d *HydrationProfile) Enable() error {
	d.IsEnabled = true
	d.EnabledAt = time.Now()
	return nil
}

func (d *HydrationProfile) Disable() error {
	d.IsEnabled = false
	d.DisabledAt = time.Now()
	return nil
}

func (d *HydrationProfile) SetDailyIntakeAmount(value int) error {
	if value < 0 {
		return apperrors.Health.InvalidIntakeAmount
	}
	d.DailyIntakeAmount = value
	return nil
}

func (d *HydrationProfile) SetHourlyIntakeAmount(value int) error {
	if value < 0 {
		return apperrors.Health.InvalidIntakeAmount
	}
	d.HourlyIntakeAmount = value
	return nil
}

func (d *HydrationProfile) ResetStreak() {
	d.CurrentStreakValue = 1
	d.CurrentStreakDate = manipulation.StartOfToday()
	if d.CurrentStreakValue > d.LongestSuccessStreakValue {
		d.LongestSuccessStreakValue = d.CurrentStreakValue
		d.LongestSuccessStreakAt = manipulation.StartOfToday()
	}
}

func (d *HydrationProfile) IncreaseStreak() {
	d.CurrentStreakValue += 1
	d.CurrentStreakDate = manipulation.StartOfToday()
	if d.CurrentStreakValue > d.LongestSuccessStreakValue {
		d.LongestSuccessStreakValue = d.CurrentStreakValue
		d.LongestSuccessStreakAt = manipulation.StartOfToday()
	}
}

func (d *HydrationProfile) SetHighestIntakeAmount(value int) error {
	if value < 0 {
		return apperrors.Health.InvalidIntakeAmount
	}

	if value > d.HighestIntakeAmountValue {
		d.HighestIntakeAmountValue = value
		d.HighestIntakeAmountAt = time.Now()
	}
	return nil
}
