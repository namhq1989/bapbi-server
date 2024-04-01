package model

import (
	"time"

	"github.com/namhq1989/bapbi-server/internal/database"
	apperrors "github.com/namhq1989/bapbi-server/internal/utils/error"
	"github.com/namhq1989/bapbi-server/pkg/health/domain"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type DailyHydrationReport struct {
	ID           primitive.ObjectID `bson:"_id"`
	UserID       primitive.ObjectID `bson:"userId"`
	GoalAmount   int                `bson:"goalAmount"`
	IntakeAmount int                `bson:"intakeAmount"`
	IntakeTimes  int                `bson:"intakeTimes"`
	IsAchieved   bool               `bson:"isAchieved"`
	Date         time.Time          `bson:"date"`
}

func (d DailyHydrationReport) ToDomain() domain.DailyHydrationReport {
	return domain.DailyHydrationReport{
		ID:           d.ID.Hex(),
		UserID:       d.UserID.Hex(),
		GoalAmount:   d.GoalAmount,
		IntakeAmount: d.IntakeAmount,
		IntakeTimes:  d.IntakeTimes,
		IsAchieved:   d.IsAchieved,
		Date:         d.Date,
	}
}

func (d DailyHydrationReport) FromDomain(report domain.DailyHydrationReport) (*DailyHydrationReport, error) {
	id, err := database.ObjectIDFromString(report.ID)
	if err != nil {
		return nil, apperrors.Common.InvalidID
	}

	userID, err := database.ObjectIDFromString(report.UserID)
	if err != nil {
		return nil, apperrors.User.InvalidUserID
	}

	return &DailyHydrationReport{
		ID:           id,
		UserID:       userID,
		GoalAmount:   report.GoalAmount,
		IntakeAmount: report.IntakeAmount,
		IntakeTimes:  report.IntakeTimes,
		IsAchieved:   report.IsAchieved,
		Date:         report.Date,
	}, nil
}
