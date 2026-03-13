package event

import (
	"github.com/Gleipnir-Technology/nidus-sync/config"
	"time"
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
	EventTypeModified
	EventTypeHeartbeat
	EventTypeSudo
)

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
