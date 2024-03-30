package domain

import (
	"github.com/namhq1989/bapbi-server/internal/database"
	"time"

	"github.com/namhq1989/bapbi-server/internal/utils/appcontext"
	apperrors "github.com/namhq1989/bapbi-server/internal/utils/error"
)

type HealthProfileRepository interface {
	CreateHealthProfile(ctx *appcontext.AppContext, profile HealthProfile) (*HealthProfile, error)
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

func NewHealthProfile(userID string, weightInKg int, heightInCm int, wakeUpHour int, bedTimeHour int) (*HealthProfile, error) {
	if userID == "" {
		return nil, apperrors.Common.InvalidID
	}

	if heightInCm == 0 {
		return nil, apperrors.Common.InvalidHeight
	}

	if weightInKg == 0 {
		return nil, apperrors.Common.InvalidWeight
	}

	if wakeUpHour == 0 || bedTimeHour == 0 {
		return nil, apperrors.Common.InvalidWakingHours
	}

	return &HealthProfile{
		ID:          database.NewStringID(),
		UserID:      userID,
		WeightInKg:  weightInKg,
		HeightInCm:  heightInCm,
		WakeUpHour:  wakeUpHour,
		BedtimeHour: bedTimeHour,
		CreatedAt:   time.Now(),
	}, nil
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

func (d *HealthProfile) SetBMI(value float64) error {
	if value < 0 {
		return apperrors.Health.InvalidBMI
	}
	d.BMI = value
	return nil
}

func (d *HealthProfile) SetWakingTimes(wakeUpHour int, bedtimeHour int) error {
	if wakeUpHour == 0 || bedtimeHour == 0 {
		return apperrors.Common.InvalidWakingHours
	}
	d.WakeUpHour = wakeUpHour
	d.BedtimeHour = bedtimeHour
	return nil
}

func (d *HealthProfile) SetUpdatedAt() {
	d.UpdatedAt = time.Now()
}
