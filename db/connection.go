package db

import (
	"context"
	"database/sql"
	"embed"
	"errors"
	"fmt"
	"io/fs"

	//"github.com/georgysavva/scany/v2/pgxscan"
	//"github.com/jackc/pgx/v5"
	"github.com/Gleipnir-Technology/bob"
	"github.com/go-jet/jet/v2/postgres"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
	"github.com/rs/zerolog/log"
	"github.com/stephenafamo/scan"
	pgxgeom "github.com/twpayne/pgx-geom"
)

var ErrNoRows = pgx.ErrNoRows

//go:embed migrations/*.sql
var embedMigrations embed.FS

type pginstance struct {
	BobDB   bob.DB
	PGXPool *pgxpool.Pool
}

var (
	PGInstance *pginstance
)

func ExecuteNone(ctx context.Context, stmt postgres.Statement) error {
	query, args := stmt.Sql()

	_, err := PGInstance.PGXPool.Query(ctx, query, args...)
	return err
}
func ExecuteNoneTx(ctx context.Context, txn Ex, stmt postgres.Statement) error {
	query, args := stmt.Sql()

	_, err := txn.Query(ctx, query, args...)
	return err
}
func ExecuteNoneTxBob(ctx context.Context, txn bob.Tx, stmt postgres.Statement) error {
	query, args := stmt.Sql()

	_, err := txn.QueryContext(ctx, query, args...)
	return err
}
func ExecuteOne[T any](ctx context.Context, stmt postgres.Statement) (T, error) {
	query, args := stmt.Sql()

	var result T
	row, err := PGInstance.PGXPool.Query(ctx, query, args...)
	if err != nil {
		return result, fmt.Errorf("execute query: %w", err)
	}
	var collected *T
	collected, err = pgx.CollectOneRow(row, pgx.RowToAddrOfStructByPos[T])
	if err != nil || collected == nil {
		return result, fmt.Errorf("collect row: %w", err)
	}
	return *collected, nil
}
func ExecuteOneTx[T any](ctx context.Context, txn Ex, stmt postgres.Statement) (T, error) {
	query, args := stmt.Sql()

	//result, err := scan.One(ctx, txn, scan.StructMapper[T](), query, args...)
	row, err := txn.Query(ctx, query, args...)
	var result T
	if err != nil {
		return result, fmt.Errorf("txn query: %w", err)
	}
	var collected *T
	collected, err = pgx.CollectOneRow(row, pgx.RowToAddrOfStructByPos[T])
	if err != nil || collected == nil {
		return result, fmt.Errorf("collect row: %w", err)
	}
	return *collected, nil
}
func ExecuteOneTxBob[T any](ctx context.Context, txn bob.Tx, stmt postgres.Statement) (T, error) {
	query, args := stmt.Sql()

	return scan.One(ctx, txn, scan.StructMapper[T](), query, args...)
}
func ExecuteMany[T any](ctx context.Context, stmt postgres.Statement) ([]T, error) {
	query, args := stmt.Sql()

	rows, err := PGInstance.PGXPool.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("execute query: %w", err)
	}
	collected, err := pgx.CollectRows(rows, pgx.RowToAddrOfStructByPos[T])
	if err != nil {
		return []T{}, fmt.Errorf("collect rows: %w", err)
	}
	results := make([]T, len(collected))
	for i, c := range collected {
		if c == nil {
			return results, fmt.Errorf("null collected")
		}
		results[i] = *c
	}
	return results, nil
}
func ExecuteManyTx[T any](ctx context.Context, txn Ex, stmt postgres.Statement) ([]T, error) {
	query, args := stmt.Sql()

	rows, err := txn.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("execute query: %w", err)
	}
	collected, err := pgx.CollectRows(rows, pgx.RowToAddrOfStructByPos[T])
	if err != nil {
		return []T{}, fmt.Errorf("collect rows: %w", err)
	}
	results := make([]T, len(collected))
	for i, c := range collected {
		if c == nil {
			return results, fmt.Errorf("null collected")
		}
		results[i] = *c
	}
	return results, nil
}
func doMigrations(connection_string string) error {
	log.Debug().Str("dsn", connection_string).Msg("Connecting to database")
	db, err := sql.Open("pgx", connection_string)
	if err != nil {
		return fmt.Errorf("Failed to open database connection: %w", err)
	}
	defer func() {
		err := db.Close()
		if err != nil {
			log.Error().Err(err).Msg("failed to close database connection")
		}
	}()
	row := db.QueryRowContext(context.Background(), "SELECT version()")
	var val string
	if err := row.Scan(&val); err != nil {
		return fmt.Errorf("Failed to get database version query result: %w", err)
	}
	log.Info().Str("version", val).Msg("Connected to database")

	fsys, err := fs.Sub(embedMigrations, "migrations")
	if err != nil {
		return fmt.Errorf("Failed to get migrations embedded directory: %w", err)
	}
	provider, err := goose.NewProvider(goose.DialectPostgres, db, fsys)
	if err != nil {
		return fmt.Errorf("Failed to create goose provider: %w", err)
	}
	//goose.SetBaseFS(embedMigrations)

	current, target, err := provider.GetVersions(context.Background())
	if err != nil {
		return fmt.Errorf("Faield to get goose versions: %w", err)
	}
	log.Info().Int("current", int(current)).Int("target", int(target)).Msg("Migration status")
	results, err := provider.Up(context.Background())
	if err != nil {
		return fmt.Errorf("Failed to run migrations: %w", err)
	}
	if len(results) > 0 {
		for _, r := range results {
			log.Info().Int("version", int(r.Source.Version)).Str("direction", r.Direction).Msg("Migration done")
		}
	} else {
		log.Info().Msg("No migrations necessary.")
	}
	return nil
}

