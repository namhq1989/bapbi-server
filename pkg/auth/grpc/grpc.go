package grpc

import (
	"github.com/namhq1989/bapbi-server/internal/genproto/authpb"
	"github.com/namhq1989/bapbi-server/internal/utils/appcontext"
	"github.com/namhq1989/bapbi-server/pkg/auth/application"
	"google.golang.org/grpc"
)

type server struct {
	app application.App
	authpb.UnimplementedAuthServiceServer
}

var _ authpb.AuthServiceServer = (*server)(nil)

func RegisterServer(_ *appcontext.AppContext, app application.App, registrar grpc.ServiceRegistrar) error {
	authpb.RegisterAuthServiceServer(registrar, server{app: app})
	return nil
}
