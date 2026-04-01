package api

import (
	"github.com/Gleipnir-Technology/nidus-sync/auth"
	"github.com/Gleipnir-Technology/nidus-sync/platform/file"
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
	r.Handle("/communication", authenticatedHandlerJSON(listCommunication)).Methods("GET")
	r.Handle("/configuration/integration/arcgis", authenticatedHandlerJSONPost(postConfigurationIntegrationArcgis)).Methods("POST")
	r.Handle("/events", auth.NewEnsureAuth(streamEvents)).Methods("GET")
	r.Handle("/image/{uuid}", auth.NewEnsureAuth(apiImagePost)).Methods("POST")
	r.Handle("/image/{uuid}/content", auth.NewEnsureAuth(apiImageContentGet)).Methods("GET")
	r.Handle("/image/{uuid}/content", auth.NewEnsureAuth(apiImageContentPost)).Methods("POST")
	r.Handle("/leads", authenticatedHandlerJSON(listLead)).Methods("GET")
	r.Handle("/leads", authenticatedHandlerJSONPost(postLeads)).Methods("POST")
	r.Handle("/mosquito-source", auth.NewEnsureAuth(apiMosquitoSource)).Methods("GET")
	r.Handle("/publicreport/invalid", authenticatedHandlerJSONPost(postPublicreportInvalid)).Methods("POST")
	r.Handle("/publicreport/signal", authenticatedHandlerJSONPost(postPublicreportSignal)).Methods("POST")
	r.Handle("/publicreport/message", authenticatedHandlerJSONPost(postPublicreportMessage)).Methods("POST")
	r.Handle("/review/pool", authenticatedHandlerJSONPost(postReviewPool)).Methods("POST")
	r.Handle("/review-task", authenticatedHandlerJSON(listReviewTask)).Methods("GET")
	r.Handle("/service-request", auth.NewEnsureAuth(apiServiceRequest)).Methods("GET")
	r.Handle("/signal", authenticatedHandlerJSON(listSignal)).Methods("GET")
	r.Handle("/sudo/email", authenticatedHandlerJSONPost(postSudoEmail)).Methods("POST")
	r.Handle("/sudo/sms", authenticatedHandlerJSONPost(postSudoSMS)).Methods("POST")
	r.Handle("/sudo/sse", authenticatedHandlerJSONPost(postSudoSSE)).Methods("POST")
	r.Handle("/trap-data", auth.NewEnsureAuth(apiTrapData)).Methods("GET")
	r.Handle("/tile/{z}/{y}/{x}", auth.NewEnsureAuth(getTile)).Methods("GET")
	r.Handle("/upload/pool/flyover", authenticatedHandlerPostMultipart(postUploadPoolFlyoverCreate, file.CollectionCSV)).Methods("POST")
	r.Handle("/upload/pool/custom", authenticatedHandlerPostMultipart(postUploadPoolCustomCreate, file.CollectionCSV)).Methods("POST")
	r.Handle("/upload", authenticatedHandlerJSON(getUploadList)).Methods("GET")
	r.Handle("/upload/{id}", authenticatedHandlerJSON(getUploadByID)).Methods("GET")
	r.Handle("/upload/{id}/commit", authenticatedHandlerJSONPost(postUploadCommit)).Methods("POST")
	r.Handle("/upload/{id}/discard", authenticatedHandlerJSONPost(postUploadDiscard)).Methods("POST")
	r.Handle("/user/self", authenticatedHandlerJSON(getUserSelf)).Methods("GET")
	r.Handle("/user/suggestion", authenticatedHandlerJSON(listUserSuggestion)).Methods("GET")
	r.Handle("/user", authenticatedHandlerJSON(listUser)).Methods("GET")
	r.Handle("/user/{id}", authenticatedHandlerJSONPut(userPut)).Methods("PUT")

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
