package report

import (
	"context"
	"crypto/rand"
	"fmt"
	"math/big"
	"strconv"
	"strings"
	"time"

	"github.com/Gleipnir-Technology/bob"
	"github.com/aarondl/opt/omit"
	"github.com/aarondl/opt/omitnull"
	//"github.com/Gleipnir-Technology/bob/dialect/psql"
	//"github.com/Gleipnir-Technology/bob/dialect/psql/sm"
	//"github.com/Gleipnir-Technology/bob/dialect/psql/um"
	"github.com/Gleipnir-Technology/nidus-sync/background"
	"github.com/Gleipnir-Technology/nidus-sync/db"
	"github.com/Gleipnir-Technology/nidus-sync/db/models"
	"github.com/Gleipnir-Technology/nidus-sync/db/sql"
	"github.com/Gleipnir-Technology/nidus-sync/platform/email"
	"github.com/Gleipnir-Technology/nidus-sync/platform/text"
	"github.com/rs/zerolog/log"
	//"github.com/stephenafamo/scan"
)

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

func RegisterNotificationEmail(ctx context.Context, txn bob.Tx, report_id string, destination string) *ErrorWithCode {
	some_report, err := findSomeReport(ctx, report_id)
	if err != nil {
		return err
	}
	e := email.EnsureInDB(ctx, destination)
	if e != nil {
		return newInternalError(e, "Failed to ensure phone is in DB")
	}
	err = some_report.addNotificationEmail(ctx, txn, destination)
	if err != nil {
		return err
	}
	background.ReportSubscriptionConfirmationEmail(destination, report_id)
	return nil
}

func RegisterNotificationPhone(ctx context.Context, txn bob.Tx, report_id string, phone text.E164) *ErrorWithCode {
	some_report, err := findSomeReport(ctx, report_id)
	if err != nil {
		return err
	}
	e := text.EnsureInDB(ctx, db.PGInstance.BobDB, phone)
	if e != nil {
		return newInternalError(e, "Failed to ensure phone is in DB")
	}
	err = some_report.addNotificationPhone(ctx, txn, phone)
	if err != nil {
		return err
	}
	background.ReportSubscriptionConfirmationText(phone, report_id)
	return nil
}

func RegisterSubscriptionEmail(ctx context.Context, txn bob.Tx, destination string) *ErrorWithCode {
	e := email.EnsureInDB(ctx, destination)
	if e != nil {
		return newInternalError(e, "Failed to ensure email is in DB")
	}
	setter := models.PublicreportSubscribeEmailSetter{
		Created: omit.From(time.Now()),
		Deleted: omitnull.FromPtr[time.Time](nil),
		//DistrictID:   omit.FromPtr[int32](nil),
		EmailAddress: omit.From(destination),
	}
	_, err := models.PublicreportSubscribeEmails.Insert(&setter).Exec(ctx, txn)
	if err != nil {
		log.Error().Err(err).Msg("Failed to save new subscription email row")
		return newInternalError(err, "Failed to save new subscription email row")
	}

	return nil
}
func RegisterSubscriptionPhone(ctx context.Context, txn bob.Tx, phone text.E164) *ErrorWithCode {
	e := text.EnsureInDB(ctx, db.PGInstance.BobDB, phone)
	if e != nil {
		return newInternalError(e, "Failed to ensure phone is in DB")
	}
	setter := models.PublicreportSubscribePhoneSetter{
		Created: omit.From(time.Now()),
		Deleted: omitnull.FromPtr[time.Time](nil),
		//DistrictID:   omitnull.FromPtr[int32](nil),
		PhoneE164: omit.From(text.PhoneString(phone)),
	}
	_, err := models.PublicreportSubscribePhones.Insert(&setter).Exec(ctx, txn)
	if err != nil {
		log.Error().Err(err).Msg("Failed to save new subscription phone row")
		return newInternalError(err, "Failed to save new subscription phone row")
	}
	return nil
}

func SaveReporter(ctx context.Context, txn bob.Tx, report_id string, name string, email string, phone *text.E164, has_consent bool) *ErrorWithCode {
	some_report, err := findSomeReport(ctx, report_id)
	if err != nil {
		return err
	}
	if name != "" {
		err = some_report.updateReporterName(ctx, txn, name)
		if err != nil {
			return err
		}
	}
	if phone != nil {
		err = some_report.updateReporterPhone(ctx, txn, *phone)
		if err != nil {
			return err
		}
	}
	if email != "" {
		err = some_report.updateReporterEmail(ctx, txn, email)
		if err != nil {
			return err
		}
	}
	err = some_report.updateReporterConsent(ctx, txn, has_consent)
	if err != nil {
		return err
	}
	return nil
}
func findSomeReport(ctx context.Context, report_id string) (result SomeReport, err *ErrorWithCode) {
	rows, e := sql.PublicreportIDTable(report_id).All(ctx, db.PGInstance.BobDB)
	if e != nil {
		log.Error().Err(e).Str("report_id", report_id).Msg("failed to query report ID table")
		return result, newErrorWithCode("internal-error", "Failed to query report ID table: %w", e)
	}
	switch len(rows) {
	case 0:
		return result, newErrorWithCode("invalid-report-id", "No reports match the provided ID")
	case 1:
		break
	default:
		log.Error().Err(e).Str("report_id", report_id).Msg("More than one report with the provided ID, which shouldn't happen")
		return result, newErrorWithCode("internal-error", "More than one report with the provided ID, which shouldn't happen")
	}
	row := rows[0]
	report_id_str := row.ReportIds[0]
	t, e := strconv.ParseInt(report_id_str, 10, 32)
	if e != nil {
		log.Error().Err(e).Str("report_id_str", report_id_str).Msg("Unable to parse integer reponse from database")
		return result, newErrorWithCode("internal-error", "Unable to parse integer response from database")
	}

	switch row.FoundInTables[0] {
	case "nuisance":
		return newNuisance(ctx, report_id, int32(t))
	case "pool":
		return newPool(ctx, report_id, int32(t))
	default:
		log.Error().Err(e).Str("table_name", row.FoundInTables[0]).Msg("Unrecognized table")
		return Nuisance{}, newErrorWithCode("internal-error", fmt.Sprintf("Unrecognized table '%s'", row.FoundInTables[0]))
	}
}
