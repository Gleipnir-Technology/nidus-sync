package event

import (
	"encoding/json"
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

func (e Event) MarshalJSON() ([]byte, error) {
	to_marshal := make(map[string]any, 0)
	to_marshal["resource"] = e.Resource
	to_marshal["time"] = e.Time
	to_marshal["type"] = e.Type.String()
	to_marshal["uri"] = e.URI
	return json.Marshal(to_marshal)
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
	TypeUnknown = iota
	TypeFileCSV
	TypeNoteAudio
	TypeNoteImage
	TypeReviewTask
	TypeRMONuisance
	TypeRMOReport
	TypeRMOWater
	TypeSignal
)

func Created(t ResourceType, organization_id int32, uri_id string) {
	go Send(Envelope{
		Event: Event{
			Resource: resourceString(t),
			Time:     time.Now(),
			Type:     EventTypeCreated,
			URI:      makeURI(t, uri_id),
		},
		OrganizationID: organization_id,
	})
}
func Updated(t ResourceType, organization_id int32, uri_id string) {
	go Send(Envelope{
		Event: Event{
			Resource: resourceString(t),
			Time:     time.Now(),
			Type:     EventTypeUpdated,
			URI:      makeURI(t, uri_id),
		},
		OrganizationID: organization_id,
	})
}
func Send(env Envelope) {
	chanEvents <- env
}
func resourceString(t ResourceType) string {
	switch t {
	case TypeFileCSV:
		return "sync:filecsv"
	case TypeNoteAudio:
		return "sync:note:audio"
	case TypeNoteImage:
		return "sync:note:image"
	case TypeReviewTask:
		return "sync:review_task"
	case TypeRMONuisance:
		return "rmo:nuisance"
	case TypeRMOReport:
		return "rmo:report"
	case TypeRMOWater:
		return "rmo:water"
	case TypeSignal:
		return "sync:signal"
	default:
		return "unknown"
	}
}
func makeURI(t ResourceType, id string) string {
	switch t {
	case TypeFileCSV:
		return config.MakeURLNidus("/upload/%s", id)
	case TypeNoteAudio:
		return config.MakeURLNidus("/note/%s", id)
	case TypeNoteImage:
		return config.MakeURLNidus("/note/%s", id)
	case TypeReviewTask:
		return config.MakeURLNidus("/review/%s", id)
	case TypeRMONuisance:
		return config.MakeURLReport("/report/%s", id)
	case TypeRMOWater:
		return config.MakeURLReport("/report/%s", id)
	case TypeSignal:
		return config.MakeURLReport("/signal/%s", id)
	default:
		return config.MakeURLReport("/unknown")
	}
}
