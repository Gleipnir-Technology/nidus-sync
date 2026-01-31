package report

import (
	"context"
	"crypto/rand"
	"fmt"
	"math/big"
	"strings"

	"github.com/Gleipnir-Technology/bob/dialect/psql"
	"github.com/Gleipnir-Technology/bob/dialect/psql/um"
	"github.com/Gleipnir-Technology/nidus-sync/background"
	"github.com/Gleipnir-Technology/nidus-sync/db"
	"github.com/Gleipnir-Technology/nidus-sync/db/sql"
	"github.com/nyaruka/phonenumbers"
	"github.com/rs/zerolog/log"
)

type E164 = phonenumbers.PhoneNumber

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
	report_id string
	type_     string
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

func RegisterNotifications(ctx context.Context, report_id string, email string, phone *E164) *ErrorWithCode {
	result, err := psql.Update(
		um.Table("publicreport.quick"),
		um.SetCol("reporter_email").ToArg(email),
		um.SetCol("reporter_phone").ToArg(phone),
		um.Where(psql.Quote("public_id").EQ(psql.Arg(report_id))),
	).Exec(ctx, db.PGInstance.BobDB)
	if err != nil {
		return newErrorWithCode("internal-error", "Failed to update report: %w", err)
	}
	rowcount, err := result.RowsAffected()
	if err != nil {
		return newErrorWithCode("internal-error", "Failed to get rows affected: %w", err)
	}
	if rowcount != 1 {
		log.Warn().Str("report_id", report_id).Msg("updated more than one row, which is a programmer error")
	}
	if email != "" {
		background.ReportSubscriptionConfirmationEmail(email, report_id)
	}
	if phone != nil {
		background.ReportSubscriptionConfirmationText(*phone, report_id)
	}
	return nil
}

func RegisterSubscriptionEmail(ctx context.Context, email string) *ErrorWithCode {
	log.Warn().Msg("RegisterSubscription not implemented yet")
	return nil
}
func RegisterSubscriptionPhone(ctx context.Context, phone *E164) *ErrorWithCode {
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
