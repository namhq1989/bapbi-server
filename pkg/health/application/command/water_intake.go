package command

import (
	"github.com/namhq1989/bapbi-server/internal/utils/appcontext"
	apperrors "github.com/namhq1989/bapbi-server/internal/utils/error"
	"github.com/namhq1989/bapbi-server/pkg/health/domain"
	"github.com/namhq1989/bapbi-server/pkg/health/dto"
)

type WaterIntakeHandler struct {
	hydrationProfileRepository domain.HydrationProfileRepository
	waterIntakeRepository      domain.WaterIntakeLogRepository
	queueRepository            domain.QueueRepository
}

func NewWaterIntakeHandler(
	hydrationProfileRepository domain.HydrationProfileRepository,
	waterIntakeRepository domain.WaterIntakeLogRepository,
	queueRepository domain.QueueRepository,
) WaterIntakeHandler {
	return WaterIntakeHandler{
		hydrationProfileRepository: hydrationProfileRepository,
		waterIntakeRepository:      waterIntakeRepository,
		queueRepository:            queueRepository,
	}
}

func (h WaterIntakeHandler) WaterIntake(ctx *appcontext.AppContext, performerID string, req dto.WaterIntakeRequest) (*dto.WaterIntakeResponse, error) {
	ctx.Logger().Info("new water intake", appcontext.Fields{"performerID": performerID, "amount": req.Amount, "intakeAt": req.IntakeAt})

	// find hydration profile
	ctx.Logger().Text("find hydration profile")
	hydrationProfile, err := h.hydrationProfileRepository.FindHydrationProfileByUserID(ctx, performerID)
	if err != nil {
		ctx.Logger().Error("failed to find hydration profile", err, appcontext.Fields{})
		return nil, err
	}
	if hydrationProfile == nil || !hydrationProfile.IsEnabled {
		ctx.Logger().Error("hydration profile not found or disabled", nil, appcontext.Fields{})
		return nil, apperrors.Health.HydrationProfileNotFound
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
