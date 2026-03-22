package sync

import (
	"github.com/Gleipnir-Technology/nidus-sync/api"
	"github.com/Gleipnir-Technology/nidus-sync/auth"
	"github.com/Gleipnir-Technology/nidus-sync/static"
	"github.com/go-chi/chi/v5"
)

func Router() chi.Router {
	r := chi.NewRouter()

	// Unauthenticated endpoints
	r.Get("/arcgis/oauth/begin", getArcgisOauthBegin)
	r.Get("/arcgis/oauth/callback", getArcgisOauthCallback)
	r.Get("/mailer/{code}", getMailer)
	r.Get("/mailer/{code}/preview", getMailerPreview)
	r.Get("/district", getDistrict)

	// Mock endpoints
	r.Get("/mock", renderMockList)
	addMock(r, "/mock/report", "sync/mock/report.html")
	addMock(r, "/mock/report/{code}", "sync/mock/report-detail.html")
	addMock(r, "/mock/report/{code}/confirm", "sync/mock/report-confirmation.html")
	addMock(r, "/mock/report/{code}/contribute", "sync/mock/report-contribute.html")
	addMock(r, "/mock/report/{code}/evidence", "sync/mock/report-evidence.html")
	addMock(r, "/mock/report/{code}/schedule", "sync/mock/report-schedule.html")
	addMock(r, "/mock/report/{code}/update", "sync/mock/report-update.html")

	// Utility endpoints
	r.Get("/privacy", getPrivacy)
	r.Get("/qr-code/report/{code}", getQRCodeReport)
	r.Get("/qr-code/mailer/{code}", getQRCodeMailer)
	r.Get("/signin", getSignin)
	r.Post("/signin", postSignin)
	r.Get("/signup", getSignup)
	r.Post("/signup", postSignup)
	r.Get("/template-test", getTemplateTest)

	// Authenticated endpoints
	r.Route("/api", api.AddRoutes)

	r.Method("GET", "/", authenticatedHandler(getRoot))
	r.Method("GET", "/communication", authenticatedHandler(getRoot))
	r.Method("GET", "/configuration", authenticatedHandler(getRoot))
	r.Method("GET", "/configuration/user", authenticatedHandler(getRoot))
	r.Method("GET", "/configuration/user/add", authenticatedHandler(getRoot))
	r.Method("GET", "/intelligence", authenticatedHandler(getRoot))
	r.Method("GET", "/operations", authenticatedHandler(getRoot))
	r.Method("GET", "/planning", authenticatedHandler(getRoot))
	r.Method("GET", "/review", authenticatedHandler(getRoot))
	r.Method("GET", "/review/pool", authenticatedHandler(getRoot))
	r.Method("GET", "/review/site", authenticatedHandler(getRoot))

	r.Method("GET", "/admin", authenticatedHandler(getAdminDash))
	r.Method("GET", "/cell/{cell}", authenticatedHandler(getCellDetails))
	r.Method("GET", "/configuration/integration", authenticatedHandler(getConfigurationIntegration))
	r.Method("GET", "/configuration/integration/arcgis", authenticatedHandler(getConfigurationIntegrationArcgis))
	r.Method("POST", "/configuration/integration/arcgis", authenticatedHandlerPost(postConfigurationIntegrationArcgis))
	r.Method("GET", "/configuration/organization", authenticatedHandler(getConfigurationOrganization))
	r.Method("GET", "/configuration/pesticide", authenticatedHandler(getConfigurationPesticide))
	r.Method("GET", "/configuration/pesticide/add", authenticatedHandler(getConfigurationPesticideAdd))
	r.Method("GET", "/configuration/upload", authenticatedHandler(getUploadList))
	r.Method("GET", "/configuration/upload/pool", authenticatedHandler(getUploadPool))
	r.Method("GET", "/configuration/upload/pool/flyover", authenticatedHandler(getUploadPoolFlyoverCreate))
	r.Method("POST", "/configuration/upload/pool/flyover", authenticatedHandlerPostMultipart(postUploadPoolFlyoverCreate))
	r.Method("GET", "/configuration/upload/pool/custom", authenticatedHandler(getUploadPoolCustomCreate))
	r.Method("POST", "/configuration/upload/pool/custom", authenticatedHandlerPostMultipart(postUploadPoolCustomCreate))
	r.Method("GET", "/configuration/upload/{id}", authenticatedHandler(getUploadByID))
	r.Method("POST", "/configuration/upload/{id}/commit", authenticatedHandlerPost(postUploadCommit))
	r.Method("POST", "/configuration/upload/{id}/discard", authenticatedHandlerPost(postUploadDiscard))
	r.Method("GET", "/download", authenticatedHandler(getDownloadList))
	r.Method("GET", "/layout-test", authenticatedHandler(getLayoutTest))
	r.Method("GET", "/message", authenticatedHandler(getMessageList))
	r.Method("GET", "/notification", authenticatedHandler(getNotificationList))
	r.Method("GET", "/oauth/refresh", authenticatedHandler(getOAuthRefresh))
	r.Method("GET", "/parcel", authenticatedHandler(getParcel))
	r.Method("GET", "/pool", authenticatedHandler(getPoolList))
	r.Method("GET", "/pool/create", authenticatedHandler(getPoolCreate))
	r.Method("GET", "/pool/{id}", authenticatedHandler(getPoolByID))
	r.Method("GET", "/radar", authenticatedHandler(getRadar))
	r.Method("GET", "/service-request", authenticatedHandler(getServiceRequestList))
	r.Method("GET", "/service-request/{id}", authenticatedHandler(getServiceRequestDetail))
	r.Method("GET", "/signout", auth.NewEnsureAuth(getSignout))
	r.Method("GET", "/source/{globalid}", authenticatedHandler(getSource))
	r.Method("GET", "/stadia", authenticatedHandler(getStadia))
	r.Method("GET", "/sudo", authenticatedHandler(getSudo))
	r.Method("POST", "/sudo/email", authenticatedHandlerPost(postSudoEmail))
	r.Method("POST", "/sudo/sms", authenticatedHandlerPost(postSudoSMS))
	r.Method("POST", "/sudo/sse", authenticatedHandlerPost(postSudoSSE))
	r.Method("GET", "/trap/{globalid}", authenticatedHandler(getTrap))
	r.Method("GET", "/text/{destination}", authenticatedHandler(getTextMessages))
	r.Method("GET", "/tile/gps", auth.NewEnsureAuth(getTileGPS))

	static.AddStaticRoute(r, "/static")
	return r
}
