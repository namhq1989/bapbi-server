package command

import (
	"github.com/namhq1989/bapbi-server/internal/utils/appcontext"
	apperrors "github.com/namhq1989/bapbi-server/internal/utils/error"
	"github.com/namhq1989/bapbi-server/pkg/health/domain"
	"github.com/namhq1989/bapbi-server/pkg/health/dto"
)

type WaterIntakeHandler struct {
	drinkWaterProfileRepository domain.DrinkWaterProfileRepository
	waterIntakeRepository       domain.WaterIntakeLogRepository
	queueRepository             domain.QueueRepository
}

func NewWaterIntakeHandler(
	drinkWaterProfileRepository domain.DrinkWaterProfileRepository,
	waterIntakeRepository domain.WaterIntakeLogRepository,
	queueRepository domain.QueueRepository,
) WaterIntakeHandler {
	return WaterIntakeHandler{
		drinkWaterProfileRepository: drinkWaterProfileRepository,
		waterIntakeRepository:       waterIntakeRepository,
		queueRepository:             queueRepository,
	}
}

func (h WaterIntakeHandler) WaterIntake(ctx *appcontext.AppContext, performerID string, req dto.WaterIntakeRequest) (*dto.WaterIntakeResponse, error) {
	ctx.Logger().Info("new water intake", appcontext.Fields{"performerID": performerID, "amount": req.Amount, "intakeAt": req.IntakeAt})

	// find drink water profile
	ctx.Logger().Text("find drink water profile")
	dwProfile, err := h.drinkWaterProfileRepository.FindDrinkWaterProfileByUserID(ctx, performerID)
	if err != nil {
		ctx.Logger().Error("failed to find drink water profile", err, appcontext.Fields{})
		return nil, err
	}
	if dwProfile == nil || !dwProfile.IsEnabled {
		ctx.Logger().Error("drink water profile not found or disabled", nil, appcontext.Fields{})
		return nil, apperrors.Health.DrinkWaterProfileNotFound
	}

	// convert to domain model
	ctx.Logger().Text("convert to domain model")
	log, err := domain.NewWaterIntakeLog(performerID, req.Amount, req.IntakeAt)
	if err != nil {
		ctx.Logger().Error("failed to convert to domain model", err, appcontext.Fields{})
		return nil, err
	}

	// persist to database
	ctx.Logger().Text("persist to database")
	if err = h.waterIntakeRepository.CreateWaterIntakeLog(ctx, *log); err != nil {
		ctx.Logger().Error("failed to persist to database", err, appcontext.Fields{})
		return nil, err
	}

	// add to queue
	if err = h.queueRepository.EnqueueNewWaterIntakeLog(ctx, *log); err != nil {
		ctx.Logger().Error("failed to add to queue", err, appcontext.Fields{})
		return nil, err
	}

	// respond
	return nil, nil
}
