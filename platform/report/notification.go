package report

import (
	"context"
	"crypto/rand"
	"fmt"
	"math/big"
	"strings"
	"time"

	"github.com/Gleipnir-Technology/bob"
	"github.com/Gleipnir-Technology/nidus-sync/db"
	"github.com/Gleipnir-Technology/nidus-sync/db/models"
	"github.com/Gleipnir-Technology/nidus-sync/platform/email"
	"github.com/Gleipnir-Technology/nidus-sync/platform/text"
	"github.com/Gleipnir-Technology/nidus-sync/platform/types"
	"github.com/aarondl/opt/omit"
	"github.com/aarondl/opt/omitnull"
	"github.com/rs/zerolog/log"
)

func DistrictForReport(ctx context.Context, report_id string) (*models.Organization, error) {
	report, err := reportByPublicID(ctx, db.PGInstance.BobDB, report_id)
	if err != nil {
		return nil, fmt.Errorf("Failed to find report %s: %w", report_id, err)
	}
	result, e := models.FindOrganization(ctx, db.PGInstance.BobDB, report.OrganizationID)
	if e != nil {
		return nil, fmt.Errorf("Failed to load organization %d: %w", report.OrganizationID, e)
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

func RegisterNotificationEmail(ctx context.Context, txn bob.Executor, report_id string, destination string) *ErrorWithCode {
	report, e := reportByPublicID(ctx, db.PGInstance.BobDB, report_id)
	if e != nil {
		return newInternalError(e, "Failed to find report")
	}
	e = email.EnsureInDB(ctx, destination)
	if e != nil {
		return newInternalError(e, "Failed to ensure phone is in DB")
	}
	err := addNotificationEmail(ctx, txn, report, destination)
	if err != nil {
		return err
	}
	email.SendReportConfirmation(ctx, destination, report_id)
	return nil
}

func RegisterNotificationPhone(ctx context.Context, txn bob.Executor, report_id string, phone types.E164) *ErrorWithCode {
	report, e := reportByPublicID(ctx, db.PGInstance.BobDB, report_id)
	if e != nil {
		return newInternalError(e, "Failed to find report")
	}
	e = text.EnsureInDB(ctx, db.PGInstance.BobDB, phone)
	if e != nil {
		return newInternalError(e, "Failed to ensure phone is in DB")
	}
	err := addNotificationPhone(ctx, txn, report, phone)
	if err != nil {
		return err
	}
	text.ReportSubscriptionConfirmationText(ctx, db.PGInstance.BobDB, phone, report.PublicID)
	return nil
}

func RegisterSubscriptionEmail(ctx context.Context, txn bob.Executor, destination string) *ErrorWithCode {
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
func RegisterSubscriptionPhone(ctx context.Context, txn bob.Executor, phone types.E164) *ErrorWithCode {
	e := text.EnsureInDB(ctx, db.PGInstance.BobDB, phone)
	if e != nil {
		return newInternalError(e, "Failed to ensure phone is in DB")
	}
	setter := models.PublicreportSubscribePhoneSetter{
		Created: omit.From(time.Now()),
		Deleted: omitnull.FromPtr[time.Time](nil),
		//DistrictID:   omitnull.FromPtr[int32](nil),
		PhoneE164: omit.From(phone.PhoneString()),
	}
	_, err := models.PublicreportSubscribePhones.Insert(&setter).Exec(ctx, txn)
	if err != nil {
		log.Error().Err(err).Msg("Failed to save new subscription phone row")
		return newInternalError(err, "Failed to save new subscription phone row")
	}
	return nil
}

func SaveReporter(ctx context.Context, txn bob.Executor, report_id string, name string, email string, phone *types.E164, has_consent bool) *ErrorWithCode {
	report, e := reportByPublicID(ctx, db.PGInstance.BobDB, report_id)
	if e != nil {
		return newInternalError(e, "Failed to find report")
	}
	if name != "" {
		err := updateReporterName(ctx, txn, report, name)
		if err != nil {
			return err
		}
	}
	if phone != nil {
		err := updateReporterPhone(ctx, txn, report, *phone)
		if err != nil {
			return err
		}
	}
	if email != "" {
		err := updateReporterEmail(ctx, txn, report, email)
		if err != nil {
			return err
		}
	}
	err := updateReporterConsent(ctx, txn, report, has_consent)
	if err != nil {
		return err
	}
	return nil
}
func reportByPublicID(ctx context.Context, txn bob.Executor, public_id string) (*models.PublicreportReport, error) {
	return models.PublicreportReports.Query(
		models.SelectWhere.PublicreportReports.PublicID.EQ(public_id),
	).One(ctx, txn)
}
func addNotificationEmail(ctx context.Context, txn bob.Executor, report *models.PublicreportReport, email string) *ErrorWithCode {
	setter := models.PublicreportNotifyEmailSetter{
		Created:      omit.From(time.Now()),
		Deleted:      omitnull.FromPtr[time.Time](nil),
		EmailAddress: omit.From(email),
		ReportID:     omit.From(report.ID),
	}
	_, err := models.PublicreportNotifyEmails.Insert(&setter).Exec(ctx, txn)
	if err != nil {
		return newInternalError(err, "Failed to save new notification email row")
	}
	return nil
}
func addNotificationPhone(ctx context.Context, txn bob.Executor, report *models.PublicreportReport, phone types.E164) *ErrorWithCode {
	var err error
	setter := models.PublicreportNotifyPhoneSetter{
		Created:   omit.From(time.Now()),
		Deleted:   omitnull.FromPtr[time.Time](nil),
		PhoneE164: omit.From(phone.PhoneString()),
		ReportID:  omit.From(report.ID),
	}
	_, err = models.PublicreportNotifyPhones.Insert(&setter).Exec(ctx, txn)
	if err != nil {
		return newInternalError(err, "Failed to save new notification phone row")
	}
	return nil
}
func updateReporterConsent(ctx context.Context, txn bob.Executor, report *models.PublicreportReport, has_consent bool) *ErrorWithCode {
	return updateReportCol(ctx, txn, report, &models.PublicreportReportSetter{
		ReporterContactConsent: omitnull.From(has_consent),
	})
}
func updateReporterEmail(ctx context.Context, txn bob.Executor, report *models.PublicreportReport, email string) *ErrorWithCode {
	return updateReportCol(ctx, txn, report, &models.PublicreportReportSetter{
		ReporterEmail: omit.From(email),
	})
}
func updateReporterName(ctx context.Context, txn bob.Executor, report *models.PublicreportReport, name string) *ErrorWithCode {
	return updateReportCol(ctx, txn, report, &models.PublicreportReportSetter{
		ReporterName: omit.From(name),
	})
}
func updateReportCol(ctx context.Context, txn bob.Executor, report *models.PublicreportReport, setter *models.PublicreportReportSetter) *ErrorWithCode {
	err := report.Update(ctx, txn, setter)
	if err != nil {
		log.Error().Err(err).Str("public_id", report.PublicID).Int32("report_id", report.ID).Msg("Failed to update report")
		return newInternalError(err, "Failed to update nuisance report in the database")
	}
	return nil
}
func updateReporterPhone(ctx context.Context, txn bob.Executor, report *models.PublicreportReport, phone types.E164) *ErrorWithCode {
	err := report.Update(ctx, txn, &models.PublicreportReportSetter{
		ReporterPhone: omit.From(phone.PhoneString()),
	})
	if err != nil {
		return newInternalError(err, "Failed to update report: %w", err)
	}
	return nil
}
