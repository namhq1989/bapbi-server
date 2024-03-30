package model

import (
	"time"

	"github.com/namhq1989/bapbi-server/internal/database"
	apperrors "github.com/namhq1989/bapbi-server/internal/utils/error"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/namhq1989/bapbi-server/pkg/health/domain"
)

type HealthProfile struct {
	ID          primitive.ObjectID `bson:"_id"`
	UserID      primitive.ObjectID `bson:"userId"`
	WeightInKg  int                `bson:"weightInKg"`
	HeightInCm  int                `bson:"heightInCm"`
	BMI         float64            `bson:"bmi"`
	WakeUpHour  int                `bson:"wakeUpHour"`
	BedtimeHour int                `bson:"bedtimeHour"`
	CreatedAt   time.Time          `bson:"createdAt"`
	UpdatedAt   time.Time          `bson:"updatedAt"`
}

func (m HealthProfile) ToDomain() domain.HealthProfile {
	return domain.HealthProfile{
		ID:          m.ID.Hex(),
		UserID:      m.UserID.Hex(),
		WeightInKg:  m.WeightInKg,
		HeightInCm:  m.HeightInCm,
		BMI:         m.BMI,
		WakeUpHour:  m.WakeUpHour,
		BedtimeHour: m.BedtimeHour,
		CreatedAt:   m.CreatedAt,
		UpdatedAt:   m.UpdatedAt,
	}
}

func (m HealthProfile) FromDomain(profile domain.HealthProfile) (*HealthProfile, error) {
	id, err := database.ObjectIDFromString(profile.ID)
	if err != nil {
		return nil, apperrors.Common.InvalidID
	}

	userID, err := database.ObjectIDFromString(profile.UserID)
	if err != nil {
		return nil, apperrors.User.InvalidUserID
	}

	return &HealthProfile{
		ID:          id,
		UserID:      userID,
		WeightInKg:  profile.WeightInKg,
		HeightInCm:  profile.HeightInCm,
		BMI:         profile.BMI,
		WakeUpHour:  profile.WakeUpHour,
		BedtimeHour: profile.BedtimeHour,
		CreatedAt:   profile.CreatedAt,
		UpdatedAt:   profile.UpdatedAt,
	}, nil
}
