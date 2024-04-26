package hub

import (
	"github.com/namhq1989/bapbi-server/internal/genproto/authpb"
	"github.com/namhq1989/bapbi-server/internal/utils/appcontext"
	"github.com/namhq1989/bapbi-server/pkg/auth/domain"
)

type IsAdminHandler struct {
	userHub domain.UserHub
}

func NewIsAdminHandler(userHub domain.UserHub) IsAdminHandler {
	return IsAdminHandler{
		userHub: userHub,
	}
}

func (h IsAdminHandler) IsAdmin(ctx *appcontext.AppContext, req *authpb.IsAdminRequest) (*authpb.IsAdminResponse, error) {
	user, err := h.userHub.GetOneByID(ctx, req.GetId())
	if err != nil {
		return nil, err
	}

	return &authpb.IsAdminResponse{IsAdmin: user.IsAdmin}, nil
}
