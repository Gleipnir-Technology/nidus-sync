package platform

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/Gleipnir-Technology/bob"
	"github.com/Gleipnir-Technology/bob/dialect/psql"
	"github.com/Gleipnir-Technology/bob/dialect/psql/sm"
	"github.com/Gleipnir-Technology/nidus-sync/db"
	"github.com/Gleipnir-Technology/nidus-sync/db/models"
	"github.com/Gleipnir-Technology/nidus-sync/db/sql"
	"github.com/google/uuid"
	"github.com/uber/h3-go/v4"
)

type Inspection struct {
	Action     string
	Date       *time.Time
	Notes      string
	Location   string
	LocationID uuid.UUID
}

func BreedingSourcesByCell(ctx context.Context, org Organization, c h3.Cell) ([]BreedingSourceSummary, error) {
	boundary, err := c.Boundary()
	if err != nil {
		return nil, fmt.Errorf("Failed to get cell boundary: %w", err)
	}
	geom_query := gisStatement(boundary)
	rows, err := org.model.Pointlocations(
		sm.Where(
			psql.F("ST_Within", "geospatial", geom_query),
		),
		sm.OrderBy("lasttreatdate"),
	).All(ctx, db.PGInstance.BobDB)
	if err != nil {
		return nil, fmt.Errorf("Failed to query rows: %w", err)
	}
	return toBreedingSourceSummary(rows), nil
}
func SourceByGlobalID(ctx context.Context, org Organization, id uuid.UUID) (*BreedingSourceDetail, error) {
	row, err := org.model.Pointlocations(
		models.SelectWhere.FieldseekerPointlocations.Globalid.EQ(id),
	).One(ctx, db.PGInstance.BobDB)
	if err != nil {
		return nil, fmt.Errorf("Failed to get point location: %w", err)
	}
	return toBreedingSource(row)
}

func TrapsBySource(ctx context.Context, org Organization, sourceID uuid.UUID) ([]TrapNearby, error) {
	locations, err := sql.TrapLocationBySourceID(org.ID(), sourceID).All(ctx, db.PGInstance.BobDB)
	if err != nil {
		return nil, fmt.Errorf("Failed to query rows: %w", err)
	}

	location_ids := make([]uuid.UUID, 0)
	var args []bob.Expression
	for _, location := range locations {
		location_ids = append(location_ids, location.TrapLocationGlobalid)
		args = append(args, psql.Arg(location.TrapLocationGlobalid))
	}
	trap_data, err := sql.TrapDataByLocationIDRecent(org.ID(), location_ids).All(ctx, db.PGInstance.BobDB)
	if err != nil {
		return nil, fmt.Errorf("Failed to query trap data: %w", err)
	}

	counts, err := sql.TrapCountByLocationID(org.ID(), location_ids).All(ctx, db.PGInstance.BobDB)
	if err != nil {
		return nil, fmt.Errorf("Failed to query trap counts: %w", err)
	}

	traps, err := toTemplateTrapsNearby(locations, trap_data, counts)
	if err != nil {
		return nil, fmt.Errorf("Failed to convert trap data: %w", err)
	}
	return traps, nil
}

func TreatmentsBySource(ctx context.Context, org Organization, sourceID uuid.UUID) ([]Treatment, error) {
	rows, err := org.model.Treatments(
		sm.Where(
			models.FieldseekerTreatments.Columns.Pointlocid.EQ(psql.Arg(sourceID)),
		),
		sm.OrderBy("enddatetime").Desc(),
	).All(ctx, db.PGInstance.BobDB)
	if err != nil {
		return nil, fmt.Errorf("Failed to query rows: %w", err)
	}
	return toTreatment(rows)
}

func TrapByGlobalId(ctx context.Context, org Organization, id uuid.UUID) (*Trap, error) {
	trap_location, err := org.model.Traplocations(
		sm.Where(models.FieldseekerTraplocations.Columns.Globalid.EQ(psql.Arg(id))),
	).One(ctx, db.PGInstance.BobDB)
	if err != nil {
		return nil, fmt.Errorf("Failed to get trap location: %w", err)
	}

	trap_data, err := sql.TrapDataByLocationIDRecent(org.ID(), []uuid.UUID{id}).All(ctx, db.PGInstance.BobDB)
	if err != nil {
		return nil, fmt.Errorf("Failed to query trap data: %w", err)
	}

	counts, err := sql.TrapCountByLocationID(org.ID(), []uuid.UUID{id}).All(ctx, db.PGInstance.BobDB)
	if err != nil {
		return nil, fmt.Errorf("Failed to query trap counts: %w", err)
	}
	result, err := toTrap(trap_location, trap_data, counts)
	if err != nil {
		return nil, fmt.Errorf("to trap: %w", err)
	}
	return &result, err
}

