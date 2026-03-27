package api

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"

	"github.com/Gleipnir-Technology/nidus-sync/auth"
	"github.com/Gleipnir-Technology/nidus-sync/platform/file"
)

func AddRoutes(r chi.Router) {
	r.Use(render.SetContentType(render.ContentTypeJSON))
	// Unauthenticated endpoints
	r.Post("/signin", handlerJSONPost(postSignin))
	r.Post("/signup", handlerJSONPost(postSignup))
	// Authenticated endpoints
	r.Method("POST", "/audio/{uuid}", auth.NewEnsureAuth(apiAudioPost))
	r.Method("POST", "/audio/{uuid}/content", auth.NewEnsureAuth(apiAudioContentPost))
	r.Method("GET", "/client/ios", auth.NewEnsureAuth(handleClientIos))
	r.Method("GET", "/communication", authenticatedHandlerJSON(listCommunication))
	r.Method("POST", "/configuration/integration/arcgis", authenticatedHandlerJSONPost(postConfigurationIntegrationArcgis))
	r.Method("GET", "/events", auth.NewEnsureAuth(streamEvents))
	r.Method("POST", "/image/{uuid}", auth.NewEnsureAuth(apiImagePost))
	r.Method("GET", "/image/{uuid}/content", auth.NewEnsureAuth(apiImageContentGet))
	r.Method("POST", "/image/{uuid}/content", auth.NewEnsureAuth(apiImageContentPost))
	r.Method("GET", "/leads", authenticatedHandlerJSON(listLead))
	r.Method("POST", "/leads", authenticatedHandlerJSONPost(postLeads))
	r.Method("GET", "/mosquito-source", auth.NewEnsureAuth(apiMosquitoSource))
	r.Method("POST", "/publicreport/invalid", authenticatedHandlerJSONPost(postPublicreportInvalid))
	r.Method("POST", "/publicreport/signal", authenticatedHandlerJSONPost(postPublicreportSignal))
	r.Method("POST", "/publicreport/message", authenticatedHandlerJSONPost(postPublicreportMessage))
	r.Method("POST", "/review/pool", authenticatedHandlerJSONPost(postReviewPool))
	r.Method("GET", "/review-task/pool", authenticatedHandlerJSON(listReviewTaskPool))
	r.Method("GET", "/service-request", auth.NewEnsureAuth(apiServiceRequest))
	r.Method("GET", "/signal", authenticatedHandlerJSON(listSignal))
	r.Method("POST", "/sudo/email", authenticatedHandlerJSONPost(postSudoEmail))
	r.Method("POST", "/sudo/sms", authenticatedHandlerJSONPost(postSudoSMS))
	r.Method("POST", "/sudo/sse", authenticatedHandlerJSONPost(postSudoSSE))
	r.Method("GET", "/trap-data", auth.NewEnsureAuth(apiTrapData))
	r.Method("GET", "/tile/{z}/{y}/{x}", auth.NewEnsureAuth(getTile))
	r.Method("POST", "/upload/pool/flyover", authenticatedHandlerPostMultipart(postUploadPoolFlyoverCreate, file.CollectionCSV))
	r.Method("POST", "/upload/pool/custom", authenticatedHandlerPostMultipart(postUploadPoolCustomCreate, file.CollectionCSV))
	r.Method("GET", "/upload/{id}", authenticatedHandlerJSON(getUploadByID))
	r.Method("POST", "/upload/{id}/commit", authenticatedHandlerJSONPost(postUploadCommit))
	r.Method("POST", "/upload/{id}/discard", authenticatedHandlerJSONPost(postUploadDiscard))
	r.Method("GET", "/user/self", authenticatedHandlerJSON(getUserSelf))
	r.Method("GET", "/user/suggestion", authenticatedHandlerJSON(listUserSuggestion))
	r.Method("GET", "/user", authenticatedHandlerJSON(listUser))

	// Unauthenticated endpoints
	r.Get("/district", apiGetDistrict)
	r.Get("/district/{slug}/logo", apiGetDistrictLogo)
	r.Get("/compliance-request/image/pool/{public_id}", getComplianceRequestImagePool)
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
