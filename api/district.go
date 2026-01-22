package api

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/Gleipnir-Technology/nidus-sync/db"
	"github.com/Gleipnir-Technology/nidus-sync/db/models"
	"github.com/Gleipnir-Technology/nidus-sync/platform"
	"github.com/Gleipnir-Technology/nidus-sync/userfile"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

func apiGetDistrict(w http.ResponseWriter, r *http.Request) {
	var latStr, lngStr string
	err := r.ParseForm()
	if err != nil {
		render.Render(w, r, errRender(fmt.Errorf("Failed to parse GET form: %w", err)))
		return
	} else {
		latStr = r.FormValue("lat")
		lngStr = r.FormValue("lng")
	}
	lat, err := strconv.ParseFloat(latStr, 64)
	if err != nil {
		render.Render(w, r, errRender(fmt.Errorf("Failed to parse lat as float: %w", err)))
		return
	}
	lng, err := strconv.ParseFloat(lngStr, 64)
	if err != nil {
		render.Render(w, r, errRender(fmt.Errorf("Failed to parse lng as float: %w", err)))
		return
	}
	district, _, err := platform.DistrictForLocation(r.Context(), lng, lat)
	if err != nil {
		render.Render(w, r, errRender(fmt.Errorf("Failed to get district: %w", err)))
		return
	}
	if district == nil {
		http.NotFound(w, r)
		return
	}
	d := ResponseDistrict{
		Agency:  district.Agency.GetOr(""),
		Manager: district.GeneralMG.GetOr(""),
		Phone:   district.Phone1.GetOr(""),
		Website: district.Website.GetOr(""),
	}
	if err := render.Render(w, r, d); err != nil {
		render.Render(w, r, errRender(err))
	}
}

func apiGetDistrictLogo(w http.ResponseWriter, r *http.Request) {
	id_str := chi.URLParam(r, "id")
	org_id, err := strconv.ParseInt(id_str, 10, 32)
	if err != nil {
		render.Render(w, r, errRender(fmt.Errorf("%s is not a recognized organization ID: %w", id_str, err)))
		return
	}
	ctx := r.Context()
	org, err := models.FindOrganization(ctx, db.PGInstance.BobDB, int32(org_id))
	if err != nil {
		http.Error(w, "Organization not found", http.StatusNotFound)
		return
	}
	if org.LogoUUID.IsNull() {
		http.Error(w, "Logo not found", http.StatusNotFound)
		return
	}
	userfile.ImageFileContentWriteLogo(w, org.LogoUUID.MustGet())
}
