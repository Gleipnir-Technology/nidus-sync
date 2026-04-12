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
	//"github.com/Gleipnir-Technology/nidus-sync/db/models"
	//"github.com/google/uuid"
	//"github.com/rs/zerolog/log"
	"github.com/stephenafamo/scan"
)

func water(ctx context.Context, public_id string, report *types.PublicReport) (*types.PublicReportWater, error) {
	row, err := bob.One(ctx, db.PGInstance.BobDB, psql.Select(
		sm.Columns(
			"access_comments",
			"access_gate",
			"access_fence",
			"access_locked",
			"access_dog",
			"access_other",
			"comments",
			"has_adult",
			"has_backyard_permission",
			"has_larvae",
			"has_pupae",
			"is_reporter_confidential",
			"is_reporter_owner",
			"owner_email AS \"owner.email\"",
			"owner_name AS \"owner.name\"",
			"owner_phone AS \"owner.phone\"",
			"report_id",
		),
		sm.From("publicreport.water"),
		sm.Where(psql.Quote("report_id").EQ(
			psql.Arg(report.ID),
		)),
	), scan.StructMapper[types.PublicReportWater]())
	if err != nil {
		return nil, fmt.Errorf("query water: %w", err)
	}
	copyReportContent(report, &row.PublicReport)
	return &row, nil
}
