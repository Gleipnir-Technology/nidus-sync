package platform

import (
	"time"

	"github.com/Gleipnir-Technology/nidus-sync/config"
	"github.com/Gleipnir-Technology/nidus-sync/platform/event"
)

type Envelope = event.Envelope
type Event = event.Event

const EventTypeHeartbeat = event.EventTypeHeartbeat

func EventShutdown() {
	event.Shutdown()
}
func SetEventChannel(chan_events chan<- Envelope) {
	event.SetEventChannel(chan_events)
}
func SudoEvent(org_id int32, resource, type_, uri_path string) {
	event_type := event.EventTypeFromString(type_)
	go event.Send(event.Envelope{
		Event: Event{
			Resource: resource,
			Time:     time.Now(),
			Type:     event_type,
			URI:      config.MakeURLNidus(uri_path),
		},
		OrganizationID: org_id,
	})
}
