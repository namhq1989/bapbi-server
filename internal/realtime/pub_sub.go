package realtime

import (
	"context"

	"github.com/ably/ably-go/ably"
	"github.com/namhq1989/bapbi-server/internal/utils/appcontext"
	apperrors "github.com/namhq1989/bapbi-server/internal/utils/error"
)

func (r Realtime) SubscribeEvent(ctx *appcontext.AppContext, channel, event string, handler func(Message)) error {
	c := r.ably.Channels.Get(channel)
	if c == nil {
		return apperrors.Realtime.ChannelNotFound
	}

	_, err := c.Subscribe(ctx.Context(), event, func(message *ably.Message) {
		handler(Message{
			ID:       message.ID,
			ClientID: message.ClientID,
			Name:     message.Name,
			Data:     message.Data,
		})
	})
	if err != nil {
		ctx.Logger().Error("failed to subscribe event", err, appcontext.Fields{"channel": channel, "event": event})
		return apperrors.Realtime.CannotSubscribeEvent
	}

	return nil
}

func (r Realtime) PublishEvent(ctx *appcontext.AppContext, channel, event string, data interface{}) error {
	c := r.ably.Channels.Get(channel)
	if c == nil {
		return apperrors.Realtime.ChannelNotFound
	}

	if err := c.Publish(context.Background(), event, data); err != nil {
		ctx.Logger().Error("failed to publish event", err, appcontext.Fields{"channel": channel, "event": event})
		return apperrors.Realtime.CannotPublishEvent
	}
	return nil
}
