package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/Gleipnir-Technology/nidus-sync/platform"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

var connectionsSSE map[*ConnectionSSE]bool = make(map[*ConnectionSSE]bool, 0)

type ConnectionSSE struct {
	chanEvent      chan platform.Event
	id             uuid.UUID
	organizationID int32
	userID         int32
}

type Message struct {
	Resource string    `json:"resource"`
	Time     time.Time `json:"time"`
	Type     string    `json:"type"`
	URI      string    `json:"uri"`
}

func (c *ConnectionSSE) SendEvent(w http.ResponseWriter, m platform.Event) error {
	return send(w, Message{
		Resource: m.Resource,
		Time:     m.Time,
		Type:     m.Type.String(),
		URI:      m.URI,
	})
}
func (c *ConnectionSSE) SendHeartbeat(w http.ResponseWriter, t time.Time) error {
	return send(w, platform.Event{
		Resource: "clock",
		Time:     t,
		Type:     platform.EventTypeHeartbeat,
		URI:      "",
	})
}
func SetEventChannel(chan_envelopes <-chan platform.Envelope) {
	go func() {
		for envelope := range chan_envelopes {
			for conn, _ := range connectionsSSE {
				if conn.organizationID == envelope.OrganizationID {
					log.Debug().Int("type", int(envelope.Event.Type)).Int32("env-org", envelope.OrganizationID).Msg("pushed event to client")
					conn.chanEvent <- envelope.Event
				} else if conn.userID == envelope.UserID {
					log.Debug().Int("type", int(envelope.Event.Type)).Int32("env-user", envelope.UserID).Msg("pushed event to user")
					conn.chanEvent <- envelope.Event
				} else {
					log.Debug().Int("type", int(envelope.Event.Type)).Int32("env-org", envelope.OrganizationID).Int32("conn-org", conn.organizationID).Msg("skipped event, bad org")
				}

			}
		}
	}()
}
func send[T any](w http.ResponseWriter, msg T) error {
	jsonData, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("marshaling json: %w", err)
	}
	// Write in SSE format: "data: <json>\n\n"
	_, err = fmt.Fprintf(w, "data: %s\n\n", jsonData)
	if err != nil {
		return fmt.Errorf("writing SSE message: %w", err)
	}

	w.(http.Flusher).Flush()
	return nil
}
func streamEvents(w http.ResponseWriter, r *http.Request, u platform.User) {
	// Set headers for SSE
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	uid, err := uuid.NewUUID()
	if err != nil {
		log.Error().Err(err).Msg("failed to create uuid")
	}
	connection := ConnectionSSE{
		chanEvent:      make(chan platform.Event),
		id:             uid,
		organizationID: u.Organization.ID,
		userID:         int32(u.ID),
	}
	connectionsSSE[&connection] = true
	log.Debug().Int32("org", u.Organization.ID).Int("user", u.ID).Str("id", uid.String()).Msg("connected SSE client")

	// Send an initial connected event
	fmt.Fprintf(w, "event: connected\ndata: {\"status\": \"connected\", \"time\": \"%s\"}\n\n", time.Now().Format(time.RFC3339))
	w.(http.Flusher).Flush()

	// Keep the connection open with a ticker sending periodic events
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	// Use a channel to detect when the client disconnects
	done := r.Context().Done()

	// Keep connection open until client disconnects
	for {
		select {
		case <-done:
			log.Debug().Int32("org", u.Organization.ID).Int("user", u.ID).Str("id", uid.String()).Msg("Client closed connection")
			delete(connectionsSSE, &connection)
			return
		case t := <-ticker.C:
			// Send a heartbeat message
			err = connection.SendHeartbeat(w, t)
			if err != nil {
				log.Error().Err(err).Msg("Failed to send heartbeat")
			}
		case e := <-connection.chanEvent:
			err = connection.SendEvent(w, e)
			if err != nil {
				log.Error().Err(err).Msg("Failed to send heartbeat")
			}
		}
	}
}
