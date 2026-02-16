package sync

import (
	"net/http"

	"github.com/Gleipnir-Technology/nidus-sync/db"
	"github.com/Gleipnir-Technology/nidus-sync/db/models"
	"github.com/Gleipnir-Technology/nidus-sync/h3utils"
	"github.com/Gleipnir-Technology/nidus-sync/html"
	"github.com/go-chi/chi/v5"
	"github.com/uber/h3-go/v4"
)

type contentCell struct {
	BreedingSources []BreedingSourceSummary
	CellBoundary    h3.CellBoundary
	Inspections     []Inspection
	MapData         ComponentMap
	Traps           []TrapSummary
	Treatments      []Treatment
	URL             ContentURL
	User            User
}

func getCellDetails(w http.ResponseWriter, r *http.Request, user *models.User) {
	cell_str := chi.URLParam(r, "cell")
	if cell_str == "" {
		respondError(w, "There should always be a cell", nil, http.StatusBadRequest)
		return
	}
	c, err := HexToInt64(cell_str)
	if err != nil {
		respondError(w, "Cannot convert provided cell to uint64", err, http.StatusBadRequest)
		return
	}
	ctx := r.Context()
	org, err := user.Organization().One(ctx, db.PGInstance.BobDB)
	if err != nil {
		respondError(w, "Failed to get org", err, http.StatusInternalServerError)
		return
	}
	userContent, err := contentForUser(ctx, user)
	if err != nil {
		respondError(w, "Failed to get user", err, http.StatusInternalServerError)
		return
	}
	center, err := h3.Cell(c).LatLng()
	if err != nil {
		respondError(w, "Failed to get center", err, http.StatusInternalServerError)
		return
	}
	boundary, err := h3.Cell(c).Boundary()
	if err != nil {
		respondError(w, "Failed to get boundary", err, http.StatusInternalServerError)
		return
	}
	inspections, err := inspectionsByCell(ctx, org, h3.Cell(c))
	if err != nil {
		respondError(w, "Failed to get inspections by cell", err, http.StatusInternalServerError)
		return
	}
	geojson, err := h3utils.H3ToGeoJSON([]h3.Cell{h3.Cell(c)})
	if err != nil {
		respondError(w, "Failed to get boundaries", err, http.StatusInternalServerError)
		return
	}
	resolution := h3.Cell(c).Resolution()
	sources, err := breedingSourcesByCell(ctx, org, h3.Cell(c))
	if err != nil {
		respondError(w, "Failed to get sources", err, http.StatusInternalServerError)
		return
	}
	traps, err := trapsByCell(ctx, org, h3.Cell(c))
	if err != nil {
		respondError(w, "Failed to get traps", err, http.StatusInternalServerError)
		return
	}

	treatments, err := treatmentsByCell(ctx, org, h3.Cell(c))
	if err != nil {
		respondError(w, "Failed to get treatments", err, http.StatusInternalServerError)
		return
	}
	data := contentCell{
		BreedingSources: sources,
		CellBoundary:    boundary,
		Inspections:     inspections,
		MapData: ComponentMap{
			Center: h3.LatLng{
				Lat: center.Lat,
				Lng: center.Lng,
			},
			GeoJSON: geojson,
			Zoom:    resolution + 5,
		},
		Traps:      traps,
		Treatments: treatments,
		URL:        newContentURL(),
		User:       userContent,
	}
	html.RenderOrError(w, "sync/cell.html", &data)
}
