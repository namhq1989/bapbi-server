package command

import (
	"github.com/namhq1989/bapbi-server/internal/utils/appcontext"
	apperrors "github.com/namhq1989/bapbi-server/internal/utils/error"
	"github.com/namhq1989/bapbi-server/pkg/health/domain"
	"github.com/namhq1989/bapbi-server/pkg/health/dto"
)

type DisableDrinkWaterProfileHandler struct {
	healthProfileRepository     domain.HealthProfileRepository
	drinkWaterProfileRepository domain.DrinkWaterProfileRepository
}

func NewDisableDrinkWaterProfileHandler(healthProfileRepository domain.HealthProfileRepository, drinkWaterProfileRepository domain.DrinkWaterProfileRepository) DisableDrinkWaterProfileHandler {
	return DisableDrinkWaterProfileHandler{
		healthProfileRepository:     healthProfileRepository,
		drinkWaterProfileRepository: drinkWaterProfileRepository,
	}
}

func (h DisableDrinkWaterProfileHandler) DisableDrinkWaterProfile(ctx *appcontext.AppContext, performerID string, _ dto.DisableDrinkWaterProfileRequest) (*dto.DisableDrinkWaterProfileResponse, error) {
	ctx.Logger().Info("disable drink water profile", appcontext.Fields{
		"performerID": performerID,
	})

	// find drink water profile in db
	ctx.Logger().Text("find drink water profile in db")
	drinkWaterProfile, err := h.drinkWaterProfileRepository.FindDrinkWaterProfileByUserID(ctx, performerID)
	if err != nil {
		ctx.Logger().Error("failed to find drink water profile in db", err, appcontext.Fields{})
		return nil, err
	}
	if drinkWaterProfile == nil {
		ctx.Logger().Error("drink water profile not found", nil, appcontext.Fields{})
		return nil, apperrors.Health.DrinkWaterProfileNotFound
	}

	// disable
	if err = drinkWaterProfile.Disable(); err != nil {
		ctx.Logger().Error("failed to disable drink water profile", err, appcontext.Fields{})
		return nil, err
	}

	ctx.Logger().Text("update drink water profile to database")
	if err = h.drinkWaterProfileRepository.UpdateDrinkWaterProfile(ctx, *drinkWaterProfile); err != nil {
		ctx.Logger().Error("failed to update drink water profile to database", err, appcontext.Fields{})
		return nil, err
	}

	ctx.Logger().Text("done disable drink water profile")
	return nil, nil
}
