package report

import (
	"context"
	//"crypto/rand"
	//"fmt"
	//"math/big"
	//"strconv"
	//"strings"
	"time"

	"github.com/Gleipnir-Technology/bob"
	"github.com/Gleipnir-Technology/bob/dialect/psql"
	"github.com/Gleipnir-Technology/bob/dialect/psql/sm"
	"github.com/Gleipnir-Technology/bob/dialect/psql/um"
	"github.com/aarondl/opt/omit"
	"github.com/aarondl/opt/omitnull"
	//"github.com/Gleipnir-Technology/nidus-sync/background"
	"github.com/Gleipnir-Technology/nidus-sync/db"
	"github.com/Gleipnir-Technology/nidus-sync/db/models"
	//"github.com/Gleipnir-Technology/nidus-sync/db/sql"
	"github.com/Gleipnir-Technology/nidus-sync/platform/text"
	"github.com/rs/zerolog/log"
	"github.com/stephenafamo/scan"
)

type Nuisance struct {
	id             int32
	publicReportID string
	row            *models.PublicreportNuisance
}

func (sr Nuisance) PublicReportID() string {
	return sr.publicReportID
}
func (sr Nuisance) addNotificationEmail(ctx context.Context, txn bob.Tx, email string) *ErrorWithCode {
	setter := models.PublicreportNotifyEmailNuisanceSetter{
		Created:      omit.From(time.Now()),
		Deleted:      omitnull.FromPtr[time.Time](nil),
		NuisanceID:   omit.From(sr.id),
		EmailAddress: omit.From(email),
	}
	_, err := models.PublicreportNotifyPhoneNuisances.Insert(&setter).Exec(ctx, txn)
	if err != nil {
		return newErrorWithCode("internal-error", "Failed to save new notification row")
	}
	return nil
}
func (sr Nuisance) addNotificationPhone(ctx context.Context, txn bob.Tx, phone text.E164) *ErrorWithCode {
	var err error
	setter := models.PublicreportNotifyPhoneNuisanceSetter{
		Created:    omit.From(time.Now()),
		Deleted:    omitnull.FromPtr[time.Time](nil),
		NuisanceID: omit.From(sr.id),
		PhoneE164:  omit.From(text.PhoneString(phone)),
	}
	_, err = models.PublicreportNotifyPhoneNuisances.Insert(&setter).Exec(ctx, txn)
	if err != nil {
		return newErrorWithCode("internal-error", "Failed to save new notification row")
	}
	return nil
}
func (sr Nuisance) districtID(ctx context.Context) *int32 {
	type _Row struct {
		OrganizationID *int32
	}
	row, err := bob.One(ctx, db.PGInstance.BobDB, psql.Select(
		sm.From("publicreport.nuisance"),
		sm.Columns("organization_id"),
		sm.Where(psql.Quote("public_id").EQ(psql.Arg(sr.publicReportID))),
	), scan.StructMapper[_Row]())
	if err != nil {
		log.Warn().Err(err).Msg("Failed to query for organization_id")
		return nil
	}
	return row.OrganizationID
}
func (sr Nuisance) reportID() int32 {
	return sr.id
}
func (sr Nuisance) updateReporterConsent(ctx context.Context, txn bob.Tx, has_consent bool) *ErrorWithCode {
	setter := models.PublicreportNuisanceSetter{
		ReporterContactConsent: omitnull.From(has_consent),
	}
	_, err := models.PublicreportNotifyPhoneNuisances.Insert(&setter).Exec(ctx, txn)
	if err != nil {
		return newErrorWithCode("internal-error", "Failed to save new notification row")
	}
	return nil
}
func (sr Nuisance) updateReporterEmail(ctx context.Context, txn bob.Tx, email string) *ErrorWithCode {
	result, err := psql.Update(
		um.Table("publicreport.nuisance"),
		um.SetCol("reporter_email").ToArg(email),
		um.Where(psql.Quote("public_id").EQ(psql.Arg(sr.publicReportID))),
	).Exec(ctx, txn)
	if err != nil {
		return newErrorWithCode("internal-error", "Failed to update report: %w", err)
	}
	rowcount, err := result.RowsAffected()
	if err != nil {
		return newErrorWithCode("internal-error", "Failed to get rows affected: %w", err)
	}
	if rowcount != 1 {
		log.Warn().Str("public_report_id", sr.publicReportID).Msg("updated more than one row, which is a programmer error")
	}
	return nil
}
func (sr Nuisance) updateReporterName(ctx context.Context, txn bob.Tx, name string) *ErrorWithCode {
	result, err := psql.Update(
		um.Table("publicreport.nuisance"),
		um.SetCol("reporter_name").ToArg(name),
		um.Where(psql.Quote("public_id").EQ(psql.Arg(sr.publicReportID))),
	).Exec(ctx, txn)
	if err != nil {
		return newErrorWithCode("internal-error", "Failed to update report: %w", err)
	}
	rowcount, err := result.RowsAffected()
	if err != nil {
		return newErrorWithCode("internal-error", "Failed to get rows affected: %w", err)
	}
	if rowcount != 1 {
		log.Warn().Str("public_report_id", sr.publicReportID).Msg("updated more than one row, which is a programmer error")
	}
	return nil
}
func (sr Nuisance) updateReporterPhone(ctx context.Context, txn bob.Tx, phone text.E164) *ErrorWithCode {
	result, err := psql.Update(
		um.Table("publicreport.nuisance"),
		um.SetCol("reporter_phone").ToArg(text.PhoneString(phone)),
		um.Where(psql.Quote("public_id").EQ(psql.Arg(sr.publicReportID))),
	).Exec(ctx, txn)
	if err != nil {
		return newErrorWithCode("internal-error", "Failed to update report: %w", err)
	}
	rowcount, err := result.RowsAffected()
	if err != nil {
		return newErrorWithCode("internal-error", "Failed to get rows affected: %w", err)
	}
	if rowcount != 1 {
		log.Warn().Str("public_report_id", sr.publicReportID).Msg("updated more than one row, which is a programmer error")
	}
	return nil
}

func newNuisance(ctx context.Context, public_id string, report_id int32) (Nuisance, *ErrorWithCode) {
	row, err := models.FindPublicreportNuisance(ctx, db.PGInstance.BobDB, report_id)
	if err != nil {
		return Nuisance{}, newErrorWithCode("internal-error", "Failed to find nuisance report %d: %w", public_id, err)
	}
	return Nuisance{
		id:             report_id,
		publicReportID: public_id,
		row:            row,
	}, nil
}
