package publicreport

import (
	"context"
	"fmt"
	"time"

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
	AccessComments         string         `db:"access_comments" json:"access_comments"`
	AccessGate             bool           `db:"access_gate" json:"access_gate"`
	AccessFence            bool           `db:"access_fence" json:"access_fence"`
	AccessLocked           bool           `db:"access_locked" json:"access_locked"`
	AccessDog              bool           `db:"access_dog" json:"access_dog"`
	AccessOther            bool           `db:"access_other" json:"access_other"`
	Address                types.Address  `db:"address" json:"address"`
	AddressRaw             string         `db:"address_raw" json:"address_raw"`
	Comments               string         `db:"comments" json:"comments"`
	Created                time.Time      `db:"created" json:"created"`
	HasAdult               bool           `db:"has_adult" json:"has_adult"`
	HasBackyardPermission  bool           `db:"has_backyard_permission" json:"has_backyard_permission"`
	HasLarvae              bool           `db:"has_larvae" json:"has_larvae"`
	HasPupae               bool           `db:"has_pupae" json:"has_pupae"`
	ID                     int32          `db:"id" json:"-"`
	Images                 []types.Image  `db:"-" json:"images"`
	IsReporterConfidential bool           `db:"is_reporter_confidential" json:"is_reporter_confidential"`
	IsReporterOwner        bool           `db:"is_reporter_owner" json:"is_reporter_owner"`
	Location               types.Location `db:"location" json:"location"`
	Owner                  types.Contact  `db:"owner" json:"owner"`
	PublicID               string         `db:"public_id" json:"public_id"`
	Reporter               types.Contact  `db:"reporter" json:"reporter"`
	ReporterContactConsent *bool          `db:"reporter_contact_consent" json:"reporter_contact_consent"`
	Status                 string         `db:"status" json:"status"`
}

func WaterReportForOrganization(ctx context.Context, org_id int32) ([]Water, error) {
	reports, err := bob.All(ctx, db.PGInstance.BobDB, psql.Select(
		sm.Columns(
			"access_comments",
			"access_gate",
			"access_fence",
			"access_locked",
			"access_dog",
			"access_other",
			"access_gate AS address_raw",
			"address_country AS \"address.country\"",
			"address_locality AS \"address.locality\"",
			"address_number AS \"address.number\"",
			"address_postal_code AS \"address.postal_code\"",
			"address_region AS \"address.region\"",
			"address_street AS \"address.street\"",
			"comments",
			"created",
			"has_adult",
			"has_backyard_permission",
			"has_larvae",
			"has_pupae",
			"id",
			"is_reporter_confidential",
			"is_reporter_owner",
			"ST_Y(location::geometry::geometry(point, 4326)) AS \"location.latitude\"",
			"ST_X(location::geometry::geometry(point, 4326)) AS \"location.longitude\"",
			"owner_email AS \"owner.email\"",
			"owner_name AS \"owner.name\"",
			"owner_phone AS \"owner.phone\"",
			"public_id",
			"reporter_email AS \"reporter.email\"",
			"reporter_name AS \"reporter.name\"",
			"reporter_phone AS \"reporter.phone\"",
			"reporter_contact_consent",
			"status",
		),
		sm.From("publicreport.water"),
		sm.Where(psql.Quote("publicreport", "water", "organization_id").EQ(psql.Arg(org_id))),
	), scan.StructMapper[Water]())
	if err != nil {
		return nil, fmt.Errorf("get reports: %w", err)
	}
	report_ids := make([]int32, len(reports))
	for i, report := range reports {
		report_ids[i] = report.ID
	}
	images_by_id, err := loadImagesForReportWater(ctx, org_id, report_ids)
	if err != nil {
		return nil, fmt.Errorf("images for report: %w", err)
	}
	for i := range reports {
		reports[i].Images = images_by_id[reports[i].ID]
	}
	return reports, nil
}
func WaterReportForOrganizationCount(ctx context.Context, org_id int32) (uint, error) {
	type _Row struct {
		Count uint `db:"count"`
	}
	row, err := bob.One(ctx, db.PGInstance.BobDB, psql.Select(
		sm.Columns(
			"COUNT(*) AS count",
		),
		sm.From("publicreport.water"),
		sm.Where(psql.Quote("publicreport", "water", "organization_id").EQ(psql.Arg(org_id))),
	), scan.StructMapper[_Row]())
	if err != nil {
		return 0, fmt.Errorf("query count: %w", err)
	}
	return row.Count, nil
}