func TrapsByCell(ctx context.Context, org Organization, c h3.Cell) (results []TrapSummary, err error) {
	boundary, err := c.Boundary()
	if err != nil {
		return results, fmt.Errorf("Failed to get cell boundary: %w", err)
	}
	geom_query := gisStatement(boundary)
	rows, err := org.model.Traplocations(
		sm.Where(
			psql.F("ST_Within", "geospatial", geom_query),
		),
		sm.OrderBy("objectid"),
	).All(ctx, db.PGInstance.BobDB)
	if err != nil {
		return results, fmt.Errorf("Failed to query rows: %w", err)
	}
	return toTemplateTrapSummary(rows)
}

func TreatmentsByCell(ctx context.Context, org Organization, c h3.Cell) ([]Treatment, error) {
	var results []Treatment
	boundary, err := c.Boundary()
	if err != nil {
		return results, fmt.Errorf("Failed to get cell boundary: %w", err)
	}
	geom_query := gisStatement(boundary)
	rows, err := org.model.Treatments(
		sm.Where(
			psql.F("ST_Within", "geospatial", geom_query),
		),
		sm.OrderBy("pointlocid"),
		sm.OrderBy("enddatetime"),
	).All(ctx, db.PGInstance.BobDB)
	if err != nil {
		return results, fmt.Errorf("Failed to query rows: %w", err)
	}
	return toTreatment(rows)
}
func InspectionsByCell(ctx context.Context, org Organization, c h3.Cell) ([]Inspection, error) {
	var results []Inspection

	boundary, err := c.Boundary()
	if err != nil {
		return results, fmt.Errorf("Failed to get cell boundary: %w", err)
	}
	geom_query := gisStatement(boundary)
	rows, err := org.model.Mosquitoinspections(
		sm.Where(
			psql.F("ST_Within", "geospatial", geom_query),
		),
		sm.OrderBy("pointlocid"),
		sm.OrderBy("enddatetime"),
	).All(ctx, db.PGInstance.BobDB)
	if err != nil {
		return results, fmt.Errorf("Failed to query rows: %w", err)
	}
	return toTemplateInspection(rows)
}
func InspectionsBySource(ctx context.Context, org Organization, sourceID uuid.UUID) ([]Inspection, error) {
	var results []Inspection

	rows, err := org.model.Mosquitoinspections(
		sm.Where(
			models.FieldseekerMosquitoinspections.Columns.Pointlocid.EQ(psql.Arg(sourceID)),
		),
		sm.OrderBy("enddatetime").Desc(),
	).All(ctx, db.PGInstance.BobDB)
	if err != nil {
		return results, fmt.Errorf("Failed to query rows: %w", err)
	}
	return toTemplateInspection(rows)
}
func toBreedingSourceSummary(points []*models.FieldseekerPointlocation) []BreedingSourceSummary {
	results := make([]BreedingSourceSummary, len(points))
	for i, r := range points {
		var last_inspected *time.Time
		if !r.Lastinspectdate.IsNull() {
			l := r.Lastinspectdate.MustGet()
			last_inspected = &l
		}
		var last_treat_date *time.Time
		if !r.Lasttreatdate.IsNull() {
			l := r.Lasttreatdate.MustGet()
			last_treat_date = &l
		}
		results[i] = BreedingSourceSummary{
			ID:            r.Globalid,
			LastInspected: last_inspected,
			LastTreated:   last_treat_date,
			Type:          r.Habitat.GetOr("none"),
		}
	}
	return results
}
func gisStatement(cb h3.CellBoundary) string {
	var content strings.Builder
	for i, p := range cb {
		if i != 0 {
			content.WriteString(", ")
		}
		content.WriteString(fmt.Sprintf("%f %f", p.Lng, p.Lat))
	}
	// Repeat the first coordinate to close the polygon
	content.WriteString(fmt.Sprintf(", %f %f", cb[0].Lng, cb[0].Lat))
	return fmt.Sprintf("ST_GeomFromText('POLYGON((%s))', 3857)", content.String())
}
