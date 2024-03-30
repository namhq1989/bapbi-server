package model

import (
	"time"

	"github.com/namhq1989/bapbi-server/internal/database"
	apperrors "github.com/namhq1989/bapbi-server/internal/utils/error"
	"github.com/namhq1989/bapbi-server/pkg/health/domain"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type DrinkWaterProfile struct {
	ID                        primitive.ObjectID `bson:"_id"`
	UserID                    primitive.ObjectID `bson:"userId"`
	IsEnabled                 bool               `bson:"isEnabled"`
	DailyIntakeAmount         int                `bson:"dailyIntakeAmount"`
	HourlyIntakeAmount        int                `bson:"hourlyIntakeAmount"`
	CurrentStreak             int                `bson:"currentStreak"`
	LongestSuccessStreakValue int                `bson:"longestSuccessStreakValue"`
	LongestSuccessStreakAt    time.Time          `bson:"longestSuccessStreakAt"`
	HighestIntakeAmountValue  int                `bson:"highestIntakeAmountValue"`
	HighestIntakeAmountAt     time.Time          `bson:"highestIntakeAmountAt"`
	EnableAt                  time.Time          `bson:"enableAt"`
}

func (m DrinkWaterProfile) ToDomain() domain.DrinkWaterProfile {
	return domain.DrinkWaterProfile{
		ID:                        m.ID.Hex(),
		UserID:                    m.UserID.Hex(),
		IsEnabled:                 m.IsEnabled,
		DailyIntakeAmount:         m.DailyIntakeAmount,
		HourlyIntakeAmount:        m.HourlyIntakeAmount,
		CurrentStreak:             m.CurrentStreak,
		LongestSuccessStreakValue: m.LongestSuccessStreakValue,
		LongestSuccessStreakAt:    m.LongestSuccessStreakAt,
		HighestIntakeAmountValue:  m.HighestIntakeAmountValue,
		HighestIntakeAmountAt:     m.HighestIntakeAmountAt,
		EnableAt:                  m.EnableAt,
	}
}

func (m DrinkWaterProfile) FromDomain(profile domain.DrinkWaterProfile) (*DrinkWaterProfile, error) {
	id, err := database.ObjectIDFromString(profile.ID)
	if err != nil {
		return nil, apperrors.Common.InvalidID
	}

	userID, err := database.ObjectIDFromString(profile.UserID)
	if err != nil {
		return nil, apperrors.User.InvalidUserID
	}

	return &DrinkWaterProfile{
		ID:                        id,
		UserID:                    userID,
		IsEnabled:                 profile.IsEnabled,
		DailyIntakeAmount:         profile.DailyIntakeAmount,
		HourlyIntakeAmount:        profile.HourlyIntakeAmount,
		CurrentStreak:             profile.CurrentStreak,
		LongestSuccessStreakValue: profile.LongestSuccessStreakValue,
		LongestSuccessStreakAt:    profile.LongestSuccessStreakAt,
		HighestIntakeAmountValue:  profile.HighestIntakeAmountValue,
		HighestIntakeAmountAt:     profile.HighestIntakeAmountAt,
		EnableAt:                  profile.EnableAt,
	}, nil
}
