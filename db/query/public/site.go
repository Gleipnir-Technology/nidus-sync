package public

import (
	"context"
	"errors"
	"fmt"
	//"time"

	//"github.com/Gleipnir-Technology/bob"
	"github.com/Gleipnir-Technology/nidus-sync/db"
	"github.com/Gleipnir-Technology/nidus-sync/db/gen/nidus-sync/public/model"
	"github.com/Gleipnir-Technology/nidus-sync/db/gen/nidus-sync/public/table"
	"github.com/go-jet/jet/v2/postgres"
)

func SiteFromAddressIDForOrg(ctx context.Context, txn db.Ex, address_id int64, org_id int64) (*model.Site, error) {
	statement := table.Site.SELECT(
		table.Site.AllColumns,
	).FROM(table.Site).
		WHERE(table.Site.AddressID.EQ(postgres.Int(address_id)).AND(
			table.Site.OrganizationID.EQ(postgres.Int(org_id))))
	result, err := db.ExecuteOneTx[model.Site](ctx, txn, statement)
	if err != nil {
		if errors.Is(err, db.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("query: %w", err)
	}
	return &result, nil
}
func SiteFromIDForOrg(ctx context.Context, txn db.Ex, comm_id int64, org_id int64) (model.Site, error) {
	statement := table.Site.SELECT(
		table.Site.AllColumns,
	).FROM(table.Site).
		WHERE(table.Site.ID.EQ(postgres.Int(comm_id)).AND(
			table.Site.OrganizationID.EQ(postgres.Int(org_id))))
	return db.ExecuteOneTx[model.Site](ctx, txn, statement)
}
