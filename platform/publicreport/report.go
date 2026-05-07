package publicreport

import (
	"context"
	"fmt"

	"github.com/Gleipnir-Technology/bob"
	"github.com/Gleipnir-Technology/bob/dialect/psql"
	"github.com/Gleipnir-Technology/bob/dialect/psql/sm"
	"github.com/Gleipnir-Technology/nidus-sync/db"
	querypublic "github.com/Gleipnir-Technology/nidus-sync/db/query/public"
	querypublicreport "github.com/Gleipnir-Technology/nidus-sync/db/query/publicreport"
	modelpublic "github.com/Gleipnir-Technology/nidus-sync/db/gen/nidus-sync/public/model"
	modelpublicreport "github.com/Gleipnir-Technology/nidus-sync/db/gen/nidus-sync/publicreport/model"
	"github.com/Gleipnir-Technology/nidus-sync/platform/types"
	//"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"github.com/stephenafamo/scan"
)

func ByIDCompliance(ctx context.Context, public_id string, is_public bool) (*types.PublicReportCompliance, error) {
	report, err := byID(ctx, public_id, is_public)
	if err != nil {
		return nil, fmt.Errorf("base report byid: %w", err)
	}
	if report == nil {
		return nil, nil
	}
	return compliance(ctx, public_id, *report)
}
func ByIDNuisance(ctx context.Context, public_id string, is_public bool) (*types.PublicReportNuisance, error) {
	report, err := byID(ctx, public_id, is_public)
	if err != nil {
		return nil, fmt.Errorf("base report byid: %w", err)
	}
	if report == nil {
		return nil, nil
	}
	return nuisance(ctx, public_id, *report)
}
func ByIDWater(ctx context.Context, public_id string, is_public bool) (*types.PublicReportWater, error) {
	report, err := byID(ctx, public_id, is_public)
	if err != nil {
		return nil, fmt.Errorf("base report byid: %w", err)
	}
	if report == nil {
		return nil, nil
	}
	return water(ctx, public_id, *report)
}
func UnreviewedForOrganization(ctx context.Context, txn db.Ex, org_id int64, is_public bool) ([]types.PublicReport, error) {
	reports, err := querypublicreport.ReportsUnreviewedForOrganization(ctx, txn, org_id)
	if err != nil {
		return nil, fmt.Errorf("reports unreviewed: %w", err)
	}
	return reportQueryToRows(ctx, reports, is_public)
}
func byID(ctx context.Context, public_id string, is_public bool) (*types.PublicReport, error) {
	report, err := querypublicreport.ReportFromPublicID(ctx, db.PGInstance.PGXPool, public_id)
	if err != nil {
		return nil, fmt.Errorf("query report from public ID: %w", err)
	}
	if report == nil {
		return nil, nil
	}
	reports, err := reportQueryToRows(ctx, []modelpublicreport.Report{*report}, is_public)
	if err != nil {
		return nil, fmt.Errorf("query to rows: %w", err)
	}
	log.Debug().Str("public_id", public_id).Int("len", len(reports)).Msg("querying for publicreport by ID")
	if len(reports) != 1 {
		return nil, nil
	}
	return &reports[0], nil
}
func reportQueryToRows(ctx context.Context, reports []modelpublicreport.Report, is_public bool) ([]types.PublicReport, error) {
	address_ids := make([]int64, 0)
	report_ids := make([]int32, len(reports))
	for i, report := range reports {
		report_ids[i] = report.ID
		if report.AddressID != nil {
			address_ids = append(address_ids, int64(*report.AddressID))
		}
	}
	images_by_id, err := loadImagesForReport(ctx, report_ids)
	if err != nil {
		return nil, fmt.Errorf("images for report: %w", err)
	}
	logs_by_report_id, err := logEntriesByReportID(ctx, report_ids, is_public)
	if err != nil {
		return nil, fmt.Errorf("log entries for reports: %w", err)
	}
	addresses, err := querypublic.AddressesFromIDs(ctx, db.PGInstance.PGXPool, address_ids)
	if err != nil {
		return nil, fmt.Errorf("addresses for reports: %w", err)
	}
	addresses_by_id := make(map[int64]modelpublic.Address, 0)
	for _, address := range addresses {
		addresses_by_id[int64(address.ID)] = address
	}

	results := make([]types.PublicReport, len(reports))
	for i, row := range reports {
		var images []types.Image
		images, ok := images_by_id[row.ID]
		if !ok {
			images = []types.Image{}
		}
		logs, ok := logs_by_report_id[row.ID]
		if !ok {
			return nil, fmt.Errorf("impossible, missing logs for %d", row.ID)
		}
		var location *types.Location
		if row.Location == nil {
			location = nil
		}
		var address *types.Address
		if row.AddressID != nil {
			addr, ok := addresses_by_id[int64(*row.AddressID)]
			if !ok {
				return nil, fmt.Errorf("impossible, missing address %d", row.AddressID)
			}
			a := types.AddressFromModel(addr)
			address = &a
		}
		if address == nil {
			return nil, fmt.Errorf("nil address: %w", err)
		}
		results[i] = types.PublicReport{
			Address: *address,
			Concerns: nil,
			Created: row.Created,
			ID: row.ID,
			Images: images,
			Location: location,
			Log: logs,
			DistrictID: &row.OrganizationID,
			District: nil,
			PublicID: row.PublicID,
			Reporter: types.Contact{
				CanSMS: &row.ReporterPhoneCanSms,
				Email: &row.ReporterEmail,
				HasEmail: row.ReporterEmail != "",
				HasPhone: row.ReporterPhone != "",
				Name: &row.ReporterName,
				Phone: &row.ReporterPhone,
			},
			Status: row.Status.String(),
			Type: row.ReportType.String(),
			URI: "",
		}
	}
	return results, nil
}
func Reports(ctx context.Context, org_id int64, ids []int64, is_public bool) ([]types.PublicReport, error) {
	reports, err := querypublicreport.ReportsFromIDsForOrg(ctx, db.PGInstance.PGXPool, ids, org_id)
	if err != nil {
		return []types.PublicReport{}, fmt.Errorf("reports from ID for org: %w", err)
	}
	return reportQueryToRows(ctx, reports, is_public)
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
func copyReportContent(src types.PublicReport, dst *types.PublicReport) {
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
