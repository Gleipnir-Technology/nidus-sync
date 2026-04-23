package sync

import (
	"github.com/Gleipnir-Technology/nidus-sync/static"
	"github.com/gorilla/mux"
)

func Router(r *mux.Router) {
	// Unauthenticated endpoints
	r.HandleFunc("/oauth/arcgis/begin", getArcgisOauthBegin)
	r.HandleFunc("/oauth/arcgis/callback", getArcgisOauthCallback)
	r.HandleFunc("/mailer/pool/random", getMailerPoolRandom)
	r.HandleFunc("/mailer/mode-1", getMailer1)
	r.HandleFunc("/mailer/mode-2", getMailer2)
	r.HandleFunc("/mailer/mode-3/{code}", getMailer3)
	r.HandleFunc("/mailer/mode-1/preview", getMailer1Preview)
	r.HandleFunc("/mailer/mode-2/preview", getMailer2Preview)
	r.HandleFunc("/mailer/mode-3/{code}/preview", getMailer3Preview)

	// Utility endpoints
	r.HandleFunc("/privacy", getPrivacy)
	r.HandleFunc("/template-test", getTemplateTest)

	//r.HandleFunc("/", getRoot)
	//r.HandleFunc("/_/*", getRoot)

	static.AddStaticRoute(r, "/static")
	r.PathPrefix("/").Handler(static.SinglePageApp("static/gen/sync")).Methods("GET")
}
