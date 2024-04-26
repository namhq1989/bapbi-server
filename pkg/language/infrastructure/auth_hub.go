package infrastructure

import (
	"github.com/namhq1989/bapbi-server/internal/genproto/authpb"
	"github.com/namhq1989/bapbi-server/internal/utils/appcontext"
)

type AuthHub struct {
	client authpb.AuthServiceClient
}

func NewAuthHub(client authpb.AuthServiceClient) *AuthHub {
	return &AuthHub{
		client: client,
	}
}

func (r AuthHub) IsAdmin(ctx *appcontext.AppContext, userID string) (bool, error) {
	resp, err := r.client.IsAdmin(ctx.Context(), &authpb.IsAdminRequest{
		Id: userID,
	})
	if err != nil {
		ctx.Logger().Error("failed to check user is admin or not", err, appcontext.Fields{"id": userID})
		return false, err
	}

	return resp.GetIsAdmin(), nil
}
