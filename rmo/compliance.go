package rmo

import (
	"net/http"

	"github.com/Gleipnir-Technology/nidus-sync/html"
)

func getDistrictCompliance(w http.ResponseWriter, r *http.Request) {
	district, err := districtBySlug(r)
	if err != nil {
		respondError(w, "Failed to lookup organization", err, http.StatusBadRequest)
		return
	}
	html.RenderOrError(
		w,
		"rmo/district-compliance.html",
		ContentNuisance{
			District: newContentDistrict(district),
			URL:      makeContentURL(nil),
		},
	)
}

func getDistrictComplianceAddress(w http.ResponseWriter, r *http.Request) {
	district, err := districtBySlug(r)
	if err != nil {
		respondError(w, "Failed to lookup organization", err, http.StatusBadRequest)
		return
	}
	html.RenderOrError(
		w,
		"rmo/district-compliance-address.html",
		ContentNuisance{
			District: newContentDistrict(district),
			URL:      makeContentURL(nil),
		},
	)
}

func getDistrictComplianceConcern(w http.ResponseWriter, r *http.Request) {
	district, err := districtBySlug(r)
	if err != nil {
		respondError(w, "Failed to lookup organization", err, http.StatusBadRequest)
		return
	}
	html.RenderOrError(
		w,
		"rmo/district-compliance-concern.html",
		ContentNuisance{
			District: newContentDistrict(district),
			URL:      makeContentURL(nil),
		},
	)
}

func getDistrictComplianceEvidence(w http.ResponseWriter, r *http.Request) {
	district, err := districtBySlug(r)
	if err != nil {
		respondError(w, "Failed to lookup organization", err, http.StatusBadRequest)
		return
	}
	html.RenderOrError(
		w,
		"rmo/district-compliance-evidence.html",
		ContentNuisance{
			District: newContentDistrict(district),
			URL:      makeContentURL(nil),
		},
	)
}
