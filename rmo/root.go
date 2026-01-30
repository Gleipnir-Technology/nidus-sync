package rmo

import (
	"fmt"
	"net/http"

	"github.com/Gleipnir-Technology/nidus-sync/config"
	"github.com/Gleipnir-Technology/nidus-sync/html"
	"github.com/rs/zerolog/log"
)

type ContentPrivacy struct {
	Address   string
	Company   string
	Site      string
	URLReport string
}
type ContentRoot struct{}

var (
	PrivacyT = buildTemplate("privacy", "base")
	RootT    = buildTemplate("root", "base")
	TermsT   = buildTemplate("terms", "base")
)

func getPrivacy(w http.ResponseWriter, r *http.Request) {
	html.RenderOrError(
		w,
		PrivacyT,
		ContentPrivacy{
			Address:   "2726 S Quinn Ave, Gilbert, AZ, USA",
			Company:   "Gleipnir LLC",
			Site:      "Report Mosquitoes Online",
			URLReport: config.MakeURLReport("/"),
		},
	)
}
func getRoot(w http.ResponseWriter, r *http.Request) {
	html.RenderOrError(
		w,
		RootT,
		ContentRoot{},
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
		ContentRoot{},
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
