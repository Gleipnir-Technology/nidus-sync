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
	//"github.com/Gleipnir-Technology/nidus-sync/db/models"
	"github.com/Gleipnir-Technology/nidus-sync/platform/types"
	//"github.com/google/uuid"
	//"github.com/rs/zerolog/log"
	"github.com/stephenafamo/scan"
)

type Report struct {
	Log        []LogEntry      `db:"-" json:"log"`
	Address    types.Address   `db:"address" json:"address"`
	AddressRaw string          `db:"address_raw" json:"address_raw"`
	Created    time.Time       `db:"created" json:"created"`
	ID         int32           `db:"id" json:"-"`
	Images     []types.Image   `db:"images" json:"images"`
	Location   *types.Location `db:"location" json:"location"`
	Nuisance   *Nuisance       `db:"nuisance" json:"nuisance"`
	PublicID   string          `db:"public_id" json:"public_id"`
	Reporter   types.Contact   `db:"reporter" json:"reporter"`
	Status     string          `db:"status" json:"status"`
	Type       string          `db:"report_type" json:"type"`
	Water      *Water          `db:"water" json:"water"`
}

func ReportsForOrganization(ctx context.Context, org_id int32) ([]*Report, error) {
	rows, err := bob.All(ctx, db.PGInstance.BobDB, psql.Select(
		sm.Columns(
			"address_country AS \"address.country\"",
			"address_locality AS \"address.locality\"",
			"address_number AS \"address.number\"",
			"address_postal_code AS \"address.postal_code\"",
			"address_raw AS address_raw",
			"address_region AS \"address.region\"",
			"address_street AS \"address.street\"",
			"created",
			"id",
			"COALESCE(ST_Y(location::geometry::geometry(point, 4326)), 0.0) AS \"location.latitude\"",
			"COALESCE(ST_X(location::geometry::geometry(point, 4326)), 0.0) AS \"location.longitude\"",
			"public_id",
			"report_type",
			"reporter_email AS \"reporter.email\"",
			"reporter_name AS \"reporter.name\"",
			"reporter_phone AS \"reporter.phone\"",
			"status",
		),
		sm.From("publicreport.report"),
		sm.Where(psql.Quote("publicreport", "report", "organization_id").EQ(psql.Arg(org_id))),
		sm.Where(psql.Quote("publicreport", "report", "reviewed").IsNull()),
	), scan.StructMapper[Report]())

	if err != nil {
		return nil, fmt.Errorf("get reports: %w", err)
	}
	report_ids := make([]int32, len(rows))
	for i, row := range rows {
		report_ids[i] = row.ID
	}
	images_by_id, err := loadImagesForReport(ctx, org_id, report_ids)
	if err != nil {
		return nil, fmt.Errorf("images for report: %w", err)
	}
	logs_by_report_id, err := logEntriesByReportID(ctx, report_ids)
	if err != nil {
		return nil, fmt.Errorf("log entries for reports: %w", err)
	}
	nuisances_by_report_id, err := nuisancesByReportID(ctx, report_ids)
	if err != nil {
		return nil, fmt.Errorf("nuisances: %w", err)
	}
	waters_by_report_id, err := watersByReportID(ctx, report_ids)
	if err != nil {
		return nil, fmt.Errorf("waters: %w", err)
	}

	results := make([]*Report, len(rows))
	for i, row := range rows {
		images, ok := images_by_id[row.ID]
		if ok {
			row.Images = images
		} else {
			row.Images = []types.Image{}
		}
		row.Log = logs_by_report_id[row.ID]
		row.Nuisance = nuisances_by_report_id[row.ID]
		row.Water = waters_by_report_id[row.ID]
		if row.Location.Latitude == 0.0 || row.Location.Longitude == 0.0 {
			row.Location = nil
		}
		results[i] = &row
	}
	return results, nil
}
func ReportsForOrganizationCount(ctx context.Context, org_id int32) (uint, error) {
	type _Row struct {
		Count uint `db:"count"`
	}
	row, err := bob.One(ctx, db.PGInstance.BobDB, psql.Select(
		sm.Columns(
			"COUNT(*) AS count",
		),
		sm.From("publicreport.report"),
		sm.Where(psql.Quote("publicreport", "report", "organization_id").EQ(psql.Arg(org_id))),
	), scan.StructMapper[_Row]())
	if err != nil {
		return 0, fmt.Errorf("query count: %w", err)
	}
	return row.Count, nil
}
