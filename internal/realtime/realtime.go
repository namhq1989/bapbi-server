package realtime

import (
	"fmt"

	"github.com/ably/ably-go/ably"
)

type Realtime struct {
	ably *ably.Realtime
}

type Message struct {
	ID       string
	ClientID string
	Name     string
	Action   string
	Data     interface{}
}

func NewRealtimeClient(apiKey string) *Realtime {
	client, err := ably.NewRealtime(ably.WithKey(apiKey), ably.WithAutoConnect(true))
	if err != nil {
		panic(fmt.Errorf("failed to create ably client: %s", err.Error()))
	}

	fmt.Printf("⚡️ [realtime]: connected \n")

	return &Realtime{ably: client}
}
