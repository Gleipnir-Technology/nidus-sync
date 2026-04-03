package rmo

import (
	"net/http"

	"github.com/Gleipnir-Technology/nidus-sync/html"
)

type ContentWater struct {
	District *ContentDistrict
	URL      ContentURL
}

func getWater(w http.ResponseWriter, r *http.Request) {
	html.RenderOrError(
		w,
		"rmo/water.html",
		ContentWater{
			District: nil,
			URL:      makeContentURL(nil),
		},
	)
}
func getWaterDistrict(w http.ResponseWriter, r *http.Request) {
	district, err := districtBySlug(r)
	if err != nil {
		respondError(w, "Failed to lookup organization", err, http.StatusBadRequest)
		return
	}
	html.RenderOrError(
		w,
		"rmo/water.html",
		ContentWater{
			District: newContentDistrict(district),
			URL:      makeContentURL(district),
		},
	)
}
