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
)

//go:embed prepared_functions/*.sql
var sqlFiles embed.FS

// PrepareStatements reads all embedded SQL files and executes them
// against the provided database connection. This is intended for
// preparing statements that will be used later.
func prepareStatements(ctx context.Context) error {
	// Get a list of all embedded SQL files
	entries, err := sqlFiles.ReadDir(".")
	if err != nil {
		return fmt.Errorf("failed to read SQL directory: %w", err)
	}

	// Process each SQL file
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".sql") {
			continue
		}

		// Read the SQL file content
		content, err := sqlFiles.ReadFile(entry.Name())
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
	return nil
}
