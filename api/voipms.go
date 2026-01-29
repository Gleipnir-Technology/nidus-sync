package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	//"github.com/Gleipnir-Technology/nidus-sync/config"
	"github.com/Gleipnir-Technology/nidus-sync/platform/text"
	"github.com/rs/zerolog/log"
)

/*
	{
	  "data": {
	    "id": 101252305,
	    "event_type": "message.received",
	    "record_type": "event",
	    "payload": {
	      "id": 101252305,
	      "record_type": "message",
	      "from": {
		"phone_number": "+18016984649"
	      },
	      "to": [
		{
		  "phone_number": "+15593720139",
		  "status": "webhook_delivered"
		}
	      ],
	      "text": "test3",
	      "received_at": "2026-01-29T20:16:23.000000+00:00",
	      "type": "SMS",
	      "media": []
	    }
	  }
	}
*/
type VoipMSStatusPhoneFrom struct {
	PhoneNumber string `json:"phone_number"`
}
type VoipMSStatusPhoneTo struct {
	PhoneNumber string `json:"phone_number"`
	Status      string `json:"status"`
}
type VoipMSStatusPayload struct {
	ID         int                   `json:"id"`
	RecordType string                `json:"record_type"`
	From       VoipMSStatusPhoneFrom `json:"from"`
	To         []VoipMSStatusPhoneTo `json:"to"`
	Text       string                `json:"text"`
	ReceivedAt string                `json:"received_at"`
	Type       string                `json:"type"`
	//Media []something
}
type VoipMSStatusUpdate struct {
	ID         int                 `json:"id"`
	EventType  string              `json:"event_type"`
	RecordType string              `json:"record_type"`
	Payload    VoipMSStatusPayload `json:"payload"`
}
type VoipMSTextPostBody struct {
	Data VoipMSStatusUpdate `json:"data"`
}

func voipmsTextGet(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	name := query.Get("to")
	age := query.Get("from")
	message := query.Get("message")
	files := query.Get("files")
	id := query.Get("id")
	date := query.Get("date")
	log.Info().Str("name", name).Str("age", age).Str("message", message).Str("files", files).Str("id", id).Str("date", date).Msg("Incoming text message")
}
func voipmsTextPost(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "failed to read", http.StatusInternalServerError)
		return
	}
	//debugSaveRequest(r)
	var b VoipMSTextPostBody
	err = json.Unmarshal(body, &b)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	to := "unknown"
	if len(b.Data.Payload.To) > 0 {
		to = b.Data.Payload.To[0].PhoneNumber
	}
	log.Info().Int("ID", b.Data.ID).Str("event_type", b.Data.EventType).Str("record_type", b.Data.RecordType).Str("from", b.Data.Payload.From.PhoneNumber).Str("to", to).Str("content", b.Data.Payload.Text).Msg("Text status")

	// Convert phone numbers from Voip.ms into E164 format for consistency
	go text.HandleTextMessage(b.Data.Payload.From.PhoneNumber, to, b.Data.Payload.Text)
	fmt.Fprintf(w, "ok")
}
