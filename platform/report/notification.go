package report

import (
	"context"
	"crypto/rand"
	"fmt"
	"math/big"
	"strings"

	"github.com/Gleipnir-Technology/bob"
	"github.com/Gleipnir-Technology/bob/dialect/psql"
	"github.com/Gleipnir-Technology/bob/dialect/psql/sm"
	"github.com/Gleipnir-Technology/bob/dialect/psql/um"
	"github.com/Gleipnir-Technology/nidus-sync/background"
	"github.com/Gleipnir-Technology/nidus-sync/db"
	"github.com/Gleipnir-Technology/nidus-sync/db/models"
	"github.com/Gleipnir-Technology/nidus-sync/db/sql"
	"github.com/Gleipnir-Technology/nidus-sync/platform/text"
	"github.com/rs/zerolog/log"
	"github.com/stephenafamo/scan"
)

type ErrorWithCode struct {
	code    string
	err     error
	message string
}

func (e *ErrorWithCode) Code() string {
	return e.code
}
func (e *ErrorWithCode) Error() string {
	return e.message
}

type SomeReport struct {
	reportID  string
	tableName string
}

func (sr SomeReport) districtID(ctx context.Context) *int32 {
	type _Row struct {
		OrganizationID int32
	}

	from := sm.From("no-such-table")
	switch sr.tableName {
	case "nuisance":
		from = sm.From("publicreport.nuisance")
	case "pool":
		from = sm.From("publicreport.pool")
	default:
		log.Error().Str("table-name", sr.tableName).Msg("Programmer error, non-exhaustive switch statement in SomeReport.districtID")
	}
	row, err := bob.One(ctx, db.PGInstance.BobDB, psql.Select(
		from,
		sm.Columns("organization_id"),
		sm.Where(psql.Quote("public_id").EQ(psql.Arg(sr.reportID))),
	), scan.StructMapper[_Row]())
	if err != nil {
		log.Warn().Err(err).Msg("Failed to query for organization_id")
		return nil
	}
	return &row.OrganizationID
}
func (sr SomeReport) updateReporterEmail(ctx context.Context, email string) *ErrorWithCode {
	table := um.Table("so-such-table")
	switch sr.tableName {
	case "nuisance":
		table = um.Table("publicreport.nuisance")
	case "pool":
		table = um.Table("publicreport.pool")
	default:
		return newErrorWithCode("internal-error", "Programmer error: unrecognized table")
	}
	result, err := psql.Update(
		table,
		um.SetCol("reporter_email").ToArg(email),
		um.Where(psql.Quote("public_id").EQ(psql.Arg(sr.reportID))),
	).Exec(ctx, db.PGInstance.BobDB)
	if err != nil {
		return newErrorWithCode("internal-error", "Failed to update report: %w", err)
	}
	rowcount, err := result.RowsAffected()
	if err != nil {
		return newErrorWithCode("internal-error", "Failed to get rows affected: %w", err)
	}
	if rowcount != 1 {
		log.Warn().Str("report_id", sr.reportID).Msg("updated more than one row, which is a programmer error")
	}
	return nil
}
func (sr SomeReport) updateReporterPhone(ctx context.Context, phone text.E164) *ErrorWithCode {
	table := um.Table("so-such-table")
	switch sr.tableName {
	case "nuisance":
		table = um.Table("publicreport.nuisance")
	case "pool":
		table = um.Table("publicreport.pool")
	default:
		return newErrorWithCode("internal-error", "Programmer error: unrecognized table")
	}
	result, err := psql.Update(
		table,
		um.SetCol("reporter_phone").ToArg(text.PhoneString(phone)),
		um.Where(psql.Quote("public_id").EQ(psql.Arg(sr.reportID))),
	).Exec(ctx, db.PGInstance.BobDB)
	if err != nil {
		return newErrorWithCode("internal-error", "Failed to update report: %w", err)
	}
	rowcount, err := result.RowsAffected()
	if err != nil {
		return newErrorWithCode("internal-error", "Failed to get rows affected: %w", err)
	}
	if rowcount != 1 {
		log.Warn().Str("report_id", sr.reportID).Msg("updated more than one row, which is a programmer error")
	}
	return nil
}

func DistrictForReport(ctx context.Context, report_id string) (*models.Organization, error) {
	some_report, err := findSomeReport(ctx, report_id)
	if err != nil {
		return nil, fmt.Errorf("Failed to find report %s: %w", report_id, err)
	}
	org_id := some_report.districtID(ctx)
	if org_id == nil {
		return nil, nil
	}
	result, e := models.FindOrganization(ctx, db.PGInstance.BobDB, *org_id)
	if e != nil {
		return nil, fmt.Errorf("Failed to load organization %d: %w", org_id, e)
	}
	return result, nil
}

// GenerateReportID creates a 12-character random string using only unambiguous
// capital letters and numbers
func GenerateReportID() (string, error) {
	// Define character set (no O/0, I/l/1, 2/Z to avoid confusion)
	const charset = "ABCDEFGHJKLMNPQRSTUVWXY3456789"
	const length = 12

	var builder strings.Builder
	builder.Grow(length)

	// Use crypto/rand for secure randomness
	for i := 0; i < length; i++ {
		// Generate a random index within our charset
		n, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			return "", fmt.Errorf("failed to generate random number: %w", err)
		}

		// Add the randomly selected character to our ID
		builder.WriteByte(charset[n.Int64()])
	}

	return builder.String(), nil
}

func RegisterNotificationEmail(ctx context.Context, report_id string, email string) *ErrorWithCode {
	some_report, err := findSomeReport(ctx, report_id)
	if err != nil {
		return err
	}
	err = some_report.updateReporterEmail(ctx, email)
	if err != nil {
		return err
	}
	background.ReportSubscriptionConfirmationEmail(email, report_id)
	return nil
}

func RegisterNotificationPhone(ctx context.Context, report_id string, phone text.E164) *ErrorWithCode {
	some_report, err := findSomeReport(ctx, report_id)
	if err != nil {
		return err
	}
	err = some_report.updateReporterPhone(ctx, phone)
	if err != nil {
		return err
	}
	background.ReportSubscriptionConfirmationText(phone, report_id)
	return nil
}

func RegisterSubscriptionEmail(ctx context.Context, email string) *ErrorWithCode {
	log.Warn().Msg("RegisterSubscription not implemented yet")
	return nil
}
func RegisterSubscriptionPhone(ctx context.Context, phone text.E164) *ErrorWithCode {
	log.Warn().Msg("RegisterSubscription not implemented yet")
	return nil
}

func findSomeReport(ctx context.Context, report_id string) (result SomeReport, err *ErrorWithCode) {
	rows, e := sql.PublicreportIDTable(report_id).All(ctx, db.PGInstance.BobDB)
	if e != nil {
		return result, newErrorWithCode("internal-error", "Failed to query report ID table: %w", e)
	}
	switch len(rows) {
	case 0:
		return result, newErrorWithCode("invalid-report-id", "No reports match the provided ID")
	case 1:
		break
	default:
		return result, newErrorWithCode("internal-error", "More than one report with the provided ID, which shouldn't happen")
	}
	row := rows[0]
	result.reportID = report_id
	result.tableName = row.FoundInTables[0]
	return result, nil
}

func newErrorWithCode(code string, format string, args ...any) *ErrorWithCode {
	if len(args) > 0 {
		return &ErrorWithCode{
			err:  fmt.Errorf(format, args...),
			code: code,
		}
	} else {
		return &ErrorWithCode{
			code:    code,
			err:     nil,
			message: format,
		}
	}
}
