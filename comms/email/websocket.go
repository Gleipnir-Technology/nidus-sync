package email

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/gorilla/websocket"
	"github.com/rs/zerolog/log"
)

var FORWARDEMAIL_WS_API = "wss://api.forwardemail.net/v1/ws"

func StartWebsocket(ctx context.Context, api_token string) {

	var conn *websocket.Conn
	for {
		err := ensureConnected(conn, api_token)
		if err != nil {
			log.Error().Err(err).Msg("Bailing on email websocket")
			return
		}
		select {
		case <-ctx.Done():
			return
		default:
			// Read message
			message_type, message, err := conn.ReadMessage()
			if err != nil {
				if !websocket.IsCloseError(err, websocket.CloseNormalClosure) {
					conn = nil
				}
				log.Error().Err(err).Msg("Error reading message")
			}

			// Process and log the message
			log.Info().Int("message_type", message_type).Bytes("message", message).Msg("Got email notification")
		}
	}
}

func ensureConnected(conn *websocket.Conn, api_token string) error {
	if conn != nil {
		return nil
	}
	url := FORWARDEMAIL_WS_API + "?token=" + api_token
	for {
		new_conn, _, err := websocket.DefaultDialer.Dial(url, nil)
		if err == nil {
			log.Info().Msg("Connected to mail websocket")
			*conn = *new_conn
			return nil
		}
		if errors.Is(err, websocket.ErrBadHandshake) {
			return fmt.Errorf("Bad handshake connecting to email websocket, bailing.")
		}
		log.Error().Err(err).Str("url", url).Msg("Error connecting to WebSocket")
		time.Sleep(3 * time.Second)

	}
}
