package event

import (
	"time"

	"github.com/Gleipnir-Technology/nidus-sync/config"
)

var chanEvents chan<- Envelope

type Event struct {
	Resource string    `json:"resource"`
	Time     time.Time `json:"time"`
	Type     EventType `json:"type"`
	URI      string    `json:"uri"`
}
type Envelope struct {
	OrganizationID int32
	Event          Event
}

func SetEventChannel(chan_events chan<- Envelope) {
	chanEvents = chan_events
}

type EventType int

const (
	EventTypeCreated EventType = iota
	EventTypeDeleted
	EventTypeHeartbeat
	EventTypeSudo
	EventTypeUnknown
	EventTypeUpdated
)

func (et EventType) String() string {
	switch et {
	case EventTypeCreated:
		return "created"
	case EventTypeDeleted:
		return "deleted"
	case EventTypeHeartbeat:
		return "heartbeat"
	case EventTypeSudo:
		return "sudo"
	case EventTypeUnknown:
		return "unknown"
	case EventTypeUpdated:
		return "updated"
	}
	return "unknown"
}
func EventTypeFromString(s string) EventType {
	switch s {
	case "created":
		return EventTypeCreated
	case "deleted":
		return EventTypeDeleted
	case "heartbeat":
		return EventTypeHeartbeat
	case "sudo":
		return EventTypeSudo
	case "updated":
		return EventTypeUpdated
	default:
		return EventTypeUnknown
	}
}

type ResourceType int

const (
	TypeRMONuisance = iota
	TypeRMOWater
)

func Created(type_ ResourceType, organization_id int32, uri_id string) {
	var resource string
	var uri string
	switch type_ {
	case TypeRMONuisance:
		resource = "rmo:nuisance"
		uri = config.MakeURLReport("/report/%s", uri_id)
	case TypeRMOWater:
		resource = "rmo:water"
		uri = config.MakeURLReport("/report/%s", uri_id)
	default:

	}
	go Send(Envelope{
		Event: Event{
			Resource: resource,
			Time:     time.Now(),
			Type:     EventTypeCreated,
			URI:      uri,
		},
		OrganizationID: organization_id,
	})
}
func Send(env Envelope) {
	chanEvents <- env

}
