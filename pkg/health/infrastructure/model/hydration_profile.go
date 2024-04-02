package model

import (
	"time"

	"github.com/namhq1989/bapbi-server/internal/database"
	apperrors "github.com/namhq1989/bapbi-server/internal/utils/error"
	"github.com/namhq1989/bapbi-server/pkg/health/domain"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type HydrationProfile struct {
	ID                        primitive.ObjectID `bson:"_id"`
	UserID                    primitive.ObjectID `bson:"userId"`
	IsEnabled                 bool               `bson:"isEnabled"`
	DailyIntakeAmount         int                `bson:"dailyIntakeAmount"`
	HourlyIntakeAmount        int                `bson:"hourlyIntakeAmount"`
	CurrentStreakValue        int                `bson:"currentStreakValue"`
	CurrentStreakDate         time.Time          `bson:"currentStreakDate"`
	LongestSuccessStreakValue int                `bson:"longestSuccessStreakValue"`
	LongestSuccessStreakAt    time.Time          `bson:"longestSuccessStreakAt"`
	HighestIntakeAmountValue  int                `bson:"highestIntakeAmountValue"`
	HighestIntakeAmountAt     time.Time          `bson:"highestIntakeAmountAt"`
	EnabledAt                 time.Time          `bson:"enabledAt"`
	DisabledAt                time.Time          `bson:"disabledAt"`
}

func (m HydrationProfile) ToDomain() domain.HydrationProfile {
	return domain.HydrationProfile{
		ID:                        m.ID.Hex(),
		UserID:                    m.UserID.Hex(),
		IsEnabled:                 m.IsEnabled,
		DailyIntakeAmount:         m.DailyIntakeAmount,
		HourlyIntakeAmount:        m.HourlyIntakeAmount,
		CurrentStreakValue:        m.CurrentStreakValue,
		CurrentStreakDate:         m.CurrentStreakDate,
		LongestSuccessStreakValue: m.LongestSuccessStreakValue,
		LongestSuccessStreakAt:    m.LongestSuccessStreakAt,
		HighestIntakeAmountValue:  m.HighestIntakeAmountValue,
		HighestIntakeAmountAt:     m.HighestIntakeAmountAt,
		EnabledAt:                 m.EnabledAt,
		DisabledAt:                m.DisabledAt,
	}
}

func (m HydrationProfile) FromDomain(profile domain.HydrationProfile) (*HydrationProfile, error) {
	id, err := database.ObjectIDFromString(profile.ID)
	if err != nil {
		return nil, apperrors.Common.InvalidID
	}

	userID, err := database.ObjectIDFromString(profile.UserID)
	if err != nil {
		return nil, apperrors.User.InvalidUserID
	}

	return &HydrationProfile{
		ID:                        id,
		UserID:                    userID,
		IsEnabled:                 profile.IsEnabled,
		DailyIntakeAmount:         profile.DailyIntakeAmount,
		HourlyIntakeAmount:        profile.HourlyIntakeAmount,
		CurrentStreakValue:        profile.CurrentStreakValue,
		CurrentStreakDate:         profile.CurrentStreakDate,
		LongestSuccessStreakValue: profile.LongestSuccessStreakValue,
		LongestSuccessStreakAt:    profile.LongestSuccessStreakAt,
		HighestIntakeAmountValue:  profile.HighestIntakeAmountValue,
		HighestIntakeAmountAt:     profile.HighestIntakeAmountAt,
		EnabledAt:                 profile.EnabledAt,
		DisabledAt:                profile.DisabledAt,
	}, nil
}
