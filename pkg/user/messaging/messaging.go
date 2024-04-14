package messaging

import (
	"context"
	"github.com/namhq1989/bapbi-server/internal/utils/appcontext"
	"github.com/namhq1989/bapbi-server/pkg/user/infrastructure"
)

type Messaging struct {
	realtimeService infrastructure.RealtimeService
}

func New(realtimeService infrastructure.RealtimeService) Messaging {
	return Messaging{
		realtimeService: realtimeService,
	}
}

func (w Messaging) Start() {
	// subscribe new connection
	if err := w.realtimeService.SubscribeChannelNewConnection(appcontext.New(context.Background())); err != nil {
		panic(err)
	}
}