func InitializeDatabase(ctx context.Context, uri string) error {
	log.Debug().Str("dsn", uri).Msg("Initializing database")
	needs, err := needsMigrations(uri)
	if err != nil {
		return fmt.Errorf("Failed to determine if migrations are needed: %w", err)
	}
	if needs == nil {
		return errors.New("Can't read variable 'needs' - it's nil")
	}
	if *needs {
		//return errors.New(fmt.Sprintf("Must migrate database before connecting: %t", *needs))
		log.Info().Msg("Handling database migrations")
		err = doMigrations(uri)
		if err != nil {
			return fmt.Errorf("Failed to handle migrations: %w", err)
		}
	} else {
		log.Debug().Msg("No database migrations necessary")
	}

	config, err := pgxpool.ParseConfig(uri)
	if err != nil {
		return fmt.Errorf("parse config: %w", err)
	}
	config.AfterConnect = func(ctx2 context.Context, conn *pgx.Conn) error {
		err2 := pgxgeom.Register(ctx, conn)
		if err2 != nil {
			return fmt.Errorf("pgxgeom register: %w", err2)
		}
		return nil
	}
	db, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return fmt.Errorf("new pool: %w", err)
	}
	bobDB := bob.NewDB(stdlib.OpenDBFromPool(db))
	PGInstance = &pginstance{bobDB, db}

	var current string
	query := `SELECT current_database()`
	err = PGInstance.BobDB.QueryRow(query).Scan(&current)
	if err != nil {
		return fmt.Errorf("Failed to get database current: %w", err)
	}
	return nil
}

func needsMigrations(connection_string string) (*bool, error) {
	db, err := sql.Open("pgx", connection_string)
	if err != nil {
		return nil, fmt.Errorf("Failed to open database connection: %w", err)
	}
	defer func() {
		err := db.Close()
		if err != nil {
			log.Error().Err(err).Msg("failed to close database connection")
		}
	}()
	row := db.QueryRowContext(context.Background(), "SELECT version()")
	var val string
	if err := row.Scan(&val); err != nil {
		return nil, fmt.Errorf("Failed to get database version query result: %w", err)
	}
	log.Info().Str("dsn", val).Msg("Connected to database")

	fsys, err := fs.Sub(embedMigrations, "migrations")
	if err != nil {
		return nil, fmt.Errorf("Failed to get migrations embedded directory: %w", err)
	}
	provider, err := goose.NewProvider(goose.DialectPostgres, db, fsys)
	if err != nil {
		return nil, fmt.Errorf("Failed to create goose provider: %w", err)
	}

	hasPending, err := provider.HasPending(context.Background())
	if err != nil {
		return nil, err
	}
	return &hasPending, nil
}
