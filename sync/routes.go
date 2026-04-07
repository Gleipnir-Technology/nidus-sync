package sync

import (
	"github.com/Gleipnir-Technology/nidus-sync/static"
	"github.com/gorilla/mux"
)

func Router(r *mux.Router) {
	// Unauthenticated endpoints
	r.HandleFunc("/arcgis/oauth/begin", getArcgisOauthBegin)
	r.HandleFunc("/arcgis/oauth/callback", getArcgisOauthCallback)
	r.HandleFunc("/mailer/pool/random", getMailerPoolRandom)
	r.HandleFunc("/mailer/mode-1", getMailer1)
	r.HandleFunc("/mailer/mode-2", getMailer2)
	r.HandleFunc("/mailer/mode-3/{code}", getMailer3)
	r.HandleFunc("/mailer/mode-1/preview", getMailer1Preview)
	r.HandleFunc("/mailer/mode-2/preview", getMailer2Preview)
	r.HandleFunc("/mailer/mode-3/{code}/preview", getMailer3Preview)
	r.HandleFunc("/district", getDistrict)

	// Mock endpoints
	r.HandleFunc("/mock", renderMockList)
	addMock(r, "/mock/report", "sync/mock/report.html")
	addMock(r, "/mock/report/{code}", "sync/mock/report-detail.html")
	addMock(r, "/mock/report/{code}/confirm", "sync/mock/report-confirmation.html")
	addMock(r, "/mock/report/{code}/contribute", "sync/mock/report-contribute.html")
	addMock(r, "/mock/report/{code}/evidence", "sync/mock/report-evidence.html")
	addMock(r, "/mock/report/{code}/schedule", "sync/mock/report-schedule.html")
	addMock(r, "/mock/report/{code}/update", "sync/mock/report-update.html")

	// Utility endpoints
	r.HandleFunc("/privacy", getPrivacy)
	r.HandleFunc("/qr-code/marketing", getQRCodeMarketing)
	r.HandleFunc("/qr-code/report/{code}", getQRCodeReport)
	r.HandleFunc("/qr-code/mailer/{code}", getQRCodeMailer)
	r.HandleFunc("/template-test", getTemplateTest)

	//r.HandleFunc("/", getRoot)
	//r.HandleFunc("/_/*", getRoot)

	static.AddStaticRoute(r, "/static")
	r.PathPrefix("/").Handler(static.SinglePageApp("static/gen/sync"))
}
