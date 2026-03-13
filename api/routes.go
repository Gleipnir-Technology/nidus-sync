package api

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"

	"github.com/Gleipnir-Technology/nidus-sync/auth"
)

func AddRoutes(r chi.Router) {
	// Authenticated endpoints
	r.Use(render.SetContentType(render.ContentTypeJSON))
	r.Method("POST", "/audio/{uuid}", auth.NewEnsureAuth(apiAudioPost))
	r.Method("POST", "/audio/{uuid}/content", auth.NewEnsureAuth(apiAudioContentPost))
	r.Method("GET", "/client/ios", auth.NewEnsureAuth(handleClientIos))
	r.Method("GET", "/communication", authenticatedHandlerJSON(listCommunication))
	r.Method("GET", "/events", auth.NewEnsureAuth(streamEvents))
	r.Method("POST", "/image/{uuid}", auth.NewEnsureAuth(apiImagePost))
	r.Method("GET", "/image/{uuid}/content", auth.NewEnsureAuth(apiImageContentGet))
	r.Method("POST", "/image/{uuid}/content", auth.NewEnsureAuth(apiImageContentPost))
	r.Method("GET", "/leads", authenticatedHandlerJSON(listLead))
	r.Method("POST", "/leads", authenticatedHandlerJSONPost(postLeads))
	r.Method("GET", "/mosquito-source", auth.NewEnsureAuth(apiMosquitoSource))
	r.Method("POST", "/review/pool", authenticatedHandlerJSONPost(postReviewPool))
	r.Method("GET", "/review-task/pool", authenticatedHandlerJSON(listReviewTaskPool))
	r.Method("GET", "/service-request", auth.NewEnsureAuth(apiServiceRequest))
	r.Method("GET", "/signal", authenticatedHandlerJSON(listSignal))
	r.Method("GET", "/trap-data", auth.NewEnsureAuth(apiTrapData))
	r.Method("GET", "/tile/{z}/{y}/{x}", auth.NewEnsureAuth(getTile))
	r.Method("GET", "/user", authenticatedHandlerJSON(getUser))

	// Unauthenticated endpoints
	r.Get("/district", apiGetDistrict)
	r.Get("/district/{slug}/logo", apiGetDistrictLogo)
	r.Get("/compliance-request/image/pool/{public_id}", getComplianceRequestImagePool)
	r.Post("/signin", postSignin)
	r.Post("/twilio/call", twilioCallPost)
	r.Post("/twilio/call/status", twilioCallStatusPost)
	r.Post("/twilio/message", twilioMessagePost)
	r.Post("/twilio/text", twilioTextPost)
	r.Post("/twilio/text/status", twilioTextStatusPost)
	r.Get("/voipms/text", voipmsTextGet)
	r.Post("/voipms/text", voipmsTextPost)
	r.Get("/webhook/fieldseeker", webhookFieldseeker)
	r.Post("/webhook/fieldseeker", webhookFieldseeker)
}
