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

func compliance(ctx context.Context, public_id string, report *types.PublicReport) (*types.PublicReportCompliance, error) {
	row, err := bob.One(ctx, db.PGInstance.BobDB, psql.Select(
		sm.Columns(
			"access_instructions",
			"availability_notes",
			"comments",
			"gate_code",
			"has_dog",
			"permission_type",
			"report_id",
			"report_phone_can_text",
			"wants_scheduled",
		),
		sm.From("publicreport.compliance"),
		sm.Where(psql.Quote("report_id").EQ(
			psql.Arg(report.ID),
		)),
	), scan.StructMapper[types.PublicReportCompliance]())
	if err != nil {
		return nil, fmt.Errorf("query compliance: %w", err)
	}
	copyReportContent(report, &row.PublicReport)
	return &row, nil

}
