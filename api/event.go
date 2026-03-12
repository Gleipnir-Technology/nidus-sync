package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/Gleipnir-Technology/nidus-sync/platform"
	"github.com/rs/zerolog/log"
)

func streamEvents(w http.ResponseWriter, r *http.Request, u platform.User) {
}

type MessageHeartbeat struct {
	Time time.Time `json:"time"`
}
type MessageSSE struct {
	Content any    `json:"content"`
	Type    string `json:"type"`
}
type ConnectionSSE struct {
	chanState chan MessageSSE
	id        string
}

func (c *ConnectionSSE) SendMessage(w http.ResponseWriter, m MessageSSE) error {
	return send(w, MessageSSE{
		Type: "heartbeat",
	})
}
func (c *ConnectionSSE) SendHeartbeat(w http.ResponseWriter, t time.Time) error {
	return send(w, MessageSSE{
		Content: MessageHeartbeat{
			Time: t,
		},
		Type: "heartbeat",
	})
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

type Webserver struct {
	connections map[*ConnectionSSE]bool
}

// sseHandler handles the Server-Sent Events connection
func (web *Webserver) sseHandler(w http.ResponseWriter, r *http.Request) {
	// Set headers for SSE
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	connection := ConnectionSSE{
		chanState: make(chan MessageSSE),
		id:        fmt.Sprintf("%d", time.Now().UnixNano()),
	}
	web.connections[&connection] = true
	// Send an initial connected event
	fmt.Fprintf(w, "event: connected\ndata: {\"status\": \"connected\", \"time\": \"%s\"}\n\n", time.Now().Format(time.RFC3339))
	w.(http.Flusher).Flush()

	// Keep the connection open with a ticker sending periodic events
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	// Use a channel to detect when the client disconnects
	done := r.Context().Done()

	// Keep connection open until client disconnects
	var err error
	for {
		err = nil
		select {
		case <-done:
			log.Info().Msg("Client closed connection")
			return
		case t := <-ticker.C:
			// Send a heartbeat message
			err = connection.SendHeartbeat(w, t)
			//case state := <-connection.chanState:
			//log.Debug().Msg("Sending new state to connection")
			//err = connection.SendState(w, state)
		}
		if err != nil {
			log.Error().Err(err).Msg("Failed to send state from webserver")
		}
	}
}
