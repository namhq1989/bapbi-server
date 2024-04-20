package hub

import (
	"github.com/namhq1989/bapbi-server/internal/genproto/userpb"
	"github.com/namhq1989/bapbi-server/internal/utils/appcontext"
	"github.com/namhq1989/bapbi-server/pkg/user/domain"
)

type GetUserSubscriptionPlanHandler struct {
	userHub domain.UserHub
}

func NewGetUserSubscriptionPlanHandler(userHub domain.UserHub) GetUserSubscriptionPlanHandler {
	return GetUserSubscriptionPlanHandler{
		userHub: userHub,
	}
}

func (h GetUserSubscriptionPlanHandler) GetUserSubscriptionPlan(ctx *appcontext.AppContext, req *userpb.GetUserSubscriptionPlanRequest) (*userpb.GetUserSubscriptionPlanResponse, error) {
	user, err := h.userHub.FindOneByID(ctx, req.GetId())
	if err != nil {
		return nil, err
	}

	if user == nil {
		return &userpb.GetUserSubscriptionPlanResponse{Plan: ""}, nil
	}

	return &userpb.GetUserSubscriptionPlanResponse{
		Plan: user.SubscriptionPlan.String(),
	}, nil
}
