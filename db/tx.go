package db

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/rs/zerolog/log"
	//"github.com/stephenafamo/scan"
)

type Ex interface {
	Exec(ctx context.Context, sql string, arguments ...any) (commandTag pgconn.CommandTag, err error)
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
}
type Tx struct {
	pgx.Tx
}

func (txn Tx) ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error) {
	result, err := txn.Tx.Exec(ctx, query, args...)
	if err != nil {
		return Result{}, fmt.Errorf("exec: %w", err)
	}
	return Result{
		tag: result,
	}, nil
}

/*
	func (txn Tx) QueryContext(ctx context.Context, query string, args ...any) (scan.Rows, error) {
		result, err := txn.Tx.Exec(ctx, query, args...)
		return Rows{
			tag: result,
		}, err
	}
*/
type Result struct {
	tag pgconn.CommandTag
}

func (r Result) LastInsertId() (int64, error) {
	log.Debug().Msg("queried last insert id. erroring...")
	return 0, fmt.Errorf("not implemented")
}
func (r Result) RowsAffected() (int64, error) {
	rows := r.tag.RowsAffected()
	log.Debug().Int64("rows", rows).Msg("queried rows affected")
	return rows, nil
}

type Rows struct {
}

func (r Rows) Close() error {
	log.Debug().Msg("requested close of rows")
	return nil
}
func (r Rows) Columns() ([]string, error) {
	log.Debug().Msg("requested columns")
	return []string{}, nil
}
func (r Rows) Err() error {
	log.Debug().Msg("requested err")
	return nil
}
func (r Rows) Next() bool {
	log.Debug().Msg("requested next")
	return false
}
func (r Rows) Scan(args ...any) error {
	log.Debug().Msg("requested scan")
	return fmt.Errorf("scan not implemented")
}
func BeginTxn(ctx context.Context) (Tx, error) {
	txn, err := PGInstance.PGXPool.BeginTx(ctx, pgx.TxOptions{})
	return Tx{
		Tx: txn,
	}, err
}
