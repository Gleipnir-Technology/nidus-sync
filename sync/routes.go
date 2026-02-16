package sync

import (
	"github.com/Gleipnir-Technology/nidus-sync/api"
	"github.com/Gleipnir-Technology/nidus-sync/auth"
	"github.com/Gleipnir-Technology/nidus-sync/html"
	"github.com/go-chi/chi/v5"
)

func Router() chi.Router {
	r := chi.NewRouter()
	// Root is a special endpoint that is neither authenticated nor unauthenticated
	r.Get("/", getRoot)

	// Unauthenticated endpoints
	r.Get("/arcgis/oauth/begin", getArcgisOauthBegin)
	r.Get("/arcgis/oauth/callback", getArcgisOauthCallback)
	r.Get("/district", getDistrict)

	// Mock endpoints
	r.Get("/mock", renderMock("mock-root.html"))
	r.Get("/mock/admin", renderMock("admin.html"))
	r.Get("/mock/admin/service-request", renderMock("admin-service-request.html"))
	r.Get("/mock/data-entry", renderMock("data-entry.html"))
	r.Get("/mock/data-entry/bad", renderMock("data-entry-bad.html"))
	r.Get("/mock/dispatch", renderMock("dispatch.html"))
	r.Get("/mock/dispatch-results", renderMock("dispatch-results.html"))
	r.Get("/mock/report", renderMock("report.html"))
	r.Get("/mock/report/{code}", renderMock("report-detail.html"))
	r.Get("/mock/report/{code}/confirm", renderMock("report-confirmation.html"))
	r.Get("/mock/report/{code}/contribute", renderMock("report-contribute.html"))
	r.Get("/mock/report/{code}/evidence", renderMock("report-evidence.html"))
	r.Get("/mock/report/{code}/schedule", renderMock("report-schedule.html"))
	r.Get("/mock/report/{code}/update", renderMock("report-update.html"))
	r.Get("/mock/service-request", renderMock("service-request.html"))
	r.Get("/mock/service-request/{code}", renderMock("service-request-detail.html"))
	r.Get("/mock/service-request-location", renderMock("service-request-location.html"))
	r.Get("/mock/service-request-mosquito", renderMock("service-request-mosquito.html"))
	r.Get("/mock/service-request-pool", renderMock("service-request-pool.html"))
	r.Get("/mock/service-request-quick", renderMock("service-request-quick.html"))
	r.Get("/mock/service-request-quick-confirmation", renderMock("service-request-quick-confirmation.html"))
	r.Get("/mock/service-request-updates", renderMock("service-request-updates.html"))
	r.Get("/mock/setting", renderMock("setting-mock.html"))
	r.Get("/mock/setting/integration", renderMock("setting-integration.html"))
	r.Get("/mock/setting/pesticide", renderMock("setting-pesticide.html"))
	r.Get("/mock/setting/pesticide/add", renderMock("setting-pesticide-add.html"))
	r.Get("/mock/setting/user", renderMock("setting-user.html"))
	r.Get("/mock/setting/user/add", renderMock("setting-user-add.html"))

	// Utility endpoints
	r.Get("/oauth/refresh", getOAuthRefresh)
	r.Get("/privacy", getPrivacy)
	r.Get("/qr-code/report/{code}", getQRCodeReport)
	r.Get("/signin", getSignin)
	r.Post("/signin", postSignin)
	r.Get("/signup", getSignup)
	r.Post("/signup", postSignup)
	r.Get("/template-test", getTemplateTest)

	// Authenticated endpoints
	r.Route("/api", api.AddRoutes)
	r.Method("GET", "/cell/{cell}", auth.NewEnsureAuth(getCellDetails))
	r.Method("GET", "/layout-test", auth.NewEnsureAuth(getLayoutTest))
	r.Method("GET", "/notification", auth.NewEnsureAuth(getNotificationList))
	r.Method("GET", "/pool", auth.NewEnsureAuth(getPoolList))
	r.Method("GET", "/pool/upload", auth.NewEnsureAuth(getPoolUpload))
	r.Method("GET", "/pool/upload/{id}", auth.NewEnsureAuth(getPoolUploadByID))
	r.Method("POST", "/pool/upload", auth.NewEnsureAuth(postPoolUpload))
	r.Method("GET", "/setting", auth.NewEnsureAuth(getSetting))
	r.Method("GET", "/setting/district", auth.NewEnsureAuth(getSettingDistrict))
	r.Method("GET", "/setting/integration", auth.NewEnsureAuth(getSettingIntegration))
	r.Method("GET", "/signout", auth.NewEnsureAuth(getSignout))
	r.Method("GET", "/source/{globalid}", auth.NewEnsureAuth(getSource))
	r.Method("GET", "/stadia", auth.NewEnsureAuth(getStadia))
	r.Method("GET", "/trap/{globalid}", auth.NewEnsureAuth(getTrap))
	r.Method("GET", "/text/{destination}", auth.NewEnsureAuth(getTextMessages))

	html.AddStaticRoute(r, "/static")
	return r
}
