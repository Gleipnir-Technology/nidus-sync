package sync

import (
	"context"
	"net/http"

	"github.com/Gleipnir-Technology/nidus-sync/db/models"
	"github.com/Gleipnir-Technology/nidus-sync/h3utils"
	"github.com/Gleipnir-Technology/nidus-sync/html"
	nhttp "github.com/Gleipnir-Technology/nidus-sync/http"
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
}

func getCellDetails(ctx context.Context, r *http.Request, org *models.Organization, user *models.User) (*html.Response[contentCell], *nhttp.ErrorWithStatus) {
	cell_str := chi.URLParam(r, "cell")
	if cell_str == "" {
		return nil, nhttp.NewErrorStatus(http.StatusBadRequest, "There should always be a cell")
	}
	c, err := HexToInt64(cell_str)
	if err != nil {
		return nil, nhttp.NewErrorStatus(http.StatusBadRequest, "Cannot convert provided cell to uint64")
	}
	center, err := h3.Cell(c).LatLng()
	if err != nil {
		return nil, nhttp.NewError("Failed to get center: %w", err)
	}
	boundary, err := h3.Cell(c).Boundary()
	if err != nil {
		return nil, nhttp.NewError("Failed to get boundary: %w", err)
	}
	inspections, err := inspectionsByCell(ctx, org, h3.Cell(c))
	if err != nil {
		return nil, nhttp.NewError("Failed to get inspections by cell: %w", err)
	}
	geojson, err := h3utils.H3ToGeoJSON([]h3.Cell{h3.Cell(c)})
	if err != nil {
		return nil, nhttp.NewError("Failed to get boundaries: %w", err)
	}
	resolution := h3.Cell(c).Resolution()
	sources, err := breedingSourcesByCell(ctx, org, h3.Cell(c))
	if err != nil {
		return nil, nhttp.NewError("Failed to get sources: %w", err)
	}
	traps, err := trapsByCell(ctx, org, h3.Cell(c))
	if err != nil {
		return nil, nhttp.NewError("Failed to get traps: %w", err)
	}

	treatments, err := treatmentsByCell(ctx, org, h3.Cell(c))
	if err != nil {
		return nil, nhttp.NewError("Failed to get treatments: %w", err)
	}
	return html.NewResponse("sync/cell.html", contentCell{
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
	}), nil
}
