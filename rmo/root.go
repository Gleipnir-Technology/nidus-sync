package rmo

import (
	"fmt"
	"net/http"

	"github.com/Gleipnir-Technology/nidus-sync/config"
	"github.com/Gleipnir-Technology/nidus-sync/db/models"
	"github.com/Gleipnir-Technology/nidus-sync/html"
	"github.com/rs/zerolog/log"
)

type ContentPrivacy struct {
	Address   string
	Company   string
	Site      string
	URLReport string
}
type ContentRoot struct {
	District *ContentDistrict
	URL      ContentURL
}
type ContentURL struct {
	Nuisance       string
	NuisanceSubmit string
	SubmitComplete string
	Status         string
	Tegola         string
	Water          string
	WaterSubmit    string
}

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
		"rmo/privacy.html",
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
		"rmo/root.html",
		ContentRoot{
			URL: makeContentURL(nil),
		},
	)
}
func getRootDistrict(w http.ResponseWriter, r *http.Request) {
	district, err := districtBySlug(r)
	if err != nil {
		respondError(w, "Failed to lookup organization", err, http.StatusBadRequest)
		return
	}
	html.RenderOrError(
		w,
		"rmo/root.html",
		ContentRoot{
			District: newContentDistrict(district),
			URL:      makeContentURL(district),
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
		"rmo/terms.html",
		ContentRoot{
			URL: makeContentURL(nil),
		},
	)
}

func makeContentURL(district *models.Organization) ContentURL {
	if district == nil || district.Slug.IsNull() {
		return ContentURL{
			Nuisance:       makeURL("/nuisance"),
			NuisanceSubmit: makeURL("/nuisance"),
			Status:         makeURL("/status"),
			Tegola:         config.MakeURLTegola("/"),
			Water:          makeURL("/water"),
			WaterSubmit:    makeURL("/water"),
		}
	} else {
		slug := district.Slug.MustGet()
		return ContentURL{
			Nuisance:       makeURL("/district/%s/nuisance", slug),
			NuisanceSubmit: makeURL("/nuisance", slug),
			Status:         makeURL("/status"),
			Tegola:         config.MakeURLTegola("/"),
			Water:          makeURL("/district/%s/water", slug),
			WaterSubmit:    makeURL("/water"),
		}
	}
}

func makeURL(f string, args ...string) string {
	return config.MakeURLReport(f, args...)
}

// Respond with an error that is visible to the user
func respondError(w http.ResponseWriter, m string, e error, s int) {
	log.Warn().Int("status", s).Err(e).Str("user message", m).Msg("Responding with an error")
	http.Error(w, m, s)
}
