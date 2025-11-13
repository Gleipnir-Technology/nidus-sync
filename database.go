package main

import (
	"bytes"
	"context"
	"database/sql"
	"embed"
	"errors"
	"fmt"
	"io/fs"
	"sync"

	//"github.com/georgysavva/scany/v2/pgxscan"
	//"github.com/jackc/pgx/v5"
	"github.com/Gleipnir-Technology/nidus-sync/enums"
	"github.com/Gleipnir-Technology/nidus-sync/models"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
	"github.com/rs/zerolog/log"
	"github.com/stephenafamo/bob"
	"github.com/stephenafamo/bob/dialect/psql"
	"github.com/stephenafamo/bob/dialect/psql/dialect"
	"github.com/stephenafamo/bob/dialect/psql/im"
	"github.com/uber/h3-go/v4"
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

func updateSummaryTables(ctx context.Context, org *models.Organization) {
	/*org, err := models.FindOrganization(ctx, PGInstance.BobDB, org_id)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get organization")
	}*/
	log.Info().Int("org_id", int(org.ID)).Msg("Getting point locations")
	point_locations, err := org.FSPointlocations().All(ctx, PGInstance.BobDB)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get organization")
		return
	}
	log.Info().Int("count", len(point_locations)).Msg("Summarizing point locations")

	for i := range 16 {
		log.Info().Int("resolution", i).Msg("Working summary layer")
		cellToCount := make(map[h3.Cell]int, 0)
		for _, p := range point_locations {
			cell, err := getCell(p.GeometryX, p.GeometryY, i)
			if err != nil {
				log.Error().Err(err).Msg("Failed to get cell")
				continue
			}
			//log.Info().Float64("X", p.GeometryX).Float64("Y", p.GeometryY).Str("cell", cell.String()).Msg("Converted lat/lng")
			cellToCount[cell] = cellToCount[cell] + 1
		}
		var to_insert []bob.Mod[*dialect.InsertQuery] = make([]bob.Mod[*dialect.InsertQuery], 0)
		to_insert = append(to_insert, im.Into("h3_aggregation", "cell", "resolution", "count_", "type_", "organization_id"))
		for cell, count := range cellToCount {
			to_insert = append(to_insert, im.Values(psql.Arg(cell.String(), i, count, enums.H3aggregationtypeServicerequest, org.ID)))
		}
		//to_insert = append(to_insert, im.OnConflict("h3_aggregation_cell_organization_id_type__key").DoUpdate(
		to_insert = append(to_insert, im.OnConflict("cell, organization_id, type_").DoUpdate(
			im.SetCol("count_").To(psql.Raw("EXCLUDED.count_")),
		))
		//log.Info().Str("sql", insertQueryToString(psql.Insert(to_insert...))).Msg("Updating...")
		_, err := psql.Insert(to_insert...).Exec(ctx, PGInstance.BobDB)
		if err != nil {
			log.Error().Err(err).Msg("Faild to add h3 aggregation")
		}
	}
}

func insertQueryToString(query bob.BaseQuery[*dialect.InsertQuery]) string {
	buf := new(bytes.Buffer)
	_, err := query.WriteQuery(context.TODO(), buf, 0)
	if err != nil {
		return fmt.Sprintf("Failed to write query: %v", err)
	}
	return buf.String()
}
