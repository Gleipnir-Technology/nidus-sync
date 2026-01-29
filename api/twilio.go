package api

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/Gleipnir-Technology/nidus-sync/config"
	"github.com/Gleipnir-Technology/nidus-sync/platform/text"
	"github.com/rs/zerolog/log"
	"github.com/twilio/twilio-go/twiml"
)

// Translate from Twilio's representation of a RCS message sender to our concept of a phone number
// From: rcs:dev_report_mosquitoes_online_dosrvwxm_agent
// To: +16235525879
func getDst(to string) (string, error) {

	if to == config.TwilioRCSSenderRMO {
		return config.PhoneNumberReportStr, nil
	}
	/*
		phone, err := models.FindCommsPhone(ctx, db.PGInstance.BobDB, to)
		if err != nil {
			return "", fmt.Errorf("Failed to search for dest phone %s: %w", to, err)
		}
		return phone.E164, nil
	*/
	return "", fmt.Errorf("Cannot match phone number to '%s'", to)
}

func splitPhoneSource(s string) (string, string) {
	parts := strings.Split(s, ":")
	switch len(parts) {
	case 0:
		return "this isn't", "possible"
	case 1:
		return "", s
	case 2:
		return parts[0], parts[1]
	default:
		log.Warn().Str("s", s).Msg("Got an incomprehensible number of parts of a phone number")
		return parts[0], parts[1]
	}

}

func twilioMessagePost(w http.ResponseWriter, r *http.Request) {
	message_sid := r.PostFormValue("MessageSid")
	log.Info().Str("sid", message_sid).Msg("Twilio Message POST")
	fmt.Fprintf(w, "")
}
func twilioCallPost(w http.ResponseWriter, r *http.Request) {
	called := r.PostFormValue("Called")
	tostate := r.PostFormValue("ToState")
	callercountry := r.PostFormValue("CallerCountry")
	direction := r.PostFormValue("Direction")
	callerstate := r.PostFormValue("CallerState")
	tozip := r.PostFormValue("ToZip")
	callsid := r.PostFormValue("CallSid")
	to := r.PostFormValue("To")
	callerzip := r.PostFormValue("CallerZip")
	tocountry := r.PostFormValue("ToCountry")
	stirverstat := r.PostFormValue("StirVerstat")
	//calltoken := r.PostFormValue("CallToken")
	calledzip := r.PostFormValue("CalledZip")
	apiversion := r.PostFormValue("ApiVersion")
	calledcity := r.PostFormValue("CalledCity")
	callstatus := r.PostFormValue("CallStatus")
	from := r.PostFormValue("From")
	accountsid := r.PostFormValue("AccountSid")
	calledcountry := r.PostFormValue("CalledCountry")
	callercity := r.PostFormValue("CallerCity")
	tocity := r.PostFormValue("ToCity")
	fromcountry := r.PostFormValue("FromCountry")
	caller := r.PostFormValue("Caller")
	fromcity := r.PostFormValue("FromCity")
	calledstate := r.PostFormValue("CalledState")
	fromzip := r.PostFormValue("FromZip")
	fromstate := r.PostFormValue("FromState")
	log.Info().Str("called", called).Str("tostate", tostate).Str("callercountry", callercountry).Str("direction", direction).Str("callerstate", callerstate).Str("tozip", tozip).Str("callsid", callsid).Str("to", to).Str("callerzip", callerzip).Str("tocountry", tocountry).Str("stirverstat", stirverstat).Str("calledzip", calledzip).Str("apiversion", apiversion).Str("calledcity", calledcity).Str("callstatus", callstatus).Str("from", from).Str("accountsid", accountsid).Str("calledcountry", calledcountry).Str("callercity", callercity).Str("tocity", tocity).Str("fromcountry", fromcountry).Str("caller", caller).Str("fromcity", fromcity).Str("calledstate", calledstate).Str("fromzip", fromzip).Str("fromstate", fromstate).Msg("Incoming phone call")

	say := &twiml.VoiceSay{
		Message: "Thanks for calling Report Mosquitoes Online. I'll forward you to our tech support lead, Eli",
	}
	call := &twiml.VoiceDial{
		Number: config.PhoneNumberSupportStr,
	}
	twimlResult, err := twiml.Voice([]twiml.Element{say, call})
	if err != nil {
		log.Error().Err(err).Msg("Failed to produce TWIML")
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.Header().Set("Content-Type", "text/xml")
	fmt.Fprintf(w, "%s", twimlResult)
}

func twilioCallStatusPost(w http.ResponseWriter, r *http.Request) {
	call_sid := r.PostFormValue("CallSid")
	account_sid := r.PostFormValue("AccountSid")
	from := r.PostFormValue("From")
	to := r.PostFormValue("To")
	call_status := r.PostFormValue("CallStatus")
	api_version := r.PostFormValue("ApiVersion")
	direction := r.PostFormValue("Direction")
	forwarded_from := r.PostFormValue("ForwardedFrom")
	caller_name := r.PostFormValue("CallerName")
	parent_call_sid := r.PostFormValue("ParentCallSid")
	log.Info().Str("call_sid", call_sid).Str("account_sid", account_sid).Str("from", from).Str("to", to).Str("call_status", call_status).Str("api_version", api_version).Str("direction", direction).Str("forwarded_from", forwarded_from).Str("caller_name", caller_name).Str("parent_call_sid", parent_call_sid)
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
	type_, src := splitPhoneSource(from)
	log.Info().Str("message_sid", message_sid).Str("account_sid", account_sid).Str("messaging_service_sid", messaging_service_sid).Str("from", from).Str("to_", to_).Str("body", body).Str("num_media", num_media).Str("num_segments", num_segments).Str("media_content_type0", media_content_type0).Str("media_url0", media_url0).Str("from_city", from_city).Str("from_state", from_state).Str("from_zip", from_zip).Str("from_country", from_country).Str("to_city", to_city).Str("to_state", to_state).Str("to_zip", to_zip).Str("to_country", to_country).Str("type_", type_).Msg("got text")

	twiml, _ := twiml.Messages([]twiml.Element{})

	dst, err := getDst(to_)
	if err != nil {
		log.Error().Err(err).Str("to", to_).Msg("Failed to get dst")
		return
	}

	go text.HandleTextMessage(src, dst, body)
	w.Header().Set("Content-Type", "text/xml")
	fmt.Fprintf(w, "%s", twiml)
}
func twilioTextStatusPost(w http.ResponseWriter, r *http.Request) {
	message_sid := r.PostFormValue("MessageSid")
	message_status := r.PostFormValue("MessageStatus")
	log.Info().Str("sid", message_sid).Str("status", message_status).Msg("Updated message status")
	text.UpdateMessageStatus(message_sid, message_status)
	fmt.Fprintf(w, "")
}
