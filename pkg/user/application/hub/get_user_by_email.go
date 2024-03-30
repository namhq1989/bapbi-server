package hub

import (
	"github.com/namhq1989/bapbi-server/internal/genproto/userpb"
	"github.com/namhq1989/bapbi-server/internal/utils/appcontext"
	"github.com/namhq1989/bapbi-server/pkg/user/domain"
)

type GetUserByEmailHandler struct {
	userHub domain.UserHub
}

func NewGetUserByEmailHandler(userHub domain.UserHub) GetUserByEmailHandler {
	return GetUserByEmailHandler{
		userHub: userHub,
	}
}

func (h GetUserByEmailHandler) GetUserByEmail(ctx *appcontext.AppContext, req *userpb.GetUserByEmailRequest) (*userpb.GetUserByEmailResponse, error) {
	user, err := h.userHub.FindOneByEmail(ctx, req.GetEmail())
	if err != nil {
		return nil, err
	}

	return &userpb.GetUserByEmailResponse{
		User: &userpb.User{
			Id:    user.ID,
			Name:  user.Name,
			Email: user.Email,
		},
	}, nil
}
