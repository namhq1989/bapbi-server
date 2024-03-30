package domain

import (
	"github.com/namhq1989/bapbi-server/internal/utils/appcontext"
	"time"

	"github.com/namhq1989/bapbi-server/internal/database"

	apperrors "github.com/namhq1989/bapbi-server/internal/utils/error"
)

type DrinkWaterProfileRepository interface {
	CreateDrinkWaterProfile(ctx *appcontext.AppContext, profile DrinkWaterProfile) error
	UpdateDrinkWaterProfile(ctx *appcontext.AppContext, profile DrinkWaterProfile) error
}

type DrinkWaterProfile struct {
	ID                        string
	UserID                    string
	IsEnabled                 bool
	DailyIntakeAmount         int
	HourlyIntakeAmount        int
	CurrentStreak             int
	LongestSuccessStreakValue int
	LongestSuccessStreakAt    time.Time
	HighestIntakeAmountValue  int
	HighestIntakeAmountAt     time.Time
	EnableAt                  time.Time
}

func NewDrinkWaterProfile(userID string, dailyIntakeAmount int, hourlyIntakeAmount int) (*DrinkWaterProfile, error) {
	if dailyIntakeAmount == 0 {
		return nil, apperrors.Health.InvalidDailyIntakeAmount
	}

	if hourlyIntakeAmount == 0 {
		return nil, apperrors.Health.InvalidDailyIntakeAmount
	}

	return &DrinkWaterProfile{
		ID:                        database.NewStringID(),
		UserID:                    userID,
		IsEnabled:                 true,
		DailyIntakeAmount:         dailyIntakeAmount,
		HourlyIntakeAmount:        hourlyIntakeAmount,
		CurrentStreak:             0,
		LongestSuccessStreakValue: 0,
		LongestSuccessStreakAt:    time.Time{},
		HighestIntakeAmountValue:  0,
		HighestIntakeAmountAt:     time.Time{},
		EnableAt:                  time.Now(),
	}, nil
}

func (d *DrinkWaterProfile) Enable() error {
	d.IsEnabled = true
	d.EnableAt = time.Now()
	return nil
}

func (d *DrinkWaterProfile) Disable() error {
	d.IsEnabled = false
	return nil
}

func (d *DrinkWaterProfile) SetDailyIntakeAmount(value int) error {
	if value < 0 {
		return apperrors.Health.InvalidIntakeAmount
	}
	d.DailyIntakeAmount = value
	return nil
}

func (d *DrinkWaterProfile) SetHourlyIntakeAmount(value int) error {
	if value < 0 {
		return apperrors.Health.InvalidIntakeAmount
	}
	d.HourlyIntakeAmount = value
	return nil
}

func (d *DrinkWaterProfile) ResetStreak() error {
	d.CurrentStreak = 0
	return nil
}

func (d *DrinkWaterProfile) SetStreak(value int) error {
	if value < 0 {
		return apperrors.Health.InvalidStreak

	}

	d.CurrentStreak = value
	if d.CurrentStreak > d.LongestSuccessStreakValue {
		d.LongestSuccessStreakValue = d.CurrentStreak
		d.LongestSuccessStreakAt = time.Now()
	}
	return nil
}

func (d *DrinkWaterProfile) SetHighestIntakeAmount(value int) error {
	if value < 0 {
		return apperrors.Health.InvalidIntakeAmount
	}

	if value > d.HighestIntakeAmountValue {
		d.HighestIntakeAmountValue = value
		d.HighestIntakeAmountAt = time.Now()
	}
	return nil
}
