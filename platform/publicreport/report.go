package publicreport

import (
	"context"
	"fmt"

	"github.com/Gleipnir-Technology/bob"
	"github.com/Gleipnir-Technology/bob/dialect/psql"
	"github.com/Gleipnir-Technology/bob/dialect/psql/dialect"
	"github.com/Gleipnir-Technology/bob/dialect/psql/sm"
	"github.com/Gleipnir-Technology/nidus-sync/db"
	"github.com/Gleipnir-Technology/nidus-sync/platform/types"
	//"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"github.com/stephenafamo/scan"
)

func ByID(ctx context.Context, public_id string) (*types.PublicReport, error) {
	query := reportQuery()
	query.Apply(
		sm.Where(psql.Quote("r", "public_id").EQ(psql.Arg(public_id))),
	)
	reports, err := reportQueryToRows(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("query to rows: %w", err)
	}
	log.Debug().Str("public_id", public_id).Int("len", len(reports)).Msg("querying for publicreport by ID")
	if len(reports) != 1 {
		return nil, fmt.Errorf("reports returned: %d", len(reports))
	}
	return reports[0], nil
}
func ByIDCompliance(ctx context.Context, public_id string) (*types.PublicReportCompliance, error) {
	report, err := ByID(ctx, public_id)
	if err != nil {
		return nil, fmt.Errorf("base report byid: %w", err)
	}
	return compliance(ctx, public_id, report)
}
func ByIDNuisance(ctx context.Context, public_id string) (*types.PublicReportNuisance, error) {
	report, err := ByID(ctx, public_id)
	if err != nil {
		return nil, fmt.Errorf("base report byid: %w", err)
	}
	return nuisance(ctx, public_id, report)
}
func ByIDWater(ctx context.Context, public_id string) (*types.PublicReportWater, error) {
	report, err := ByID(ctx, public_id)
	if err != nil {
		return nil, fmt.Errorf("base report byid: %w", err)
	}
	return water(ctx, public_id, report)
}
func ReportsForOrganization(ctx context.Context, org_id int32) ([]*types.PublicReport, error) {
	query := reportQuery()
	query.Apply(
		sm.Where(psql.Quote("r", "organization_id").EQ(psql.Arg(org_id))),
		sm.Where(psql.Quote("r", "reviewed").IsNull()),
	)
	return reportQueryToRows(ctx, query)
}
func reportQueryToRows(ctx context.Context, query bob.BaseQuery[*dialect.SelectQuery]) ([]*types.PublicReport, error) {
	rows, err := bob.All(ctx, db.PGInstance.BobDB, query, scan.StructMapper[types.PublicReport]())

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

	results := make([]*types.PublicReport, len(rows))
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
		if row.Location.Latitude == 0.0 || row.Location.Longitude == 0.0 {
			row.Location = nil
		}
		results[i] = &row
	}
	return results, nil
}
func Reports(ctx context.Context, org_id int32, ids []int32) ([]*types.PublicReport, error) {
	query := reportQuery()
	query.Apply(
		sm.Where(psql.Quote("r", "organization_id").EQ(psql.Arg(org_id))),
		sm.Where(psql.Quote("r", "id").EQ(psql.Any(ids))),
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
func copyReportContent(src *types.PublicReport, dst *types.PublicReport) {
	dst.Address = src.Address
	dst.Created = src.Created
	dst.ID = src.ID
	dst.Images = src.Images
	dst.Location = src.Location
	dst.Log = src.Log
	dst.DistrictID = src.DistrictID
	dst.District = src.District
	dst.PublicID = src.PublicID
	dst.Reporter = src.Reporter
	dst.Status = src.Status
	dst.Type = src.Type
	dst.URI = src.URI
}
func reportQuery() bob.BaseQuery[*dialect.SelectQuery] {
	return psql.Select(
		sm.Columns(
			"COALESCE(a.country, '') AS \"address.country\"",
			"a.id AS \"address.id\"",
			"COALESCE(a.gid, '') AS \"address.gid\"",
			"COALESCE(a.locality, '') AS \"address.locality\"",
			"COALESCE(a.number_, '') AS \"address.number\"",
			"COALESCE(a.postal_code, '') AS \"address.postal_code\"",
			"COALESCE(a.region, '') AS \"address.region\"",
			"COALESCE(a.street, '') AS \"address.street\"",
			"r.address_raw AS \"address.raw\"",
			"r.created",
			"r.id",
			"r.latlng_accuracy_value AS \"location.accuracy\"",
			"COALESCE(ST_Y(r.location::geometry::geometry(point, 4326)), 0.0) AS \"location.latitude\"",
			"COALESCE(ST_X(r.location::geometry::geometry(point, 4326)), 0.0) AS \"location.longitude\"",
			"r.organization_id",
			"r.public_id",
			"r.report_type",
			"r.reporter_email AS \"reporter.email\"",
			"r.reporter_name AS \"reporter.name\"",
			"r.reporter_phone AS \"reporter.phone\"",
			"r.status",
		),
		sm.From("publicreport.report").As("r"),
		sm.LeftJoin("address").As("a").OnEQ(
			psql.Quote("r", "address_id"),
			psql.Quote("a", "id"),
		),
	)
}
