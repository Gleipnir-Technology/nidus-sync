package background

import (
	"context"
	"fmt"

	"github.com/Gleipnir-Technology/nidus-sync/db"
	"github.com/Gleipnir-Technology/nidus-sync/db/enums"
	"github.com/Gleipnir-Technology/nidus-sync/db/models"
	"github.com/Gleipnir-Technology/nidus-sync/h3utils"
	"github.com/rs/zerolog/log"
	"github.com/stephenafamo/bob"
	"github.com/stephenafamo/bob/dialect/psql"
	"github.com/stephenafamo/bob/dialect/psql/dialect"
	"github.com/stephenafamo/bob/dialect/psql/dm"
	"github.com/stephenafamo/bob/dialect/psql/im"
	"github.com/uber/h3-go/v4"
)

func updateSummaryTables(ctx context.Context, org *models.Organization) {
	updateSummaryMosquitoSource(ctx, org)
	updateSummaryServiceRequest(ctx, org)
	updateSummaryTrap(ctx, org)
}

func aggregateAtResolution(ctx context.Context, resolution int, org_id int32, type_ enums.H3aggregationtype, cells []h3.Cell) error {
	var err error
	log.Info().Int("resolution", resolution).Str("type", string(type_)).Msg("Working summary layer")
	cellToCount := make(map[h3.Cell]int, 0)
	for _, cell := range cells {
		scaled, err := cell.Parent(resolution)
		if err != nil {
			log.Error().Err(err).Int("resolution", resolution).Msg("Failed to get cell's parent at resolution")
			continue
		}
		cellToCount[scaled] = cellToCount[scaled] + 1
	}

	_, err = models.H3Aggregations.Delete(
		dm.Where(
			psql.And(
				models.H3Aggregations.Columns.OrganizationID.EQ(psql.Arg(org_id)),
				models.H3Aggregations.Columns.Resolution.EQ(psql.Arg(resolution)),
				models.H3Aggregations.Columns.Type.EQ(psql.Arg(type_)),
			),
		),
	).Exec(ctx, db.PGInstance.BobDB)
	if err != nil {
		return fmt.Errorf("Failed to clear previous aggregation: %w", err)
	}
	var to_insert []bob.Mod[*dialect.InsertQuery] = make([]bob.Mod[*dialect.InsertQuery], 0)
	to_insert = append(to_insert, im.Into("h3_aggregation", "cell", "resolution", "count_", "type_", "organization_id", "geometry"))
	for cell, count := range cellToCount {
		polygon, err := h3utils.CellToPostgisGeometry(cell)
		if err != nil {
			log.Error().Err(err).Msg("Failed to get PostGIS geometry")
			continue
		}
		// log.Info().Str("polygon", polygon).Msg("Going to insert")
		to_insert = append(to_insert, im.Values(psql.Arg(cell.String(), resolution, count, type_, org_id), psql.F("st_geomfromtext", psql.S(polygon), 4326)))
	}
	to_insert = append(to_insert, im.OnConflict("cell, organization_id, type_").DoUpdate(
		im.SetCol("count_").To(psql.Raw("EXCLUDED.count_")),
	))
	//log.Info().Str("sql", insertQueryToString(psql.Insert(to_insert...))).Msg("Updating...")
	_, err = psql.Insert(to_insert...).Exec(ctx, db.PGInstance.BobDB)
	if err != nil {
		return fmt.Errorf("Failed to add h3 aggregation: %w", err)
	}
	return nil
}

func updateSummaryMosquitoSource(ctx context.Context, org *models.Organization) {
	point_locations, err := org.Pointlocations().All(ctx, db.PGInstance.BobDB)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get all point locations")
		return
	}
	if len(point_locations) == 0 {
		log.Info().Int("org_id", int(org.ID)).Msg("No updates to perform")
		return
	}

	cells := make([]h3.Cell, 0)
	for _, p := range point_locations {
		if p.H3cell.IsNull() {
			continue
		}
		cell, err := h3utils.ToCell(p.H3cell.MustGet())
		if err != nil {
			log.Error().Err(err).Msg("Failed to get geometry point")
			continue
		}
		cells = append(cells, cell)
	}

	for i := range 16 {
		err = aggregateAtResolution(ctx, i, org.ID, enums.H3aggregationtypeMosquitosource, cells)
		if err != nil {
			log.Error().Err(err).Int("resolution", i).Msg("Failed to aggregate mosquito source")
		}
	}
}

func updateSummaryServiceRequest(ctx context.Context, org *models.Organization) {
	service_requests, err := org.Servicerequests().All(ctx, db.PGInstance.BobDB)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get all service requests")
		return
	}
	if len(service_requests) == 0 {
		log.Info().Int("org_id", int(org.ID)).Msg("No updates to perform")
		return
	}

	cells := make([]h3.Cell, 0)
	for _, p := range service_requests {
		if p.H3cell.IsNull() {
			continue
		}
		cell, err := h3utils.ToCell(p.H3cell.MustGet())
		if err != nil {
			log.Error().Err(err).Msg("Failed to get geometry point")
			continue
		}
		cells = append(cells, cell)
	}
	for i := range 16 {
		err = aggregateAtResolution(ctx, i, org.ID, enums.H3aggregationtypeServicerequest, cells)
		if err != nil {
			log.Error().Err(err).Int("resolution", i).Msg("Failed to aggregate service request")
		}
	}
}

func updateSummaryTrap(ctx context.Context, org *models.Organization) {
	traps, err := org.Traplocations().All(ctx, db.PGInstance.BobDB)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get all trap locations")
		return
	}
	if len(traps) == 0 {
		log.Info().Int("org_id", int(org.ID)).Msg("No updates to perform")
		return
	}

	cells := make([]h3.Cell, 0)
	for _, t := range traps {
		if t.H3cell.IsNull() {
			continue
		}
		cell, err := h3utils.ToCell(t.H3cell.MustGet())
		if err != nil {
			log.Error().Err(err).Msg("Failed to get geometry point")
			continue
		}
		cells = append(cells, cell)
	}
	for i := range 16 {
		err = aggregateAtResolution(ctx, i, org.ID, enums.H3aggregationtypeTrap, cells)
		if err != nil {
			log.Error().Err(err).Int("resolution", i).Msg("Failed to aggregate trap")
		}
	}
}
