package sync

import (
	"context"
	"net/http"

	//"github.com/Gleipnir-Technology/nidus-sync/h3utils"
	"github.com/Gleipnir-Technology/nidus-sync/html"
	nhttp "github.com/Gleipnir-Technology/nidus-sync/http"
	"github.com/Gleipnir-Technology/nidus-sync/platform"
	"github.com/go-chi/chi/v5"
	"github.com/uber/h3-go/v4"
)

type contentCell struct {
	BreedingSources []platform.BreedingSourceSummary
	CellBoundary    h3.CellBoundary
	Inspections     []platform.Inspection
	Traps           []platform.TrapSummary
	Treatments      []platform.Treatment
}

func getCellDetails(ctx context.Context, r *http.Request, user platform.User) (*html.Response[contentCell], *nhttp.ErrorWithStatus) {
	cell_str := chi.URLParam(r, "cell")
	if cell_str == "" {
		return nil, nhttp.NewErrorStatus(http.StatusBadRequest, "There should always be a cell")
	}
	c, err := HexToInt64(cell_str)
	if err != nil {
		return nil, nhttp.NewErrorStatus(http.StatusBadRequest, "Cannot convert provided cell to uint64")
	}
	boundary, err := h3.Cell(c).Boundary()
	if err != nil {
		return nil, nhttp.NewError("Failed to get boundary: %w", err)
	}
	inspections, err := platform.InspectionsByCell(ctx, user.Organization, h3.Cell(c))
	if err != nil {
		return nil, nhttp.NewError("Failed to get inspections by cell: %w", err)
	}
	/*
	center, err := h3.Cell(c).LatLng()
	if err != nil {
		return nil, nhttp.NewError("Failed to get center: %w", err)
	}
	geojson, err := h3utils.H3ToGeoJSON([]h3.Cell{h3.Cell(c)})
	if err != nil {
		return nil, nhttp.NewError("Failed to get boundaries: %w", err)
	}
	resolution := h3.Cell(c).Resolution()
	*/
	sources, err := platform.BreedingSourcesByCell(ctx, user.Organization, h3.Cell(c))
	if err != nil {
		return nil, nhttp.NewError("Failed to get sources: %w", err)
	}
	traps, err := platform.TrapsByCell(ctx, user.Organization, h3.Cell(c))
	if err != nil {
		return nil, nhttp.NewError("Failed to get traps: %w", err)
	}

	treatments, err := platform.TreatmentsByCell(ctx, user.Organization, h3.Cell(c))
	if err != nil {
		return nil, nhttp.NewError("Failed to get treatments: %w", err)
	}
	return html.NewResponse("sync/cell.html", contentCell{
		BreedingSources: sources,
		CellBoundary:    boundary,
		Inspections:     inspections,
		Traps:      traps,
		Treatments: treatments,
	}), nil
}
