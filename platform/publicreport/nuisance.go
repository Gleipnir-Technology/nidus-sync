package publicreport

import (
	"context"
	"fmt"

	"github.com/Gleipnir-Technology/bob"
	"github.com/Gleipnir-Technology/bob/dialect/psql"
	"github.com/Gleipnir-Technology/bob/dialect/psql/sm"
	//"github.com/Gleipnir-Technology/nidus-sync/config"
	"github.com/Gleipnir-Technology/nidus-sync/db"
	"github.com/Gleipnir-Technology/nidus-sync/platform/types"
	//"github.com/google/uuid"
	//"github.com/rs/zerolog/log"
	"github.com/stephenafamo/scan"
)

func nuisancesByReportID(ctx context.Context, report_ids []int32) (map[int32]*types.Nuisance, error) {
	rows, err := bob.All(ctx, db.PGInstance.BobDB, psql.Select(
		sm.Columns(
			"additional_info",
			"duration",
			"is_location_backyard",
			"is_location_frontyard",
			"is_location_garden",
			"is_location_other",
			"is_location_pool",
			"report_id",
			"source_container",
			"source_description",
			"source_gutter",
			"source_stagnant",
			"tod_day",
			"tod_early",
			"tod_evening",
			"tod_night",
		),
		sm.From("publicreport.nuisance"),
		sm.Where(psql.Quote("report_id").EQ(
			psql.Any(report_ids),
		)),
	), scan.StructMapper[types.Nuisance]())
	if err != nil {
		return nil, fmt.Errorf("query nuisance: %w", err)
	}
	results := make(map[int32]*types.Nuisance, len(rows))
	for _, row := range rows {
		results[row.ReportID] = &types.Nuisance{
			AdditionalInfo:      row.AdditionalInfo,
			Duration:            row.Duration,
			IsLocationBackyard:  row.IsLocationBackyard,
			IsLocationFrontyard: row.IsLocationFrontyard,
			IsLocationGarden:    row.IsLocationGarden,
			IsLocationOther:     row.IsLocationOther,
			IsLocationPool:      row.IsLocationPool,
			SourceContainer:     row.SourceContainer,
			SourceDescription:   row.SourceDescription,
			SourceGutter:        row.SourceGutter,
			SourceStagnant:      row.SourceStagnant,
			TODDay:              row.TODDay,
			TODEarly:            row.TODEarly,
			TODEvening:          row.TODEvening,
			TODNight:            row.TODNight,
		}
	}
	return results, nil
}
