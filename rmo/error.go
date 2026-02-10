package rmo

import (
	"net/http"

	//"github.com/Gleipnir-Technology/nidus-sync/config"
	"github.com/Gleipnir-Technology/nidus-sync/html"
)

type ContentError struct {
	Code     string
	District *ContentDistrict
	URL      ContentURL
}

func getError(w http.ResponseWriter, r *http.Request) {
	code := r.FormValue("code")
	district, err := districtBySlug(r)
	if err != nil {
		//respondError(w, "Failed to lookup organization", err, http.StatusBadRequest)
		district = nil
	}
	html.RenderOrError(
		w,
		"rmo/error.html",
		ContentError{
			Code:     code,
			District: newContentDistrict(district),
			URL:      makeContentURL(nil),
		},
	)
}
