package platform

import (
	"context"
	"fmt"

	"github.com/Gleipnir-Technology/bob"
	"github.com/Gleipnir-Technology/bob/dialect/psql"
	"github.com/Gleipnir-Technology/bob/dialect/psql/sm"
	"github.com/Gleipnir-Technology/nidus-sync/db"
	"github.com/Gleipnir-Technology/nidus-sync/db/models"
	"github.com/Gleipnir-Technology/nidus-sync/platform/types"
	"github.com/stephenafamo/scan"
)

type Address = types.Address

func AddressFromComplianceReportRequestID(ctx context.Context, public_id string) (*types.Address, error) {
	row, err := bob.One(ctx, db.PGInstance.BobDB, psql.Select(
		sm.Columns(
			models.Addresses.Columns.Country,
			models.Addresses.Columns.Gid,
			models.Addresses.Columns.ID,
			models.Addresses.Columns.Locality,
			models.Addresses.Columns.LocationLatitude.As("location.latitude"),
			models.Addresses.Columns.LocationLongitude.As("location.longitude"),
			models.Addresses.Columns.Number,
			models.Addresses.Columns.PostalCode,
			models.Addresses.Columns.Region,
			models.Addresses.Columns.Street,
			models.Addresses.Columns.Unit,
		),
		//sm.From(models.Addresses.NameAs()),
		sm.From(models.ComplianceReportRequests.NameAs()),
		sm.InnerJoin(models.Leads.NameAs()).On(
			models.ComplianceReportRequests.Columns.LeadID.EQ(models.Leads.Columns.ID)),
		sm.InnerJoin(models.Sites.NameAs()).On(
			models.Leads.Columns.SiteID.EQ(models.Sites.Columns.ID)),
		sm.InnerJoin(models.Addresses.NameAs()).On(
			models.Sites.Columns.AddressID.EQ(models.Addresses.Columns.ID)),
		sm.Where(models.ComplianceReportRequests.Columns.PublicID.EQ(psql.Arg(public_id))),
	), scan.StructMapper[*types.Address]())
	if err != nil {
		return nil, fmt.Errorf("query address from compliance report request: %w", err)
	}
	return row, nil
}
