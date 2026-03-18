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

type Water struct {
	AccessComments         string        `db:"access_comments" json:"access_comments"`
	AccessGate             bool          `db:"access_gate" json:"access_gate"`
	AccessFence            bool          `db:"access_fence" json:"access_fence"`
	AccessLocked           bool          `db:"access_locked" json:"access_locked"`
	AccessDog              bool          `db:"access_dog" json:"access_dog"`
	AccessOther            bool          `db:"access_other" json:"access_other"`
	Comments               string        `db:"comments" json:"comments"`
	HasAdult               bool          `db:"has_adult" json:"has_adult"`
	HasBackyardPermission  bool          `db:"has_backyard_permission" json:"has_backyard_permission"`
	HasLarvae              bool          `db:"has_larvae" json:"has_larvae"`
	HasPupae               bool          `db:"has_pupae" json:"has_pupae"`
	IsReporterConfidential bool          `db:"is_reporter_confidential" json:"is_reporter_confidential"`
	IsReporterOwner        bool          `db:"is_reporter_owner" json:"is_reporter_owner"`
	Owner                  types.Contact `db:"owner" json:"owner"`
	ReportID               int32         `db:"report_id" json:"-"`
}

func watersByReportID(ctx context.Context, report_ids []int32) (map[int32]*Water, error) {
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
	), scan.StructMapper[Water]())
	if err != nil {
		return nil, fmt.Errorf("query water: %w", err)
	}
	results := make(map[int32]*Water, len(rows))
	for _, row := range rows {
		results[row.ReportID] = &Water{
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
