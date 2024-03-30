package model

import (
	"time"

	"github.com/namhq1989/bapbi-server/internal/database"
	apperrors "github.com/namhq1989/bapbi-server/internal/utils/error"
	"github.com/namhq1989/bapbi-server/pkg/health/domain"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type WaterIntakeLog struct {
	ID        primitive.ObjectID `bson:"_id"`
	UserID    primitive.ObjectID `bson:"userId"`
	Amount    int                `bson:"amount"`
	IntakeAt  time.Time          `bson:"intakeAt"`
	CreatedAt time.Time          `bson:"createdAt"`
}

func (m WaterIntakeLog) ToDomain() domain.WaterIntakeLog {
	return domain.WaterIntakeLog{
		ID:        m.ID.Hex(),
		UserID:    m.UserID.Hex(),
		Amount:    m.Amount,
		IntakeAt:  m.IntakeAt,
		CreatedAt: m.CreatedAt,
	}
}

func (m WaterIntakeLog) FromDomain(log domain.WaterIntakeLog) (*WaterIntakeLog, error) {
	id, err := database.ObjectIDFromString(log.ID)
	if err != nil {
		return nil, apperrors.Common.InvalidID
	}

	userID, err := database.ObjectIDFromString(log.UserID)
	if err != nil {
		return nil, apperrors.User.InvalidUserID
	}

	return &WaterIntakeLog{
		ID:        id,
		UserID:    userID,
		Amount:    log.Amount,
		IntakeAt:  log.IntakeAt,
		CreatedAt: log.CreatedAt,
	}, nil
}
