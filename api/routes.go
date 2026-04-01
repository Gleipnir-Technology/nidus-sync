package api

import (
	"github.com/Gleipnir-Technology/nidus-sync/auth"
	"github.com/Gleipnir-Technology/nidus-sync/platform/file"
	"github.com/Gleipnir-Technology/nidus-sync/resource"
	"github.com/gorilla/mux"
)

func AddRoutes(r *mux.Router) {
	//r.Use(render.SetContentType(render.ContentTypeJSON))
	// Unauthenticated endpoints
	r.HandleFunc("/signin", handlerJSONPost(postSignin))
	r.HandleFunc("/signup", handlerJSONPost(postSignup))
	// Authenticated endpoints
	r.Handle("/audio/{uuid}", auth.NewEnsureAuth(apiAudioPost)).Methods("POST")
	r.Handle("/audio/{uuid}/content", auth.NewEnsureAuth(apiAudioContentPost)).Methods("POST")
	r.Handle("/avatar", authenticatedHandlerPostMultipart(avatarPost, file.CollectionAvatar)).Methods("POST")
	r.Handle("/client/ios", auth.NewEnsureAuth(handleClientIos)).Methods("GET")
	communication := resource.Communication(r)
	r.Handle("/communication", authenticatedHandlerJSON(communication.List)).Methods("GET")
	r.Handle("/configuration/integration/arcgis", authenticatedHandlerJSONPost(postConfigurationIntegrationArcgis)).Methods("POST")
	r.Handle("/events", auth.NewEnsureAuth(streamEvents)).Methods("GET")
	r.Handle("/image/{uuid}", auth.NewEnsureAuth(apiImagePost)).Methods("POST")
	r.Handle("/image/{uuid}/content", auth.NewEnsureAuth(apiImageContentGet)).Methods("GET")
	r.Handle("/image/{uuid}/content", auth.NewEnsureAuth(apiImageContentPost)).Methods("POST")
	lead := resource.Lead(r)
	r.Handle("/leads", authenticatedHandlerJSON(lead.List)).Methods("GET")
	r.Handle("/leads", authenticatedHandlerJSONPost(lead.Create)).Methods("POST")
	r.Handle("/mosquito-source", auth.NewEnsureAuth(apiMosquitoSource)).Methods("GET")
	r.Handle("/publicreport/invalid", authenticatedHandlerJSONPost(postPublicreportInvalid)).Methods("POST")
	r.Handle("/publicreport/signal", authenticatedHandlerJSONPost(postPublicreportSignal)).Methods("POST")
	r.Handle("/publicreport/message", authenticatedHandlerJSONPost(postPublicreportMessage)).Methods("POST")
	r.Handle("/review/pool", authenticatedHandlerJSONPost(postReviewPool)).Methods("POST")
	review_task := resource.ReviewTask(r)
	r.Handle("/review-task", authenticatedHandlerJSON(review_task.List)).Methods("GET")
	r.Handle("/service-request", auth.NewEnsureAuth(apiServiceRequest)).Methods("GET")
	signal := resource.Signal(r)
	r.Handle("/signal", authenticatedHandlerJSON(signal.List)).Methods("GET")
	r.Handle("/sudo/email", authenticatedHandlerJSONPost(postSudoEmail)).Methods("POST")
	r.Handle("/sudo/sms", authenticatedHandlerJSONPost(postSudoSMS)).Methods("POST")
	r.Handle("/sudo/sse", authenticatedHandlerJSONPost(postSudoSSE)).Methods("POST")
	r.Handle("/trap-data", auth.NewEnsureAuth(apiTrapData)).Methods("GET")
	r.Handle("/tile/{z}/{y}/{x}", auth.NewEnsureAuth(getTile)).Methods("GET")
	upload := resource.Upload(r)
	r.Handle("/upload/pool/flyover", authenticatedHandlerPostMultipart(upload.PoolFlyoverCreate, file.CollectionCSV)).Methods("POST")
	r.Handle("/upload/pool/custom", authenticatedHandlerPostMultipart(upload.PoolCustomCreate, file.CollectionCSV)).Methods("POST")
	r.Handle("/upload", authenticatedHandlerJSON(upload.List)).Methods("GET")
	r.Handle("/upload/{id}", authenticatedHandlerJSON(upload.ByIDGet)).Methods("GET")
	r.Handle("/upload/{id}/commit", authenticatedHandlerJSONPost(upload.Commit)).Methods("POST")
	r.Handle("/upload/{id}/discard", authenticatedHandlerJSONPost(upload.Discard)).Methods("POST")

	user := resource.User(r)
	r.Handle("/user/self", authenticatedHandlerJSON(user.SelfGet)).Methods("GET")
	r.Handle("/user/suggestion", authenticatedHandlerJSON(user.SuggestionGet)).Methods("GET")
	r.Handle("/user", authenticatedHandlerJSONSlice(user.List)).Methods("GET")
	r.Handle("/user/{id}", authenticatedHandlerJSON(user.ByIDGet)).Methods("GET").Name("user.ByIDGet")
	r.Handle("/user/{id}", authenticatedHandlerJSONPut(user.ByIDPut)).Methods("PUT")

	// Unauthenticated endpoints
	r.HandleFunc("/district", apiGetDistrict).Methods("GET")
	r.HandleFunc("/district/{slug}/logo", apiGetDistrictLogo).Methods("GET")
	r.HandleFunc("/compliance-request/image/pool/{public_id}", getComplianceRequestImagePool).Methods("GET")
	r.HandleFunc("/twilio/call", twilioCallPost).Methods("POST")
	r.HandleFunc("/twilio/call/status", twilioCallStatusPost).Methods("POST")
	r.HandleFunc("/twilio/message", twilioMessagePost).Methods("POST")
	r.HandleFunc("/twilio/text", twilioTextPost).Methods("POST")
	r.HandleFunc("/twilio/text/status", twilioTextStatusPost).Methods("POST")
	r.HandleFunc("/voipms/text", voipmsTextGet).Methods("GET")
	r.HandleFunc("/voipms/text", voipmsTextPost).Methods("POST")
	r.HandleFunc("/webhook/fieldseeker", webhookFieldseeker).Methods("GET")
	r.HandleFunc("/webhook/fieldseeker", webhookFieldseeker).Methods("POST")
}
