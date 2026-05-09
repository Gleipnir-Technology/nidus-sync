package api

import (
	"net/http"

	"github.com/Gleipnir-Technology/nidus-sync/db"
	"github.com/Gleipnir-Technology/nidus-sync/db/models"
	"github.com/Gleipnir-Technology/nidus-sync/platform/file"
	"github.com/gorilla/mux"
)

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
