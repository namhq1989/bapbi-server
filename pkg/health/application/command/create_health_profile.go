package command

import (
	"github.com/namhq1989/bapbi-server/internal/utils/appcontext"
	apperrors "github.com/namhq1989/bapbi-server/internal/utils/error"
	"github.com/namhq1989/bapbi-server/pkg/health/domain"
	"github.com/namhq1989/bapbi-server/pkg/health/dto"
)

type CreateHealthProfileHandler struct {
	healthProfileRepository domain.HealthProfileRepository
}

func NewCreateHealthProfileHandler(healthProfileRepository domain.HealthProfileRepository) CreateHealthProfileHandler {
	return CreateHealthProfileHandler{
		healthProfileRepository: healthProfileRepository,
	}
}

func (h CreateHealthProfileHandler) CreateHealthProfile(ctx *appcontext.AppContext, performerID string, req dto.CreateHealthProfileRequest) (*dto.CreateHealthProfileResponse, error) {
	ctx.Logger().Info("create new health profile", appcontext.Fields{
		"performerID": performerID,
		"height":      req.Height,
		"weight":      req.Weight,
		"wakeupHour":  req.WakeUpHour,
		"bedtimeHour": req.BedtimeHour,
	})

	// new profile from domain
	ctx.Logger().Text("new profile from domain")
	domainProfile, err := domain.NewHealthProfile(performerID, req.Weight, req.Height, req.WakeUpHour, req.BedtimeHour)
	if err != nil {
		ctx.Logger().Error("failed to create new domain profile", err, appcontext.Fields{})
		return nil, err
	}

	// find profile in db, return error if existed
	ctx.Logger().Text("find profile in db")
	dbProfile, err := h.healthProfileRepository.FindHealthProfileByUserID(ctx, performerID)
	if err != nil {
		ctx.Logger().Error("failed to find profile in db", err, appcontext.Fields{})
		return nil, err
	} else if dbProfile != nil {
		ctx.Logger().Error("profile already existed", nil, appcontext.Fields{})
		return nil, apperrors.Common.AlreadyExisted
	}

	// insert to db
	ctx.Logger().Text("insert profile to db")
	if err = h.healthProfileRepository.CreateHealthProfile(ctx, *domainProfile); err != nil {
		ctx.Logger().Error("failed to insert profile to db", err, appcontext.Fields{})
		return nil, err
	}

	// respond
	ctx.Logger().Text("done create new health profile")
	return &dto.CreateHealthProfileResponse{ID: domainProfile.ID}, nil
}
