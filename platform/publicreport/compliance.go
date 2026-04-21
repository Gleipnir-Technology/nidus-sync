package publicreport

import (
	"context"
	"fmt"

	"github.com/Gleipnir-Technology/bob"
	"github.com/Gleipnir-Technology/bob/dialect/psql"
	"github.com/Gleipnir-Technology/bob/dialect/psql/sm"
	//"github.com/Gleipnir-Technology/nidus-sync/config"
	"github.com/Gleipnir-Technology/nidus-sync/db"
	"github.com/Gleipnir-Technology/nidus-sync/db/models"
	"github.com/Gleipnir-Technology/nidus-sync/platform/types"
	//"github.com/google/uuid"
	//"github.com/rs/zerolog/log"
	"github.com/stephenafamo/scan"
)

func compliance(ctx context.Context, public_id string, report *types.PublicReport) (*types.PublicReportCompliance, error) {
	row, err := bob.One(ctx, db.PGInstance.BobDB, psql.Select(
		sm.Columns(
			models.PublicreportCompliances.Columns.AccessInstructions,
			models.PublicreportCompliances.Columns.AvailabilityNotes,
			models.PublicreportCompliances.Columns.Comments,
			models.PublicreportCompliances.Columns.GateCode,
			models.PublicreportCompliances.Columns.HasDog,
			models.PublicreportCompliances.Columns.PermissionType,
			models.PublicreportCompliances.Columns.ReportID,
			models.PublicreportCompliances.Columns.ReportPhoneCanText,
			models.PublicreportCompliances.Columns.WantsScheduled,
		),
		sm.From("publicreport.compliance"),
		sm.Where(models.PublicreportCompliances.Columns.ReportID.EQ(
			psql.Arg(report.ID),
		)),
	), scan.StructMapper[types.PublicReportCompliance]())
	if err != nil {
		return nil, fmt.Errorf("query compliance: %w", err)
	}
	copyReportContent(report, &row.PublicReport)
	return &row, nil

}
