package grpc

import (
	"github.com/namhq1989/bapbi-server/internal/genproto/userpb"
	"github.com/namhq1989/bapbi-server/internal/utils/appcontext"
	"github.com/namhq1989/bapbi-server/pkg/user/application"
	"google.golang.org/grpc"
)

type server struct {
	app application.App
	userpb.UnimplementedUserServiceServer
}

var _ userpb.UserServiceServer = (*server)(nil)

func RegisterServer(_ *appcontext.AppContext, app application.App, registrar grpc.ServiceRegistrar) error {
	userpb.RegisterUserServiceServer(registrar, server{app: app})
	return nil
}
