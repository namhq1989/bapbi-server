package grpc

import (
	"context"

	"github.com/namhq1989/bapbi-server/internal/utils/appcontext"

	"github.com/namhq1989/bapbi-server/internal/genproto/userpb"
)

func (s server) GetUserByID(bgCtx context.Context, req *userpb.GetUserByIDRequest) (*userpb.GetUserByIDResponse, error) {
	var (
		ctx = appcontext.New(bgCtx)
	)

	return s.app.GetUserByID(ctx, req)
}

func (s server) GetUserByEmail(bgCtx context.Context, req *userpb.GetUserByEmailRequest) (*userpb.GetUserByEmailResponse, error) {
	var (
		ctx = appcontext.New(bgCtx)
	)

	return s.app.GetUserByEmail(ctx, req)
}

func (s server) CreateUser(bgCtx context.Context, req *userpb.CreateUserRequest) (*userpb.CreateUserResponse, error) {
	var (
		ctx = appcontext.New(bgCtx)
	)

	return s.app.CreateUser(ctx, req)
}
