package api

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"

	"github.com/Gleipnir-Technology/nidus-sync/auth"
)

func AddRoutes(r chi.Router) {
	// Authenticated endpoints
	r.Use(render.SetContentType(render.ContentTypeJSON))
	r.Method("GET", "/mosquito-source", auth.NewEnsureAuth(apiMosquitoSource))
	r.Method("GET", "/service-request", auth.NewEnsureAuth(apiServiceRequest))
	r.Method("GET", "/trap-data", auth.NewEnsureAuth(apiTrapData))
	r.Method("GET", "/client/ios", auth.NewEnsureAuth(handleClientIos))
	r.Method("POST", "/audio/{uuid}", auth.NewEnsureAuth(apiAudioPost))
	r.Method("POST", "/audio/{uuid}/content", auth.NewEnsureAuth(apiAudioContentPost))
	r.Method("POST", "/image/{uuid}", auth.NewEnsureAuth(apiImagePost))
	r.Method("POST", "/image/{uuid}/content", auth.NewEnsureAuth(apiImageContentPost))

	// Unauthenticated endpoints
	r.Get("/district", apiGetDistrict)
	r.Get("/district/{slug}/logo", apiGetDistrictLogo)
	r.Post("/signin", postSignin)
	r.Post("/twilio/call", twilioCallPost)
	r.Post("/twilio/call/status", twilioCallStatusPost)
	r.Post("/twilio/message", twilioMessagePost)
	r.Post("/twilio/text", twilioTextPost)
	r.Post("/twilio/text/status", twilioTextStatusPost)
	r.Get("/webhook/fieldseeker", webhookFieldseeker)
	r.Post("/webhook/fieldseeker", webhookFieldseeker)
}
