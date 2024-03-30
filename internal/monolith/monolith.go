package monolith

import (
	"github.com/labstack/echo/v4"
	"github.com/namhq1989/bapbi-server/internal/caching"
	"github.com/namhq1989/bapbi-server/internal/config"
	"github.com/namhq1989/bapbi-server/internal/monitoring"
	"github.com/namhq1989/bapbi-server/internal/queue"
	"github.com/namhq1989/bapbi-server/internal/utils/appcontext"
	appjwt "github.com/namhq1989/bapbi-server/internal/utils/jwt"
	"github.com/namhq1989/bapbi-server/internal/utils/waiter"
	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/grpc"
)

type Monolith interface {
	Config() config.Server
	Mongo() *mongo.Database
	Rest() *echo.Echo
	RPC() *grpc.Server
	Waiter() waiter.Waiter
	JWT() *appjwt.JWT
	Caching() *caching.Caching
	Monitoring() *monitoring.Monitoring
	Queue() *queue.Queue
}

type Module interface {
	Name() string
	Startup(ctx *appcontext.AppContext, monolith Monolith) error
}
