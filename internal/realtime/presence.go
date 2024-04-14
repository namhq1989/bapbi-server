package realtime

import (
	"fmt"
	"time"

	"github.com/ably/ably-go/ably"
	"github.com/namhq1989/bapbi-server/internal/utils/appcontext"
	apperrors "github.com/namhq1989/bapbi-server/internal/utils/error"
)

func (r Realtime) SubscribePresenceAllEvents(ctx *appcontext.AppContext, channel string, handler func(Message)) error {
	_, err := r.ably.Channels.Get(channel).Presence.SubscribeAll(ctx.Context(), func(message *ably.PresenceMessage) {
		handler(Message{
			ID:       message.ID,
			ClientID: message.ClientID,
			Name:     message.Name,
			Action:   message.Action.String(),
			Data:     message.Data,
		})
	})
	if err != nil {
		ctx.Logger().Error("failed to subscribe presence channel", err, appcontext.Fields{"channel": channel})
		return apperrors.Realtime.CannotSubscribeEvent
	}

	go func() {
		for {
			clients, _ := r.ably.Channels.Get(channel).Presence.Get(ctx.Context())
			for _, client := range clients {
				fmt.Println("Present client:", client)
			}

			time.Sleep(5 * time.Second)
		}
	}()

	return nil
}
