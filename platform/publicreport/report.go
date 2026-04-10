package publicreport

import (
	"context"
	"fmt"

	"github.com/Gleipnir-Technology/bob"
	"github.com/Gleipnir-Technology/bob/dialect/psql"
	"github.com/Gleipnir-Technology/bob/dialect/psql/dialect"
	"github.com/Gleipnir-Technology/bob/dialect/psql/sm"
	//"github.com/Gleipnir-Technology/nidus-sync/config"
	"github.com/Gleipnir-Technology/nidus-sync/db"
	//"github.com/Gleipnir-Technology/nidus-sync/db/models"
	"github.com/Gleipnir-Technology/nidus-sync/platform/types"
	//"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"github.com/stephenafamo/scan"
)

func ReportsForOrganization(ctx context.Context, org_id int32) ([]*types.Report, error) {
	query := reportQuery()
	query.Apply(
		sm.Where(psql.Quote("publicreport", "report", "organization_id").EQ(psql.Arg(org_id))),
		sm.Where(psql.Quote("publicreport", "report", "reviewed").IsNull()),
	)
	return reportQueryToRows(ctx, query)
}
func reportQueryToRows(ctx context.Context, query bob.BaseQuery[*dialect.SelectQuery]) ([]*types.Report, error) {
	rows, err := bob.All(ctx, db.PGInstance.BobDB, query, scan.StructMapper[types.Report]())

	if err != nil {
		return nil, fmt.Errorf("get reports: %w", err)
	}
	address_ids := make([]int32, 0)
	report_ids := make([]int32, len(rows))
	for i, row := range rows {
		report_ids[i] = row.ID
		if row.Address.ID != nil {
			address_ids = append(address_ids, *row.Address.ID)
		}
	}
	addresses_by_id, err := loadAddresses(ctx, db.PGInstance.BobDB, address_ids)
	if err != nil {
		return nil, fmt.Errorf("addresses by ID: %w", err)
	}
	images_by_id, err := loadImagesForReport(ctx, report_ids)
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

	results := make([]*types.Report, len(rows))
	for i, row := range rows {
		if row.Address.ID != nil {
			address, ok := addresses_by_id[*row.Address.ID]
			if !ok {
				log.Warn().Int32("address.id", *row.Address.ID).Msg("failed to find in addresses_by_id, which means our DB query is wrong")
			} else {
				row.Address = address
			}
		}
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
func Report(ctx context.Context, public_id string) (*types.Report, error) {
	query := reportQuery()
	query.Apply(
		sm.Where(psql.Quote("publicreport", "report", "public_id").EQ(psql.Arg(public_id))),
	)
	reports, err := reportQueryToRows(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("query to rows: %w", err)
	}
	if len(reports) != 1 {
		return nil, fmt.Errorf("reports returned: %d", len(reports))
	}
	return reports[0], nil
}
func Reports(ctx context.Context, org_id int32, ids []int32) ([]*types.Report, error) {
	query := reportQuery()
	query.Apply(
		sm.Where(psql.Quote("publicreport", "report", "organization_id").EQ(psql.Arg(org_id))),
		sm.Where(psql.Quote("publicreport", "report", "id").EQ(psql.Any(ids))),
	)
	return reportQueryToRows(ctx, query)
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
func reportQuery() bob.BaseQuery[*dialect.SelectQuery] {
	return psql.Select(
		sm.Columns(
			"address_country AS \"address.country\"",
			"address_id AS \"address.id\"",
			"address_gid AS \"address.gid\"",
			"address_locality AS \"address.locality\"",
			"address_number AS \"address.number\"",
			"address_postal_code AS \"address.postal_code\"",
			"address_raw AS \"address.raw\"",
			"address_region AS \"address.region\"",
			"address_street AS \"address.street\"",
			"created",
			"id",
			"latlng_accuracy_value AS \"location.accuracy\"",
			"COALESCE(ST_Y(location::geometry::geometry(point, 4326)), 0.0) AS \"location.latitude\"",
			"COALESCE(ST_X(location::geometry::geometry(point, 4326)), 0.0) AS \"location.longitude\"",
			"organization_id",
			"public_id",
			"report_type",
			"reporter_email AS \"reporter.email\"",
			"reporter_name AS \"reporter.name\"",
			"reporter_phone AS \"reporter.phone\"",
			"status",
		),
		sm.From("publicreport.report"),
	)
}
