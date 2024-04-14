package apperrors

import "errors"

var Realtime = struct {
	ChannelNotFound      error
	CannotSubscribeEvent error
	CannotPublishEvent   error
}{
	ChannelNotFound:      errors.New("realtime_channel_not_found"),
	CannotSubscribeEvent: errors.New("realtime_cannot_subscribe_event"),
	CannotPublishEvent:   errors.New("realtime_cannot_publish_event"),
}
