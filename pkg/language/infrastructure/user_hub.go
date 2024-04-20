package infrastructure

import (
	"github.com/namhq1989/bapbi-server/internal/genproto/userpb"
	"github.com/namhq1989/bapbi-server/internal/utils/appcontext"
	apperrors "github.com/namhq1989/bapbi-server/internal/utils/error"
	"github.com/namhq1989/bapbi-server/pkg/language/domain"
)

type UserHub struct {
	client userpb.UserServiceClient
}

func NewUserHub(client userpb.UserServiceClient) *UserHub {
	return &UserHub{
		client: client,
	}
}

func (r UserHub) GetUserPlan(ctx *appcontext.AppContext, userID string) (*domain.SubscriptionPlan, error) {
	resp, err := r.client.GetUserSubscriptionPlan(ctx.Context(), &userpb.GetUserSubscriptionPlanRequest{
		Id: userID,
	})
	if err != nil {
		ctx.Logger().Error("failed to get user subscription plan", err, appcontext.Fields{"id": userID})
		return nil, err
	}

	plan := domain.ToSubscriptionPlan(resp.GetPlan())
	if !plan.IsValid() {
		return nil, apperrors.User.InvalidSubscriptionPlan
	}

	return &plan, nil
}
