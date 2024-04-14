package infrastructure

import (
	"github.com/namhq1989/bapbi-server/internal/realtime"
	"github.com/namhq1989/bapbi-server/internal/utils/appcontext"
	"github.com/namhq1989/bapbi-server/pkg/user/domain"
)

type RealtimeService struct {
	realtime *realtime.Realtime
}

func NewRealtimeService(rt *realtime.Realtime) RealtimeService {
	return RealtimeService{
		realtime: rt,
	}
}

func (r RealtimeService) SubscribeChannelNewConnection(ctx *appcontext.AppContext) error {
	if err := r.realtime.SubscribeEvent(ctx, domain.RealtimeServiceChannels.NewConnection, "", func(m realtime.Message) {
		ctx.Logger().Print("new message from channel new_connection", appcontext.Fields{
			"id":       m.ID,
			"clientId": m.ClientID,
			"name":     m.Name,
			"data":     m.Data,
		})
	}); err != nil {
		return err
	}

	if err := r.realtime.SubscribePresenceAllEvents(ctx, domain.RealtimeServiceChannels.NewConnection, func(m realtime.Message) {
		ctx.Logger().Print("new presence from channel new_connection", appcontext.Fields{
			"id":       m.ID,
			"clientId": m.ClientID,
			"name":     m.Name,
			"data":     m.Data,
			"action":   m.Action,
		})
	}); err != nil {
		return err
	}

	return nil
}
