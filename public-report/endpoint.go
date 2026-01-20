package publicreport

import (
	"fmt"
	"net/http"

	"github.com/Gleipnir-Technology/nidus-sync/htmlpage"
	"github.com/rs/zerolog/log"
)

type ContextRegisterNotificationsComplete struct {
	ReportID string
}
type ContextRoot struct{}

var (
	PrivacyT = buildTemplate("privacy", "base")
	RootT    = buildTemplate("root", "base")
	TermsT   = buildTemplate("terms", "base")
)

func getPrivacy(w http.ResponseWriter, r *http.Request) {
	htmlpage.RenderOrError(
		w,
		PrivacyT,
		ContextRoot{},
	)
}
func getRoot(w http.ResponseWriter, r *http.Request) {
	htmlpage.RenderOrError(
		w,
		RootT,
		ContextRoot{},
	)
}

func getRobots(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "User-agent: *\n")
	fmt.Fprint(w, "Allow: /\n")
}
func getTerms(w http.ResponseWriter, r *http.Request) {
	htmlpage.RenderOrError(
		w,
		TermsT,
		ContextRoot{},
	)
}

// Respond with an error that is visible to the user
func respondError(w http.ResponseWriter, m string, e error, s int) {
	log.Warn().Int("status", s).Err(e).Str("user message", m).Msg("Responding with an error")
	http.Error(w, m, s)
}

func boolFromForm(r *http.Request, k string) bool {
	s := r.PostFormValue(k)
	if s == "on" {
		return true
	}
	return false
}

func postFormValueOrNone(r *http.Request, k string) string {
	v := r.PostFormValue(k)
	if v == "" {
		return "none"
	}
	return v
}
