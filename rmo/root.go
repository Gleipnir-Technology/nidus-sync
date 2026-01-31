package rmo

import (
	"fmt"
	"net/http"

	"github.com/Gleipnir-Technology/nidus-sync/config"
	"github.com/Gleipnir-Technology/nidus-sync/html"
	"github.com/rs/zerolog/log"
)

type ContentDistrict struct {
	Name       string
	URLLogo    string
	URLWebsite string
}
type ContentPrivacy struct {
	Address   string
	Company   string
	Site      string
	URLReport string
}
type ContentRoot struct {
	URL ContentURL
}
type ContentURL struct {
	Nuisance               string
	NuisanceSubmit         string
	NuisanceSubmitComplete string
	Status                 string
	Tegola                 string
	Water                  string
}

var (
	PrivacyT = buildTemplate("privacy", "base")
	RootT    = buildTemplate("root", "base")
	TermsT   = buildTemplate("terms", "base")
)

func boolFromForm(r *http.Request, k string) bool {
	s := r.PostFormValue(k)
	if s == "on" {
		return true
	}
	return false
}

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
		ContentRoot{
			URL: makeContentURL(),
		},
	)
}

func getRobots(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "User-agent: *\n")
	fmt.Fprint(w, "Allow: /\n")
}
func getTerms(w http.ResponseWriter, r *http.Request) {
	html.RenderOrError(
		w,
		TermsT,
		ContentRoot{
			URL: makeContentURL(),
		},
	)
}

func makeContentURL() ContentURL {
	return ContentURL{
		Nuisance:       makeURL("nuisance"),
		NuisanceSubmit: makeURL("nuisance"),
		Status:         makeURL("status"),
		Tegola:         config.MakeURLTegola("/"),
		Water:          makeURL("water"),
	}
}

func makeURL(p string) string {
	return config.MakeURLReport("/%s", p)
}

// Respond with an error that is visible to the user
func respondError(w http.ResponseWriter, m string, e error, s int) {
	log.Warn().Int("status", s).Err(e).Str("user message", m).Msg("Responding with an error")
	http.Error(w, m, s)
}
