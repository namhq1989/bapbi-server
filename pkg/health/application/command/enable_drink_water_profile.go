package command

import (
	"github.com/namhq1989/bapbi-server/internal/utils/appcontext"
	apperrors "github.com/namhq1989/bapbi-server/internal/utils/error"
	"github.com/namhq1989/bapbi-server/pkg/health/domain"
	"github.com/namhq1989/bapbi-server/pkg/health/dto"
)

type EnableDrinkWaterProfileHandler struct {
	healthProfileRepository     domain.HealthProfileRepository
	drinkWaterProfileRepository domain.DrinkWaterProfileRepository
}

func NewEnableDrinkWaterProfileHandler(healthProfileRepository domain.HealthProfileRepository, drinkWaterProfileRepository domain.DrinkWaterProfileRepository) EnableDrinkWaterProfileHandler {
	return EnableDrinkWaterProfileHandler{
		healthProfileRepository:     healthProfileRepository,
		drinkWaterProfileRepository: drinkWaterProfileRepository,
	}
}

func (h EnableDrinkWaterProfileHandler) EnableDrinkWaterProfile(ctx *appcontext.AppContext, performerID string, _ dto.EnableDrinkWaterProfileRequest) (*dto.EnableDrinkWaterProfileResponse, error) {
	ctx.Logger().Info("enable drink water profile", appcontext.Fields{
		"performerID": performerID,
	})

	// find health profile in db
	ctx.Logger().Text("find health profile in db")
	healthProfile, err := h.healthProfileRepository.FindHealthProfileByUserID(ctx, performerID)
	if err != nil {
		ctx.Logger().Error("failed to find health profile in db", err, appcontext.Fields{})
		return nil, err
	}
	if healthProfile == nil {
		ctx.Logger().Error("health profile not found", nil, appcontext.Fields{})
		return nil, apperrors.Health.HealthProfileNotFound
	}

	// find drink water profile in db
	ctx.Logger().Text("find drink water profile in db")
	drinkWaterProfile, err := h.drinkWaterProfileRepository.FindDrinkWaterProfileByUserID(ctx, performerID)
	if err != nil {
		ctx.Logger().Error("failed to find drink water profile in db", err, appcontext.Fields{})
		return nil, err
	}

	// get daily and hourly intake amount from health profile
	dailyIntakeAmount, hourlyIntakeAmount := healthProfile.GetDailyAndHourlyIntakeAmount()

	if drinkWaterProfile == nil {
		// if not found, create new drink water profile
		ctx.Logger().Text("drink water profile not found, create new")
		drinkWaterProfile, err = domain.NewDrinkWaterProfile(performerID, dailyIntakeAmount, hourlyIntakeAmount)
		if err != nil {
			ctx.Logger().Error("failed to create new drink water profile", err, appcontext.Fields{})
			return nil, err
		}

		// save
		ctx.Logger().Text("persist drink water profile to database")
		if err = h.drinkWaterProfileRepository.CreateDrinkWaterProfile(ctx, *drinkWaterProfile); err != nil {
			ctx.Logger().Error("failed to persist drink water profile to database", err, appcontext.Fields{})
			return nil, err
		}
	} else {
		// if found, update drink water profile
		ctx.Logger().Text("drink water profile found, update")
		_ = drinkWaterProfile.SetDailyIntakeAmount(dailyIntakeAmount)
		_ = drinkWaterProfile.SetHourlyIntakeAmount(hourlyIntakeAmount)

		// enable
		if err = drinkWaterProfile.Enable(); err != nil {
			ctx.Logger().Error("failed to enable drink water profile", err, appcontext.Fields{})
			return nil, err
		}

		ctx.Logger().Text("update drink water profile to database")
		if err = h.drinkWaterProfileRepository.UpdateDrinkWaterProfile(ctx, *drinkWaterProfile); err != nil {
			ctx.Logger().Error("failed to update drink water profile to database", err, appcontext.Fields{})
			return nil, err
		}
	}

	ctx.Logger().Text("done enable drink water profile")
	return nil, nil
}
