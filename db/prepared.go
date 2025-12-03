package db

import (
	"context"
	"embed"
	"fmt"
	"path/filepath"
	"strings"
	"time"

	//"github.com/stephenafamo/bob"
	//"github.com/stephenafamo/bob/dialect/psql"
	fslayer "github.com/Gleipnir-Technology/arcgis-go/fieldseeker/layer"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"github.com/stephenafamo/bob"
	"github.com/stephenafamo/bob/dialect/psql"
	"github.com/stephenafamo/scan"
)

//go:embed prepared_functions/*.sql
var sqlFiles embed.FS

// PrepareStatements reads all embedded SQL files and executes them
// against the provided database connection. This is intended for
// preparing statements that will be used later.
func prepareStatements(ctx context.Context) error {
	return nil
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
func TestPreparedQueryOld(ctx context.Context) error {
	type Skn struct {
		Result int
	}
	q := fmt.Sprintf("EXECUTE test_function(%d)", 4)
	query := psql.RawQuery(q)
	result, err := bob.One[Skn](ctx, PGInstance.BobDB, query, scan.StructMapper[Skn]())
	if err != nil {
		return fmt.Errorf("Failed to exectue test function: %w", err)
	}
	log.Info().Int("value", result.Result).Msg("got result")

	return nil
}
func TestPreparedQuery(ctx context.Context, row *fslayer.RodentLocation) error {
	q := queryStoredProcedure("fieldseeker.insert_rodentlocation",
		Uint(row.ObjectID),
		String(row.LocationName),
		String(row.Zone),
		String(row.Zone2),
		Enum(row.Habitat),
		Enum(row.Priority),
		Enum(row.Usetype),
		Enum(row.Active),
		String(row.Description),
		String(row.Accessdesc),
		String(row.Comments),
		Enum(row.Symbology),
		String(row.ExternalID),
		Timestamp(row.Nextactiondatescheduled),
		Int32(row.Locationnumber),
		Timestamp(row.LastInspectionDate),
		String(row.LastInspectionSpecies),
		String(row.LastInspectionAction),
		String(row.LastInspectionConditions),
		String(row.LastInspectionRodentEvidence),
		UUID(row.GlobalID),
		String(row.CreatedUser),
		Timestamp(row.CreatedDate),
		String(row.LastEditedUser),
		Timestamp(row.LastEditedDate),
		Timestamp(row.CreationDate),
		String(row.Creator),
		Timestamp(row.EditDate),
		String(row.Editor),
		String(row.Jurisdiction),
	)
	type InsertResultRow struct {
		Added   bool `db:"row_inserted"`
		Version int  `db:"version_num"`
	}
	type InsertResult struct {
		//Row InsertResultRow `db:"insert_rodentlocation"`
		Row string `db:"insert_rodentlocation"`
	}
	query := psql.RawQuery(q)
	log.Info().Str("query", q).Msg("querying")
	result, err := bob.One[InsertResult](ctx, PGInstance.BobDB, query, scan.StructMapper[InsertResult]())
	if err != nil {
		return fmt.Errorf("Failed to execute test function: %w", err)
	}
	//log.Info().Int("version", result.NextVersion).Msg("got result")
	//log.Info().Bool("added", result.Row.Added).Int("version", result.Row.Version).Msg("done")
	log.Info().Str("row", result.Row).Msg("done")

	return nil
}

// SqlParam is a generic struct that wraps a parameter with its SQL representation
type SqlParam interface {
	ToSql() string
}

// StringParam wraps a string parameter
type StringParam string

func (p StringParam) ToSql() string {
	escapedStr := strings.ReplaceAll(string(p), "'", "''")
	return fmt.Sprintf("'%s'", escapedStr)
}

// IntParam wraps an int parameter
type IntParam int64

func (p IntParam) ToSql() string {
	return fmt.Sprintf("%d", p)
}

// UintParam wraps a uint parameter
type UintParam uint64

func (p UintParam) ToSql() string {
	return fmt.Sprintf("%d", p)
}

type UUIDParam string

func (p UUIDParam) ToSql() string {
	return fmt.Sprintf("'%s'", p)
}

// FloatParam wraps a float parameter
type FloatParam float64

func (p FloatParam) ToSql() string {
	return fmt.Sprintf("%f", p)
}

// BoolParam wraps a boolean parameter
type BoolParam bool

func (p BoolParam) ToSql() string {
	return fmt.Sprintf("%t", p)
}

// EnumParam wraps a PostgreSQL enum parameter
type EnumParam string

func (p EnumParam) ToSql() string {
	escapedStr := strings.ReplaceAll(string(p), "'", "''")
	return fmt.Sprintf("'%s'", escapedStr)
}

// NullParam represents a NULL value
type NullParam struct{}

func (NullParam) ToSql() string {
	return "NULL"
}

// Convenience functions to create typed parameters
func String(s string) StringParam {
	return StringParam(s)
}

type Stringable interface {
	String() string
}

func Enum(e Stringable) EnumParam {
	return EnumParam(e.String())
}
func Int(i int) IntParam {
	return IntParam(i)
}
func Int32(i int32) IntParam {
	return IntParam(i)
}
func Int64(i int64) IntParam {
	return IntParam(i)
}

// Timestamp creates a PostgreSQL TIMESTAMP WITHOUT TIME ZONE parameter
func Timestamp(t time.Time) TimestampParam {
	return TimestampParam(t)
}

// Timestamptz creates a PostgreSQL TIMESTAMP WITH TIME ZONE parameter
func Timestamptz(t time.Time) TimestamptzParam {
	return TimestamptzParam(t)
}

func Uint(u uint) UintParam {
	return UintParam(u)
}
func Uint64(u uint64) UintParam {
	return UintParam(u)
}
func UUID(u uuid.UUID) UUIDParam {
	return UUIDParam(u.String())
}

func Float(f float64) FloatParam {
	return FloatParam(f)
}

func Bool(b bool) BoolParam {
	return BoolParam(b)
}

func Null() NullParam {
	return NullParam{}
}

// TimestampParam wraps a time.Time parameter for PostgreSQL TIMESTAMP WITHOUT TIME ZONE
type TimestampParam time.Time

func (p TimestampParam) ToSql() string {
	// Format as PostgreSQL timestamp without timezone
	// The format string is based on PostgreSQL's expected format
	t := time.Time(p)
	return fmt.Sprintf("'%s'::timestamp", t.Format("2006-01-02 15:04:05.999999"))
}

// TimestamptzParam wraps a time.Time parameter for PostgreSQL TIMESTAMP WITH TIME ZONE
type TimestamptzParam time.Time

func (p TimestamptzParam) ToSql() string {
	// Format as PostgreSQL timestamp with timezone
	t := time.Time(p)
	return fmt.Sprintf("'%s'::timestamptz", t.Format("2006-01-02 15:04:05.999999-07:00"))
}

func queryStoredProcedure(procedure string, params ...SqlParam) string {
	if len(params) == 0 {
		return fmt.Sprintf("SELECT %s()", procedure)
	}

	// Convert each parameter to its SQL representation
	paramStrings := make([]string, len(params))
	for i, param := range params {
		paramStrings[i] = param.ToSql()
	}

	// Join parameters and return the execute statement
	return fmt.Sprintf("SELECT %s(%s)", procedure, strings.Join(paramStrings, ", "))
}

func executeFunction(functionName string, params ...SqlParam) string {
	if len(params) == 0 {
		return fmt.Sprintf("EXECUTE %s()", functionName)
	}

	// Convert each parameter to its SQL representation
	paramStrings := make([]string, len(params))
	for i, param := range params {
		paramStrings[i] = param.ToSql()
	}

	// Join parameters and return the execute statement
	return fmt.Sprintf("EXECUTE %s(%s)", functionName, strings.Join(paramStrings, ", "))
}
