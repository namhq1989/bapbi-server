package command

import (
	"github.com/namhq1989/bapbi-server/internal/utils/appcontext"
	apperrors "github.com/namhq1989/bapbi-server/internal/utils/error"
	"github.com/namhq1989/bapbi-server/pkg/health/domain"
	"github.com/namhq1989/bapbi-server/pkg/health/dto"
)

type EnableHydrationProfileHandler struct {
	healthProfileRepository    domain.HealthProfileRepository
	hydrationProfileRepository domain.HydrationProfileRepository
}

func NewEnableHydrationProfileHandler(healthProfileRepository domain.HealthProfileRepository, hydrationProfileRepository domain.HydrationProfileRepository) EnableHydrationProfileHandler {
	return EnableHydrationProfileHandler{
		healthProfileRepository:    healthProfileRepository,
		hydrationProfileRepository: hydrationProfileRepository,
	}
}

func (h EnableHydrationProfileHandler) EnableHydrationProfile(ctx *appcontext.AppContext, performerID string, _ dto.EnableHydrationProfileRequest) (*dto.EnableHydrationProfileResponse, error) {
	ctx.Logger().Info("enable hydration profile", appcontext.Fields{
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

	// find hydration profile in db
	ctx.Logger().Text("find hydration profile in db")
	hydrationProfile, err := h.hydrationProfileRepository.FindHydrationProfileByUserID(ctx, performerID)
	if err != nil {
		ctx.Logger().Error("failed to find hydration profile in db", err, appcontext.Fields{})
		return nil, err
	}

	// get daily and hourly intake amount from health profile
	dailyIntakeAmount, hourlyIntakeAmount := healthProfile.GetDailyAndHourlyIntakeAmount()

	if hydrationProfile == nil {
		// if not found, create new hydration profile
		ctx.Logger().Text("hydration profile not found, create new")
		hydrationProfile, err = domain.NewHydrationProfile(performerID, dailyIntakeAmount, hourlyIntakeAmount)
		if err != nil {
			ctx.Logger().Error("failed to create new hydration profile", err, appcontext.Fields{})
			return nil, err
		}

		// save
		ctx.Logger().Text("persist hydration profile to database")
		if err = h.hydrationProfileRepository.CreateHydrationProfile(ctx, *hydrationProfile); err != nil {
			ctx.Logger().Error("failed to persist hydration profile to database", err, appcontext.Fields{})
			return nil, err
		}
	} else {
		// if found, update hydration profile
		ctx.Logger().Text("hydration profile found, update")
		_ = hydrationProfile.SetDailyIntakeAmount(dailyIntakeAmount)
		_ = hydrationProfile.SetHourlyIntakeAmount(hourlyIntakeAmount)

		// enable
		if err = hydrationProfile.Enable(); err != nil {
			ctx.Logger().Error("failed to enable hydration profile", err, appcontext.Fields{})
			return nil, err
		}

		ctx.Logger().Text("update hydration profile to database")
		if err = h.hydrationProfileRepository.UpdateHydrationProfile(ctx, *hydrationProfile); err != nil {
			ctx.Logger().Error("failed to update hydration profile to database", err, appcontext.Fields{})
			return nil, err
		}
	}

	ctx.Logger().Text("done enable hydration profile")
	return nil, nil
}
