package db

import (
	"context"
	"embed"
	"fmt"
	"path/filepath"
	"strings"

	//"github.com/stephenafamo/bob"
	//"github.com/stephenafamo/bob/dialect/psql"
	"github.com/rs/zerolog/log"
	"github.com/stephenafamo/bob"
	"github.com/stephenafamo/bob/dialect/psql"
)

//go:embed prepared_functions/*.sql
var sqlFiles embed.FS

// PrepareStatements reads all embedded SQL files and executes them
// against the provided database connection. This is intended for
// preparing statements that will be used later.
func prepareStatements(ctx context.Context) error {
	// Get a list of all embedded SQL files
	entries, err := sqlFiles.ReadDir("prepared_functions")
	if err != nil {
		return fmt.Errorf("failed to read SQL directory: %w", err)
	}
	log.Info().Int("len", len(entries)).Msg("Reading prepared functions")

	// Process each SQL file
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".sql") {
			log.Info().Str("name", entry.Name()).Msg("Skipping")
			continue
		}

		// Read the SQL file content
		content, err := sqlFiles.ReadFile(filepath.Join("prepared_functions", entry.Name()))
		if err != nil {
			return fmt.Errorf("failed to read SQL file %s: %w", entry.Name(), err)
		}

		// Get the statement name from the filename (without extension)
		statementName := strings.TrimSuffix(filepath.Base(entry.Name()), ".sql")

		// Execute the SQL to prepare the statement
		_, err = PGInstance.BobDB.Exec(string(content))
		if err != nil {
			return fmt.Errorf("failed to prepare statement %s: %w", statementName, err)
		}
		/*
			query := psql.RawQuery(string(content))
			stmt, err := bob.Prepare(ctx, PGInstance.BobDB, query)
			if err != nil {
				return fmt.Errorf("failed to prepare statement %s: %w", statementName, err)
			}
		*/

		log.Info().Str("statement", statementName).Msg("Prepared statement")
	}

	return nil
}
func TestPreparedQuery(ctx context.Context) error {
	query := psql.RawQuery("EXECUTE test_function")
	result, err := bob.Exec(ctx, PGInstance.BobDB, query)
	if err != nil {
		return fmt.Errorf("Failed to exectue test function: %w", err)
	}
	/*insert_id, err := result.LastInsertId()
	if err != nil {
		log.Error().Err(err).Msg("failed insert id")
		return fmt.Errorf("Failed to get insert ID: %w", err)
	}*/
	rows_affected, err := result.RowsAffected()
	if err != nil {
		log.Error().Err(err).Msg("failed rows affected")
		return fmt.Errorf("Failed to get rows affected: %w", err)
	}
	//log.Info().Int64("insert id", insert_id).Int64("rows", rows_affected).Msg("bah")
	log.Info().Int64("rows", rows_affected).Msg("got rows")

	return nil
}
