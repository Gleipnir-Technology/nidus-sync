package public

import (
	"context"
	"errors"
	"fmt"

	//"github.com/Gleipnir-Technology/bob"
	"github.com/Gleipnir-Technology/nidus-sync/db"
	//"github.com/Gleipnir-Technology/nidus-sync/db/gen/nidus-sync/public/enum"
	"github.com/Gleipnir-Technology/nidus-sync/db/gen/nidus-sync/public/model"
	"github.com/Gleipnir-Technology/nidus-sync/db/gen/nidus-sync/public/table"
	"github.com/go-jet/jet/v2/postgres"
)

func AddressFromComplianceReportRequestID(ctx context.Context, txn db.Ex, public_id string) (model.Address, error) {
	statement := table.Address.SELECT(
		table.Address.AllColumns,
	).FROM(
		table.Address,
		table.Site.INNER_JOIN(
			table.Site,
			table.Site.AddressID.EQ(table.Address.ID)),
		table.Site.INNER_JOIN(
			table.Lead,
			table.Lead.SiteID.EQ(table.Site.ID)),
		table.Lead.INNER_JOIN(
			table.ComplianceReportRequest,
			table.ComplianceReportRequest.LeadID.EQ(table.Lead.ID))).
		WHERE(table.ComplianceReportRequest.PublicID.EQ(postgres.String(public_id)))
	return db.ExecuteOne[model.Address](ctx, statement)
}
func AddressFromGID(ctx context.Context, txn db.Ex, gid string) (*model.Address, error) {
	statement := table.Address.SELECT(
		table.Address.AllColumns,
	).FROM(table.Address).
		WHERE(table.Address.Gid.EQ(postgres.String(gid)))
	result, err := db.ExecuteOneTx[model.Address](ctx, txn, statement)
	if err != nil {
		if errors.Is(err, db.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("execute one: %w", err)
	}
	return &result, nil
}
func AddressFromID(ctx context.Context, txn db.Ex, comm_id int64) (model.Address, error) {
	statement := table.Address.SELECT(
		table.Address.AllColumns,
	).FROM(table.Address).
		WHERE(table.Address.ID.EQ(postgres.Int(comm_id)))
	return db.ExecuteOne[model.Address](ctx, statement)
}
func AddressesFromIDs(ctx context.Context, txn db.Ex, address_ids []int64) ([]model.Address, error) {
	sql_ids := make([]postgres.Expression, len(address_ids))
	for i, address_id := range address_ids {
		sql_ids[i] = postgres.Int(address_id)
	}
	statement := table.Address.SELECT(
		table.Address.AllColumns,
	).FROM(table.Address).
		WHERE(table.Address.ID.IN(sql_ids...))
	return db.ExecuteManyTx[model.Address](ctx, txn, statement)
}
func AddressInsert(ctx context.Context, txn db.Ex, a model.Address) (model.Address, error) {
	statement := table.Address.
		INSERT(table.Address.MutableColumns).
		MODEL(a).
		RETURNING(table.Address.AllColumns)
	return db.ExecuteOneTx[model.Address](ctx, txn, statement)
}
func AddressInserts(ctx context.Context, txn db.Ex, addresses []model.Address) ([]model.Address, error) {
	statement := table.Address.
		INSERT(table.Address.MutableColumns).
		MODELS(addresses).
		RETURNING(table.Address.AllColumns)
	return db.ExecuteManyTx[model.Address](ctx, txn, statement)
}
