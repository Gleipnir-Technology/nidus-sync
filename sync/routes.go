package sync

import (
	"github.com/Gleipnir-Technology/nidus-sync/api"
	"github.com/Gleipnir-Technology/nidus-sync/static"
	"github.com/go-chi/chi/v5"
)

func Router() chi.Router {
	r := chi.NewRouter()

	// Unauthenticated endpoints
	r.Get("/arcgis/oauth/begin", getArcgisOauthBegin)
	r.Get("/arcgis/oauth/callback", getArcgisOauthCallback)
	r.Get("/mailer/pool/random", getMailerPoolRandom)
	r.Get("/mailer/mode-1", getMailer1)
	r.Get("/mailer/mode-2/{code}", getMailer2)
	r.Get("/mailer/mode-3/{code}", getMailer3)
	r.Get("/mailer/mode-1/preview", getMailer1Preview)
	r.Get("/mailer/mode-2/preview", getMailer2Preview)
	r.Get("/mailer/mode-3/{code}/preview", getMailer3Preview)
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
	r.Get("/template-test", getTemplateTest)

	// Authenticated endpoints
	r.Route("/api", api.AddRoutes)

	r.Get("/", getRoot)
	r.Get("/_/*", getRoot)

	static.AddStaticRoute(r, "/static")
	return r
}
