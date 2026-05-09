package db

import (
	"context"
	"fmt"

	//"github.com/Gleipnir-Technology/jet"
	"github.com/Gleipnir-Technology/jet/postgres"
)

type Updater[T postgres.Table, M any] struct {
	Columns postgres.ColumnList
	//Columns []jet.Column
	Model M
	Table T

	buildWhere func(pk_values ...interface{}) postgres.BoolExpression
}

func (u Updater[T, M]) Execute(ctx context.Context, txn Ex, pk_values ...interface{}) error {
	// We get syntax errors from the database if there are no updates to perform
	if u.Columns == nil {
		return fmt.Errorf("nil columns")
	}
	if len(u.Columns) == 0 {
		return nil
	}
	statement := u.Table.
		UPDATE(u.Columns).
		MODEL(u.Model).
		WHERE(u.buildWhere(pk_values...))
	return ExecuteNoneTx(ctx, txn, statement)
}
func (u Updater[T, M]) Has(c postgres.Column) bool {
	for _, col := range u.Columns {
		if col == c {
			return true
		}
	}
	return false
}
func (u *Updater[T, M]) Set(c postgres.Column) {
	u.Columns = append(u.Columns, c)
}
func (u *Updater[T, M]) Unset(c postgres.Column) {
	var index = -1
	for i, col := range u.Columns {
		if col == c {
			index = i
		}
	}
	if index > -1 {
		u.Columns[index] = u.Columns[len(u.Columns)-1]
		u.Columns = u.Columns[:len(u.Columns)-1]
	}
}
func NewUpdater[T postgres.Table, M any](
	table *T,
	pk_columns ...postgres.ColumnInteger,
) Updater[T, M] {
	return Updater[T, M]{
		Columns: postgres.ColumnList{},
		Table:   *table,
		buildWhere: func(pk_values ...interface{}) postgres.BoolExpression {
			conditions := make([]postgres.BoolExpression, len(pk_columns))
			for i, col := range pk_columns {
				conditions[i] = col.EQ(postgres.Int64(pk_values[i].(int64)))
			}
			return postgres.AND(conditions...)
		},
	}
}
