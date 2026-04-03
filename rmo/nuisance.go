package rmo

import (
	"fmt"
	"net/http"

	"github.com/Gleipnir-Technology/nidus-sync/html"
	"github.com/Gleipnir-Technology/nidus-sync/platform/report"
	//"github.com/rs/zerolog/log"
)

type ContentNuisance struct {
	District    *ContentDistrict
	MapboxToken string
	URL         ContentURL
}
type ContentNuisanceSubmitComplete struct {
	District *ContentDistrict
	ReportID string
	URL      ContentURL
}

func getNuisance(w http.ResponseWriter, r *http.Request) {
	html.RenderOrError(
		w,
		"rmo/nuisance.html",
		ContentNuisance{
			District: nil,
			URL:      makeContentURL(nil),
		},
	)
}
func getNuisanceDistrict(w http.ResponseWriter, r *http.Request) {
	district, err := districtBySlug(r)
	if err != nil {
		respondError(w, "Failed to lookup organization", err, http.StatusBadRequest)
		return
	}
	html.RenderOrError(
		w,
		"rmo/nuisance.html",
		ContentNuisance{
			District: newContentDistrict(district),
			URL:      makeContentURL(nil),
		},
	)
}
func getSubmitComplete(w http.ResponseWriter, r *http.Request) {
	report_id := r.URL.Query().Get("report")
	district, err := report.DistrictForReport(r.Context(), report_id)
	if err != nil {
		respondError(w, fmt.Sprintf("Failed to get district for report '%s'", report_id, err), err, http.StatusInternalServerError)
		return
	}
	html.RenderOrError(
		w,
		"rmo/submit-complete.html",
		ContentNuisanceSubmitComplete{
			District: newContentDistrict(district),
			ReportID: report_id,
			URL:      makeContentURL(nil),
		},
	)
}
