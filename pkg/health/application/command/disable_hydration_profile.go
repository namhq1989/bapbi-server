package command

import (
	"github.com/namhq1989/bapbi-server/internal/utils/appcontext"
	apperrors "github.com/namhq1989/bapbi-server/internal/utils/error"
	"github.com/namhq1989/bapbi-server/pkg/health/domain"
	"github.com/namhq1989/bapbi-server/pkg/health/dto"
)

type DisableHydrationProfileHandler struct {
	healthProfileRepository    domain.HealthProfileRepository
	hydrationProfileRepository domain.HydrationProfileRepository
}

func NewDisableHydrationProfileHandler(healthProfileRepository domain.HealthProfileRepository, hydrationProfileRepository domain.HydrationProfileRepository) DisableHydrationProfileHandler {
	return DisableHydrationProfileHandler{
		healthProfileRepository:    healthProfileRepository,
		hydrationProfileRepository: hydrationProfileRepository,
	}
}

func (h DisableHydrationProfileHandler) DisableHydrationProfile(ctx *appcontext.AppContext, performerID string, _ dto.DisableHydrationProfileRequest) (*dto.DisableHydrationProfileResponse, error) {
	ctx.Logger().Info("disable hydration profile", appcontext.Fields{
		"performerID": performerID,
	})

	// find hydration profile in db
	ctx.Logger().Text("find hydration profile in db")
	hydrationProfile, err := h.hydrationProfileRepository.FindHydrationProfileByUserID(ctx, performerID)
	if err != nil {
		ctx.Logger().Error("failed to find hydration profile in db", err, appcontext.Fields{})
		return nil, err
	}
	if hydrationProfile == nil {
		ctx.Logger().Error("hydration profile not found", nil, appcontext.Fields{})
		return nil, apperrors.Health.HydrationProfileNotFound
	}

	// disable
	if err = hydrationProfile.Disable(); err != nil {
		ctx.Logger().Error("failed to disable hydration profile", err, appcontext.Fields{})
		return nil, err
	}

	ctx.Logger().Text("update hydration profile to database")
	if err = h.hydrationProfileRepository.UpdateHydrationProfile(ctx, *hydrationProfile); err != nil {
		ctx.Logger().Error("failed to update hydration profile to database", err, appcontext.Fields{})
		return nil, err
	}

	ctx.Logger().Text("done disable hydration profile")
	return nil, nil
}
