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
		Uint("p_objectid", row.ObjectID),
		String("p_locationname", row.LocationName),
		String("p_zone", row.Zone),
		String("p_zone2", row.Zone2),
		String("p_habitat", row.Habitat),
		String("p_priority", row.Priority),
		String("p_usetype", row.Usetype),
		Int16("p_active", row.Active),
		String("p_description", row.Description),
		String("p_accessdesc", row.Accessdesc),
		String("p_comments", row.Comments),
		String("p_symbology", row.Symbology),
		String("p_externalid", row.ExternalID),
		Timestamp("p_nextactiondatescheduled", row.Nextactiondatescheduled),
		Int32("p_locationnumber", row.Locationnumber),
		Timestamp("p_lastinspectdate", row.LastInspectionDate),
		String("p_lastinspectspecies", row.LastInspectionSpecies),
		String("p_lastinspectaction", row.LastInspectionAction),
		String("p_lastinspectconditions", row.LastInspectionConditions),
		String("p_lastinspectrodentevidence", row.LastInspectionRodentEvidence),
		UUID("p_globalid", row.GlobalID),
		String("p_created_user", row.CreatedUser),
		Timestamp("p_created_date", row.CreatedDate),
		String("p_last_edited_user", row.LastEditedUser),
		Timestamp("p_last_edited_date", row.LastEditedDate),
		Timestamp("p_creationdate", row.CreationDate),
		String("p_creator", row.Creator),
		Timestamp("p_editdate", row.EditDate),
		String("p_editor", row.Editor),
		String("p_jurisdiction", row.Jurisdiction),
	)
	query := psql.RawQuery(q)
	log.Info().Str("query", q).Msg("querying")
	result, err := bob.One[InsertResultRow](ctx, PGInstance.BobDB, query, scan.StructMapper[InsertResultRow]())
	if err != nil {
		return fmt.Errorf("Failed to execute test function: %w", err)
	}
	//log.Info().Int("version", result.NextVersion).Msg("got result")
	//log.Info().Bool("added", result.Row.Added).Int("version", result.Row.Version).Msg("done")
	log.Info().Bool("inserted", result.Inserted).Int("version", result.Version).Msg("done")

	return nil
}

// SqlParam is a generic struct that wraps a parameter with its SQL representation
type SqlParam interface {
	ToSql() string
}

// StringParam wraps a string parameter
type StringParam struct {
	Name  string
	Value string
}

func (p StringParam) ToSql() string {
	// Escape quotes since we are writing text directly into the SQL query and this is a key delimiter
	escapedQuotes := strings.ReplaceAll(string(p.Value), "'", "''")
	// Escape question marks because they are a special signal for replacement to bob
	escapedQuestions := strings.ReplaceAll(escapedQuotes, "?", "\\?")
	return fmt.Sprintf("%s => '%s'::varchar", p.Name, escapedQuestions)
}

type Float64Param struct {
	Name  string
	Value float64
}

func (p Float64Param) ToSql() string {
	return fmt.Sprintf("%s => %f::double precision", p.Name, p.Value)
}

// IntParam wraps an int parameter
type Int16Param struct {
	Name  string
	Value int16
}

func (p Int16Param) ToSql() string {
	return fmt.Sprintf("%s => %d::smallint", p.Name, p.Value)
}

type Int32Param struct {
	Name  string
	Value int32
}

func (p Int32Param) ToSql() string {
	return fmt.Sprintf("%s => %d::int", p.Name, p.Value)
}

type Int64Param struct {
	Name  string
	Value int64
}

func (p Int64Param) ToSql() string {
	return fmt.Sprintf("%s => %d::bigint", p.Name, p.Value)
}

// UintParam wraps a uint parameter
type UintParam struct {
	Name  string
	Value uint
}

func (p UintParam) ToSql() string {
	return fmt.Sprintf("%s => %d::int", p.Name, p.Value)
}

type Uint32Param struct {
	Name  string
	Value uint32
}

func (p Uint32Param) ToSql() string {
	return fmt.Sprintf("%s => %d::int", p.Name, p.Value)
}

type Uint64Param struct {
	Name  string
	Value uint64
}

func (p Uint64Param) ToSql() string {
	return fmt.Sprintf("%s => %d::bigint", p.Name, p.Value)
}

type UUIDParam struct {
	Name  string
	Value string
}

func (p UUIDParam) ToSql() string {
	return fmt.Sprintf("%s => '%s'", p.Name, p.Value)
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
func String(n, s string) StringParam {
	return StringParam{
		Name:  n,
		Value: s,
	}
}

type Stringable interface {
	String() string
}

func Enum(n string, e Stringable) EnumParam {
	return EnumParam(e.String())
}
func Float64(n string, f float64) Float64Param {
	return Float64Param{n, f}
}
func Int16(n string, i int16) Int16Param {
	return Int16Param{n, i}
}
func Int32(n string, i int32) Int32Param {
	return Int32Param{n, i}
}
func Int64(n string, i int64) Int64Param {
	return Int64Param{n, i}
}

// Timestamp creates a PostgreSQL TIMESTAMP WITHOUT TIME ZONE parameter
func Timestamp(name string, t time.Time) TimestampParam {
	return TimestampParam{name, t}
}

// Timestamptz creates a PostgreSQL TIMESTAMP WITH TIME ZONE parameter
func Timestamptz(t time.Time) TimestamptzParam {
	return TimestamptzParam(t)
}

func Uint(name string, u uint) UintParam {
	return UintParam{name, u}
}
func Uint32(name string, u uint) Uint32Param {
	return Uint32Param{name, uint32(u)}
}
func Uint64(name string, u uint64) Uint64Param {
	return Uint64Param{name, u}
}
func UUID(name string, u uuid.UUID) UUIDParam {
	return UUIDParam{name, u.String()}
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
type TimestampParam struct {
	Name  string
	Value time.Time
}

func (p TimestampParam) ToSql() string {
	// Format as PostgreSQL timestamp without timezone
	// The format string is based on PostgreSQL's expected format
	t := time.Time(p.Value)
	return fmt.Sprintf("%s => '%s'::timestamp", p.Name, t.Format("2006-01-02 15:04:05.999999"))
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
		return fmt.Sprintf("SELECT * FROM %s()", procedure)
	}

	// Convert each parameter to its SQL representation
	paramStrings := make([]string, len(params))
	for i, param := range params {
		paramStrings[i] = param.ToSql()
	}

	// Join parameters and return the execute statement
	return fmt.Sprintf("SELECT * FROM %s(%s)", procedure, strings.Join(paramStrings, ", "))
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
