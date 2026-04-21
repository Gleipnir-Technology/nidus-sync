package platform

import (
	"context"
	"fmt"

	"github.com/Gleipnir-Technology/bob"
	"github.com/Gleipnir-Technology/bob/dialect/psql"
	"github.com/Gleipnir-Technology/bob/dialect/psql/dialect"
	"github.com/Gleipnir-Technology/bob/dialect/psql/sm"
	"github.com/Gleipnir-Technology/nidus-sync/db"
	"github.com/Gleipnir-Technology/nidus-sync/db/models"
	"github.com/Gleipnir-Technology/nidus-sync/platform/types"
	"github.com/stephenafamo/scan"
)

func MailerByID(ctx context.Context, user User, id int32) (*types.Mailer, error) {
	query := mailerQuery()
	query.Apply(
		sm.Where(models.ComplianceReportRequests.Columns.ID.EQ(psql.Arg(id))),
		sm.Where(
			models.Sites.Columns.OrganizationID.EQ(psql.Arg(user.Organization.ID)),
		),
	)
	mailers, err := mailerQueryToRows(ctx, query)
	if err != nil {
		return nil, err
	}
	return mailers[id], nil
}
func MailerList(ctx context.Context, user User, limit int) ([]*types.Mailer, error) {
	query := mailerQuery()
	query.Apply(
		sm.Where(
			models.Sites.Columns.OrganizationID.EQ(psql.Arg(user.Organization.ID)),
		),
		sm.OrderBy(models.ComplianceReportRequests.Columns.Created),
		sm.Limit(limit),
	)
	return mailerQueryToRows(ctx, query)
}
func mailerQuery() bob.BaseQuery[*dialect.SelectQuery] {
	return psql.Select(
		sm.Columns(
			models.Addresses.Columns.Country.As("address.country"),
			models.Addresses.Columns.Locality.As("address.locality"),
			//sm.From(psql.F("COALESCE", psql.Raw("address.location_latitude"), 0)).As("address.location.latitude"),
			//sm.From(psql.F("COALESCE", psql.Raw("address.location_longitude"), 0)).As("address.location.longitude"),
			"COALESCE(address.location_latitude, 0) AS \"address.location.latitude\"",
			"COALESCE(address.location_longitude, 0) AS \"address.location.longitude\"",
			models.Addresses.Columns.Number.As("address.number"),
			models.Addresses.Columns.PostalCode.As("address.postal_code"),
			models.Addresses.Columns.Region.As("address.region"),
			models.Addresses.Columns.Street.As("address.street"),
			models.Addresses.Columns.Unit.As("address.unit"),
			models.ComplianceReportRequests.Columns.Created.As("created"),
			models.ComplianceReportRequests.Columns.PublicID.As("compliance_report_request_id"),
			models.Sites.Columns.ID.As("site_id"),
			models.Sites.Columns.OwnerName.As("recipient"),
			"'created' AS \"status\"",
		),
		sm.From(models.ComplianceReportRequestMailers.Name()),
		sm.InnerJoin(models.ComplianceReportRequests.Name()).OnEQ(
			models.ComplianceReportRequestMailers.Columns.ComplianceReportRequestID,
			models.ComplianceReportRequests.Columns.ID,
		),
		sm.InnerJoin(models.Leads.Name()).OnEQ(
			models.ComplianceReportRequests.Columns.LeadID,
			models.Leads.Columns.ID,
		),
		sm.InnerJoin(models.Sites.Name()).OnEQ(
			models.Leads.Columns.SiteID,
			models.Sites.Columns.ID,
		),
		sm.InnerJoin(models.Addresses.Name()).OnEQ(
			models.Sites.Columns.AddressID,
			models.Addresses.Columns.ID,
		),
	)
}
func mailerQueryToRows(ctx context.Context, query bob.BaseQuery[*dialect.SelectQuery]) ([]*types.Mailer, error) {
	rows, err := bob.All(ctx, db.PGInstance.BobDB, query, scan.StructMapper[*types.Mailer]())
	if err != nil {
		return nil, fmt.Errorf("query mailers: %w", err)
	}

	return rows, nil
}
