package publicreport

import (
/*
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
*/
)

/*
func watersByReportID(ctx context.Context, report_ids []int32) (map[int32]*types.Water, error) {
	rows, err := bob.All(ctx, db.PGInstance.BobDB, psql.Select(
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
			psql.Any(report_ids),
		)),
	), scan.StructMapper[types.Water]())
	if err != nil {
		return nil, fmt.Errorf("query water: %w", err)
	}
	results := make(map[int32]*types.Water, len(rows))
	for _, row := range rows {
		results[row.ReportID] = &types.Water{
			AccessComments:         row.AccessComments,
			AccessGate:             row.AccessGate,
			AccessFence:            row.AccessFence,
			AccessLocked:           row.AccessLocked,
			AccessDog:              row.AccessDog,
			AccessOther:            row.AccessOther,
			Comments:               row.Comments,
			HasAdult:               row.HasAdult,
			HasBackyardPermission:  row.HasBackyardPermission,
			HasLarvae:              row.HasLarvae,
			HasPupae:               row.HasPupae,
			IsReporterConfidential: row.IsReporterConfidential,
			IsReporterOwner:        row.IsReporterOwner,
			Owner:                  row.Owner,
			ReportID:               row.ReportID,
		}
	}
	return results, nil
}
*/
