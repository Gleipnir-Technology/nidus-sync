package main

import (
	"context"
	"database/sql"
	"embed"
	"errors"
	"fmt"
	"io/fs"
	"sync"

	//"github.com/georgysavva/scany/v2/pgxscan"
	//"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
	"github.com/rs/zerolog/log"
	"github.com/stephenafamo/bob"
)

//go:embed migrations/*.sql
var embedMigrations embed.FS

type postgres struct {
	BobDB   bob.DB
	PGXPool *pgxpool.Pool
}

var (
	PGInstance *postgres
	pgOnce     sync.Once
)

func doMigrations(connection_string string) error {
	log.Info().Str("dsn", connection_string).Msg("Connecting to database")
	db, err := sql.Open("pgx", connection_string)
	if err != nil {
		return fmt.Errorf("Failed to open database connection: %w", err)
	}
	defer db.Close()
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

func initializeDatabase(ctx context.Context, uri string) error {
	needs, err := needsMigrations(uri)
	if err != nil {
		return fmt.Errorf("Failed to determine if migrations are needed: %v", err)
	}
	if needs == nil {
		return errors.New("Can't read variable 'needs' - it's nil")
	}
	if *needs {
		//return errors.New(fmt.Sprintf("Must migrate database before connecting: %t", *needs))
		log.Info().Msg("Handling database migrations")
		err = doMigrations(uri)
		if err != nil {
			return fmt.Errorf("Failed to handle migrations: %v", err)
		}
	} else {
		log.Info().Msg("No database migrations necessary")
	}

	pgOnce.Do(func() {
		db, e := pgxpool.New(ctx, uri)
		bobDB := bob.NewDB(stdlib.OpenDBFromPool(db))
		PGInstance = &postgres{bobDB, db}
		err = e
	})
	if err != nil {
		return fmt.Errorf("unable to create connection pool: %w", err)
	}

	var current string
	query := `SELECT current_database()`
	err = PGInstance.BobDB.QueryRow(query).Scan(&current)
	if err != nil {
		return fmt.Errorf("Failed to get database current: %w", err)
	}
	log.Info().Str("database", current).Msg("Connected to database")
	return nil
}

func needsMigrations(connection_string string) (*bool, error) {
	log.Info().Str("dsn", connection_string).Msg("Connecting to database")
	db, err := sql.Open("pgx", connection_string)
	if err != nil {
		return nil, fmt.Errorf("Failed to open database connection: %w", err)
	}
	defer db.Close()
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
