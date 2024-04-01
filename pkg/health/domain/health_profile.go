package domain

import (
	"math"
	"time"

	"github.com/namhq1989/bapbi-server/internal/database"

	"github.com/namhq1989/bapbi-server/internal/utils/appcontext"
	apperrors "github.com/namhq1989/bapbi-server/internal/utils/error"
)

type HealthProfileRepository interface {
	CreateHealthProfile(ctx *appcontext.AppContext, profile HealthProfile) error
	FindHealthProfileByUserID(ctx *appcontext.AppContext, userID string) (*HealthProfile, error)
}

type HealthProfile struct {
	ID          string
	UserID      string
	WeightInKg  int
	HeightInCm  int
	BMI         float64
	WakeUpHour  int
	BedtimeHour int
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func NewHealthProfile(userID string, weightInKg int, heightInCm int, wakeUpHour int, bedtimeHour int) (*HealthProfile, error) {
	if userID == "" {
		return nil, apperrors.Common.InvalidID
	}

	if heightInCm == 0 {
		return nil, apperrors.Health.InvalidHeight
	}

	if weightInKg == 0 {
		return nil, apperrors.Health.InvalidWeight
	}

	if wakeUpHour == 0 && bedtimeHour == 0 {
		return nil, apperrors.Health.InvalidWakingHours
	}

	if wakeUpHour >= bedtimeHour {
		return nil, apperrors.Health.InvalidWakingHours
	}

	profile := &HealthProfile{
		ID:          database.NewStringID(),
		UserID:      userID,
		WeightInKg:  weightInKg,
		HeightInCm:  heightInCm,
		WakeUpHour:  wakeUpHour,
		BedtimeHour: bedtimeHour,
		CreatedAt:   time.Now(),
	}
	if err := profile.SetBMI(weightInKg, heightInCm); err != nil {
		return nil, err
	}

	return profile, nil
}

func (d *HealthProfile) SetWeightInKg(value int) error {
	if value < 0 {
		return apperrors.Health.InvalidWeight
	}
	d.WeightInKg = value
	return nil
}

func (d *HealthProfile) SetHeightInCm(value int) error {
	if value < 0 {
		return apperrors.Health.InvalidHeight
	}
	d.HeightInCm = value
	return nil
}

func (d *HealthProfile) SetBMI(weight, height int) error {
	if weight == 0 || height == 0 {
		return apperrors.Health.InvalidBMI
	}

	// convert height from centimeters to meters
	heightInM := float64(height) / 100

	// calculate
	shift := math.Pow(10, 2)
	bmi := float64(weight) / (heightInM * heightInM)
	d.BMI = math.Round(bmi*shift) / shift

	return nil
}

func (d *HealthProfile) SetWakingTimes(wakeUpHour int, bedtimeHour int) error {
	if wakeUpHour == 0 && bedtimeHour == 0 {
		return apperrors.Health.InvalidWakingHours
	}
	if wakeUpHour >= bedtimeHour {
		return apperrors.Health.InvalidWakingHours
	}
	d.WakeUpHour = wakeUpHour
	d.BedtimeHour = bedtimeHour
	return nil
}

func (d *HealthProfile) SetUpdatedAt() {
	d.UpdatedAt = time.Now()
}

const waterMlPerKg = 40

func (d *HealthProfile) GetDailyAndHourlyIntakeAmount() (int, int) {
	wakingHours := d.BedtimeHour - d.WakeUpHour
	dailyAmount := d.WeightInKg * waterMlPerKg
	hourlyAmount := dailyAmount / wakingHours

	return dailyAmount, hourlyAmount
}
