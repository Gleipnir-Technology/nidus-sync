package api

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/Gleipnir-Technology/nidus-sync/db"
	"github.com/Gleipnir-Technology/nidus-sync/db/models"
	"github.com/Gleipnir-Technology/nidus-sync/platform"
	"github.com/Gleipnir-Technology/nidus-sync/platform/file"
	"github.com/gorilla/mux"
)

func apiGetDistrict(w http.ResponseWriter, r *http.Request) {
	var latStr, lngStr string
	err := r.ParseForm()
	if err != nil {
		if err := renderShim(w, r, errRender(fmt.Errorf("Failed to parse GET form: %w", err))); err != nil {
			http.Error(w, fmt.Sprintf("render shim: %v", err), http.StatusInternalServerError)
		}
		return
	} else {
		latStr = r.FormValue("lat")
		lngStr = r.FormValue("lng")
	}
	lat, err := strconv.ParseFloat(latStr, 64)
	if err != nil {
		if err := renderShim(w, r, errRender(fmt.Errorf("Failed to parse lat as float: %w", err))); err != nil {
			http.Error(w, fmt.Sprintf("render shim: %v", err), http.StatusInternalServerError)
		}
		return
	}
	lng, err := strconv.ParseFloat(lngStr, 64)
	if err != nil {
		renderShim(w, r, errRender(fmt.Errorf("Failed to parse lng as float: %w", err)))
		return
	}
	org, err := platform.DistrictForLocation(r.Context(), lng, lat)
	if err != nil {
		renderShim(w, r, errRender(fmt.Errorf("Failed to get district: %w", err)))
		return
	}
	if org == nil {
		http.NotFound(w, r)
		return
	}
	d := ResponseDistrict{
		Agency:  org.Name,
		Manager: org.GeneralManagerName.GetOr(""),
		Phone:   org.OfficePhone.GetOr(""),
		Website: org.Website.GetOr(""),
	}
	if err := renderShim(w, r, d); err != nil {
		renderShim(w, r, errRender(err))
	}
}

func apiGetDistrictLogo(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	slug := vars["slug"]
	ctx := r.Context()
	rows, err := models.Organizations.Query(
		models.SelectWhere.Organizations.Slug.EQ(slug),
	).All(ctx, db.PGInstance.BobDB)
	if err != nil {
		http.Error(w, "Failed to query", http.StatusInternalServerError)
		return
	}
	switch len(rows) {
	case 0:
		http.Error(w, "Organization not found", http.StatusNotFound)
		return
	case 1:
		org := rows[0]
		if org.LogoUUID.IsNull() {
			http.Error(w, "Logo not found", http.StatusNotFound)
			return
		}
		file.ImageFileToWriter(file.CollectionLogo, org.LogoUUID.MustGet(), w)
		return
	default:
		http.Error(w, "Too many organizations, this is a programmer error", http.StatusInternalServerError)
		return
	}
}
