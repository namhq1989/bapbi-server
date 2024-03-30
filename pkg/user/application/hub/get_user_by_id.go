package hub

import (
	"github.com/namhq1989/bapbi-server/internal/genproto/userpb"
	"github.com/namhq1989/bapbi-server/internal/utils/appcontext"
	"github.com/namhq1989/bapbi-server/pkg/user/domain"
)

type GetUserByIDHandler struct {
	userHub domain.UserHub
}

func NewGetUserByIDHandler(userHub domain.UserHub) GetUserByIDHandler {
	return GetUserByIDHandler{
		userHub: userHub,
	}
}

func (h GetUserByIDHandler) GetUserByID(ctx *appcontext.AppContext, req *userpb.GetUserByIDRequest) (*userpb.GetUserByIDResponse, error) {
	user, err := h.userHub.FindOneByID(ctx, req.GetId())
	if err != nil {
		return nil, err
	}

	return &userpb.GetUserByIDResponse{
		User: &userpb.User{
			Id:    user.ID,
			Name:  user.Name,
			Email: user.Email,
		},
	}, nil
}
