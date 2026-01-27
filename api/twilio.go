package api

import (
	"fmt"
	"net/http"

	"github.com/Gleipnir-Technology/nidus-sync/platform/text"
	"github.com/rs/zerolog/log"
	"github.com/twilio/twilio-go/twiml"
)

func twilioMessagePost(w http.ResponseWriter, r *http.Request) {
	message_sid := r.PostFormValue("MessageSid")
	log.Info().Str("sid", message_sid).Msg("Twilio Message POST")
	fmt.Fprintf(w, "")
}
func twilioStatusPost(w http.ResponseWriter, r *http.Request) {
	message_sid := r.PostFormValue("MessageSid")
	message_status := r.PostFormValue("MessageStatus")
	log.Info().Str("sid", message_sid).Str("status", message_status).Msg("Updated message status")
	text.UpdateMessageStatus(message_sid, message_status)
	fmt.Fprintf(w, "")
}
func twilioTextPost(w http.ResponseWriter, r *http.Request) {
	message_sid := r.PostFormValue("MessageSid")
	account_sid := r.PostFormValue("AccountSid")
	messaging_service_sid := r.PostFormValue("MessagingServiceSid")
	from := r.PostFormValue("From")
	to_ := r.PostFormValue("To")
	body := r.PostFormValue("Body")
	num_media := r.PostFormValue("NumMedia")
	num_segments := r.PostFormValue("NumSegments")
	media_content_type0 := r.PostFormValue("MediaContentType0")
	media_url0 := r.PostFormValue("MediaUrl0")
	from_city := r.PostFormValue("FromCity")
	from_state := r.PostFormValue("FromState")
	from_zip := r.PostFormValue("FromZip")
	from_country := r.PostFormValue("FromCountry")
	to_city := r.PostFormValue("ToCity")
	to_state := r.PostFormValue("ToState")
	to_zip := r.PostFormValue("ToZip")
	to_country := r.PostFormValue("ToCountry")
	log.Info().Str("message_sid", message_sid).Str("account_sid", account_sid).Str("messaging_service_sid", messaging_service_sid).Str("from", from).Str("to_", to_).Str("body", body).Str("num_media", num_media).Str("num_segments", num_segments).Str("media_content_type0", media_content_type0).Str("media_url0", media_url0).Str("from_city", from_city).Str("from_state", from_state).Str("from_zip", from_zip).Str("from_country", from_country).Str("to_city", to_city).Str("to_state", to_state).Str("to_zip", to_zip).Str("to_country", to_country).Msg("got text")

	twiml, _ := twiml.Messages([]twiml.Element{})
	go text.HandleTextMessage(from, to_, body)
	w.Header().Set("Content-Type", "text/xml")
	fmt.Fprintf(w, "%s", twiml)
}
