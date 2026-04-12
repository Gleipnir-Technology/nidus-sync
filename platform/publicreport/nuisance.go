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

func nuisance(ctx context.Context, public_id string, report *types.PublicReport) (*types.PublicReportNuisance, error) {
	row, err := bob.One(ctx, db.PGInstance.BobDB, psql.Select(
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
			psql.Arg(report.ID),
		)),
	), scan.StructMapper[types.PublicReportNuisance]())
	if err != nil {
		return nil, fmt.Errorf("query nuisance: %w", err)
	}
	copyReportContent(report, &row.PublicReport)
	return &row, nil
}
