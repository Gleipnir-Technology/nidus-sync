package publicreport

import (
	"context"
	"errors"
	"fmt"
	//"time"

	//"github.com/Gleipnir-Technology/bob"
	"github.com/Gleipnir-Technology/nidus-sync/db"
	"github.com/Gleipnir-Technology/nidus-sync/db/gen/nidus-sync/publicreport/model"
	"github.com/Gleipnir-Technology/nidus-sync/db/gen/nidus-sync/publicreport/table"
	"github.com/go-jet/jet/v2/postgres"
)

type ReportUpdater = db.Updater[table.ReportTable, model.Report]

func ReportInsert(ctx context.Context, txn db.Ex, m model.Report) (model.Report, error) {
	statement := table.Report.INSERT(table.Report.MutableColumns).
		MODEL(m).
		RETURNING(table.Report.AllColumns)
	return db.ExecuteOneTx[model.Report](ctx, txn, statement)
}
func ReportFromID(ctx context.Context, report_id int64) (model.Report, error) {
	statement := table.Report.SELECT(
		table.Report.AllColumns,
	).FROM(table.Report).
		WHERE(table.Report.ID.EQ(postgres.Int(report_id)))
	return db.ExecuteOne[model.Report](ctx, statement)
}
func ReportsFromIDs(ctx context.Context, report_ids []int64) ([]model.Report, error) {
	sql_ids := make([]postgres.Expression, len(report_ids))
	for i, report_id := range report_ids {
		sql_ids[i] = postgres.Int(report_id)
	}
	statement := table.Report.SELECT(
		table.Report.AllColumns,
	).FROM(table.Report).
		WHERE(table.Report.ID.IN(sql_ids...))
	return db.ExecuteMany[model.Report](ctx, statement)
}
func ReportsFromIDsForOrg(ctx context.Context, txn db.Ex, report_ids []int64, org_id int64) ([]model.Report, error) {
	sql_ids := make([]postgres.Expression, len(report_ids))
	for i, report_id := range report_ids {
		sql_ids[i] = postgres.Int(report_id)
	}
	statement := table.Report.SELECT(
		table.Report.AllColumns,
	).FROM(table.Report).
		WHERE(table.Report.ID.IN(sql_ids...).AND(
			table.Report.OrganizationID.EQ(postgres.Int(org_id))))
	return db.ExecuteManyTx[model.Report](ctx, txn, statement)
}
func ReportFromPublicID(ctx context.Context, txn db.Ex, public_id string) (*model.Report, error) {
	statement := table.Report.SELECT(
		table.Report.AllColumns,
	).FROM(table.Report).
		WHERE(table.Report.PublicID.EQ(postgres.String(public_id)))
	result, err := db.ExecuteOneTx[model.Report](ctx, txn, statement)
	if err != nil {
		if errors.Is(err, db.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("query: %w", err)
	}
	return &result, nil
}
func ReportFromPublicIDForOrg(ctx context.Context, txn db.Ex, public_id string, org_id int64) (*model.Report, error) {
	statement := table.Report.SELECT(
		table.Report.AllColumns,
	).FROM(table.Report).
		WHERE(table.Report.PublicID.EQ(postgres.String(public_id)).AND(
			table.Report.OrganizationID.EQ(postgres.Int(org_id))))
	result, err := db.ExecuteOneTx[model.Report](ctx, txn, statement)
	if err != nil {
		if errors.Is(err, db.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("query: %w", err)
	}
	return &result, nil
}
func ReportsUnreviewedForOrganization(ctx context.Context, txn db.Ex, org_id int64) ([]model.Report, error) {
	statement := table.Report.SELECT(
		table.Report.AllColumns,
	).FROM(table.Report).
		WHERE(table.Report.Reviewed.IS_NULL().AND(
			table.Report.OrganizationID.EQ(postgres.Int(org_id))))
	return db.ExecuteManyTx[model.Report](ctx, txn, statement)
}
