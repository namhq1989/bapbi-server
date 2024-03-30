package query

import (
	"github.com/namhq1989/bapbi-server/internal/utils/appcontext"
	apperrors "github.com/namhq1989/bapbi-server/internal/utils/error"
	"github.com/namhq1989/bapbi-server/pkg/auth/domain"
	"github.com/namhq1989/bapbi-server/pkg/auth/dto"
)

type MeHandler struct {
	userHub domain.UserHub
}

func NewMeHandler(userHub domain.UserHub) MeHandler {
	return MeHandler{
		userHub: userHub,
	}
}

func (h MeHandler) Me(ctx *appcontext.AppContext, req dto.MeRequest) (*dto.MeResponse, error) {
	ctx.Logger().Text("get me")

	// get user id
	userID := ctx.GetUserID()

	// find user
	ctx.Logger().Info("find user in database", appcontext.Fields{"userId": userID})
	user, err := h.userHub.GetOneByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		ctx.Logger().Error("user not found", nil, appcontext.Fields{"userId": userID})
		return nil, apperrors.Common.NotFound
	}

	// respond
	ctx.Logger().Text("respond user data")
	return &dto.MeResponse{
		ID:   user.ID,
		Name: user.Name,
	}, nil
}
