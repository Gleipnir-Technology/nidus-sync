package rmo

import (
	"net/http"

	"github.com/Gleipnir-Technology/nidus-sync/html"
)

type ContentCompliance struct {
	District            *ContentDistrict
	HasCompleteResponse bool
	HasUsefulInfo       bool
	ReferenceNumber     string
	URL                 ContentURL
}

func getDistrictCompliance(w http.ResponseWriter, r *http.Request) {
	district, err := districtBySlug(r)
	if err != nil {
		respondError(w, "Failed to lookup organization", err, http.StatusBadRequest)
		return
	}
	html.RenderOrError(
		w,
		"rmo/district-compliance.html",
		ContentCompliance{
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
		ContentCompliance{
			District: newContentDistrict(district),
			URL:      makeContentURL(nil),
		},
	)
}

func getDistrictComplianceComplete(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	complete := query.Get("complete")
	is_complete := complete != ""
	useful := query.Get("useful")
	is_useful := useful != ""

	district, err := districtBySlug(r)
	if err != nil {
		respondError(w, "Failed to lookup organization", err, http.StatusBadRequest)
		return
	}
	html.RenderOrError(
		w,
		"rmo/district-compliance-complete.html",
		ContentCompliance{
			District:            newContentDistrict(district),
			HasCompleteResponse: is_complete,
			HasUsefulInfo:       is_useful,
			ReferenceNumber:     "ABC-123",
			URL:                 makeContentURL(nil),
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
		ContentCompliance{
			District: newContentDistrict(district),
			URL:      makeContentURL(nil),
		},
	)
}

func getDistrictComplianceContact(w http.ResponseWriter, r *http.Request) {
	district, err := districtBySlug(r)
	if err != nil {
		respondError(w, "Failed to lookup organization", err, http.StatusBadRequest)
		return
	}
	html.RenderOrError(
		w,
		"rmo/district-compliance-contact.html",
		ContentCompliance{
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
		ContentCompliance{
			District: newContentDistrict(district),
			URL:      makeContentURL(nil),
		},
	)
}

func getDistrictCompliancePermission(w http.ResponseWriter, r *http.Request) {
	district, err := districtBySlug(r)
	if err != nil {
		respondError(w, "Failed to lookup organization", err, http.StatusBadRequest)
		return
	}
	html.RenderOrError(
		w,
		"rmo/district-compliance-permission.html",
		ContentCompliance{
			District: newContentDistrict(district),
			URL:      makeContentURL(nil),
		},
	)
}
func getDistrictComplianceProcess(w http.ResponseWriter, r *http.Request) {
	district, err := districtBySlug(r)
	if err != nil {
		respondError(w, "Failed to lookup organization", err, http.StatusBadRequest)
		return
	}
	html.RenderOrError(
		w,
		"rmo/district-compliance-process.html",
		ContentCompliance{
			District: newContentDistrict(district),
			URL:      makeContentURL(nil),
		},
	)
}

func getDistrictComplianceSubmit(w http.ResponseWriter, r *http.Request) {
	district, err := districtBySlug(r)
	if err != nil {
		respondError(w, "Failed to lookup organization", err, http.StatusBadRequest)
		return
	}
	html.RenderOrError(
		w,
		"rmo/district-compliance-submit.html",
		ContentCompliance{
			District: newContentDistrict(district),
			URL:      makeContentURL(nil),
		},
	)
}
