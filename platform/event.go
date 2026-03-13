package platform

import (
	"time"

	"github.com/Gleipnir-Technology/nidus-sync/platform/event"
)

type Envelope = event.Envelope
type Event = event.Event

const EventTypeHeartbeat = event.EventTypeHeartbeat

func SetEventChannel(chan_events chan<- Envelope) {
	event.SetEventChannel(chan_events)
}
func SudoEvent(org_id int32, content string) {
	go event.Send(event.Envelope{
		Event: Event{
			Resource: "sudo",
			Time:     time.Now(),
			Type:     event.EventTypeSudo,
			URI:      content,
		},
		OrganizationID: org_id,
	})
}
