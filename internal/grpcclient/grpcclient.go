package grpcclient

import (
	"github.com/namhq1989/bapbi-server/internal/genproto/authpb"
	"github.com/namhq1989/bapbi-server/internal/genproto/userpb"
	"github.com/namhq1989/bapbi-server/internal/utils/appcontext"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func NewUserClient(ctx *appcontext.AppContext, addr string) (userpb.UserServiceClient, error) {
	conn, err := grpc.DialContext(ctx.Context(), addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	return userpb.NewUserServiceClient(conn), nil
}

func NewAuthClient(ctx *appcontext.AppContext, addr string) (authpb.AuthServiceClient, error) {
	conn, err := grpc.DialContext(ctx.Context(), addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	return authpb.NewAuthServiceClient(conn), nil
}
