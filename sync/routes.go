package sync

import (
	"github.com/Gleipnir-Technology/nidus-sync/api"
	"github.com/Gleipnir-Technology/nidus-sync/auth"
	"github.com/Gleipnir-Technology/nidus-sync/htmlpage"
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

	r.Get("/mock", renderMock("mock-root"))
	r.Get("/mock/admin", renderMock("admin"))
	r.Get("/mock/admin/service-request", renderMock("admin-service-request"))
	r.Get("/mock/data-entry", renderMock("data-entry"))
	r.Get("/mock/data-entry/bad", renderMock("data-entry-bad"))
	r.Get("/mock/data-entry/good", renderMock("data-entry-good"))
	r.Get("/mock/dispatch", renderMock("dispatch"))
	r.Get("/mock/dispatch-results", renderMock("dispatch-results"))
	r.Get("/mock/report", renderMock("report"))
	r.Get("/mock/report/{code}", renderMock("report-detail"))
	r.Get("/mock/report/{code}/confirm", renderMock("report-confirmation"))
	r.Get("/mock/report/{code}/contribute", renderMock("report-contribute"))
	r.Get("/mock/report/{code}/evidence", renderMock("report-evidence"))
	r.Get("/mock/report/{code}/schedule", renderMock("report-schedule"))
	r.Get("/mock/report/{code}/update", renderMock("report-update"))
	r.Get("/mock/service-request", renderMock("service-request"))
	r.Get("/mock/service-request/{code}", renderMock("service-request-detail"))
	r.Get("/mock/service-request-location", renderMock("service-request-location"))
	r.Get("/mock/service-request-mosquito", renderMock("service-request-mosquito"))
	r.Get("/mock/service-request-pool", renderMock("service-request-pool"))
	r.Get("/mock/service-request-quick", renderMock("service-request-quick"))
	r.Get("/mock/service-request-quick-confirmation", renderMock("service-request-quick-confirmation"))
	r.Get("/mock/service-request-updates", renderMock("service-request-updates"))
	r.Get("/mock/setting", renderMock("setting-mock"))
	r.Get("/mock/setting/integration", renderMock("setting-integration"))
	r.Get("/mock/setting/pesticide", renderMock("setting-pesticide"))
	r.Get("/mock/setting/pesticide/add", renderMock("setting-pesticide-add"))
	r.Get("/mock/setting/user", renderMock("setting-user"))
	r.Get("/mock/setting/user/add", renderMock("setting-user-add"))

	r.Get("/oauth/refresh", getOAuthRefresh)

	r.Get("/privacy", getPrivacy)

	r.Get("/qr-code/report/{code}", getQRCodeReport)
	r.Get("/signin", getSignin)
	r.Post("/signin", postSignin)
	r.Get("/signup", getSignup)
	r.Post("/signup", postSignup)

	// Authenticated endpoints
	r.Route("/api", api.AddRoutes)
	r.Method("GET", "/cell/{cell}", auth.NewEnsureAuth(getCellDetails))
	r.Method("GET", "/settings", auth.NewEnsureAuth(getSettings))
	r.Method("GET", "/signout", auth.NewEnsureAuth(getSignout))
	r.Method("GET", "/source/{globalid}", auth.NewEnsureAuth(getSource))
	r.Method("GET", "/trap/{globalid}", auth.NewEnsureAuth(getTrap))

	htmlpage.AddStaticRoute(r, "/static")
	return r
}
