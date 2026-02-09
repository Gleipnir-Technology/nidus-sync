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
	"github.com/Gleipnir-Technology/nidus-sync/db"
	"github.com/Gleipnir-Technology/nidus-sync/db/models"
	"github.com/Gleipnir-Technology/nidus-sync/platform/text"
	"github.com/aarondl/opt/omit"
	"github.com/aarondl/opt/omitnull"
	"github.com/rs/zerolog/log"
	"github.com/stephenafamo/scan"
)

type Pool struct {
	id             int32
	publicReportID string
	row            *models.PublicreportPool
}

func (sr Pool) PublicReportID() string {
	return sr.publicReportID
}
func (sr Pool) addNotificationEmail(ctx context.Context, txn bob.Tx, email string) *ErrorWithCode {
	setter := models.PublicreportNotifyEmailPoolSetter{
		Created:      omit.From(time.Now()),
		Deleted:      omitnull.FromPtr[time.Time](nil),
		PoolID:       omit.From(sr.id),
		EmailAddress: omit.From(email),
	}
	_, err := models.PublicreportNotifyEmailPools.Insert(&setter).Exec(ctx, txn)
	if err != nil {
		log.Error().Err(err).Msg("Failed to save new notification email row")
		return newInternalError(err, "Failed to save new notification email row")
	}
	return nil
}
func (sr Pool) addNotificationPhone(ctx context.Context, txn bob.Tx, phone text.E164) *ErrorWithCode {
	setter := models.PublicreportNotifyPhonePoolSetter{
		Created:   omit.From(time.Now()),
		Deleted:   omitnull.FromPtr[time.Time](nil),
		PoolID:    omit.From(sr.id),
		PhoneE164: omit.From(text.PhoneString(phone)),
	}
	_, err := models.PublicreportNotifyPhonePools.Insert(&setter).Exec(ctx, txn)
	if err != nil {
		log.Error().Err(err).Msg("Failed to save new notification phone row")
		return newInternalError(err, "Failed to save new notification phone row")
	}
	return nil
}
func (sr Pool) districtID(ctx context.Context) *int32 {
	type _Row struct {
		OrganizationID *int32
	}

	row, err := bob.One(ctx, db.PGInstance.BobDB, psql.Select(
		sm.From("publicreport.pool"),
		sm.Columns("organization_id"),
		sm.Where(psql.Quote("public_id").EQ(psql.Arg(sr.publicReportID))),
	), scan.StructMapper[_Row]())
	if err != nil {
		log.Warn().Err(err).Msg("Failed to query for organization_id")
		return nil
	}
	return row.OrganizationID
}
func (sr Pool) reportID() int32 {
	return sr.id
}
func (sr Pool) updateReporterConsent(ctx context.Context, txn bob.Tx, has_consent bool) *ErrorWithCode {
	return sr.updateReportCol(ctx, txn, &models.PublicreportPoolSetter{
		ReporterContactConsent: omitnull.From(has_consent),
	})
}
func (sr Pool) updateReporterEmail(ctx context.Context, txn bob.Tx, email string) *ErrorWithCode {
	return sr.updateReportCol(ctx, txn, &models.PublicreportPoolSetter{
		ReporterEmail: omit.From(email),
	})
}
func (sr Pool) updateReporterName(ctx context.Context, txn bob.Tx, name string) *ErrorWithCode {
	return sr.updateReportCol(ctx, txn, &models.PublicreportPoolSetter{
		ReporterName: omit.From(name),
	})
}
func (sr Pool) updateReportCol(ctx context.Context, txn bob.Tx, setter *models.PublicreportPoolSetter) *ErrorWithCode {
	err := sr.row.Update(ctx, txn, setter)
	if err != nil {
		log.Error().Err(err).Str("public_id", sr.publicReportID).Int32("report_id", sr.id).Msg("Failed to update report")
		return newInternalError(err, "Failed to update pool report in the database")
	}
	return nil
}
func (sr Pool) updateReporterPhone(ctx context.Context, txn bob.Tx, phone text.E164) *ErrorWithCode {
	return sr.updateReportCol(ctx, txn, &models.PublicreportPoolSetter{
		ReporterPhone: omit.From(text.PhoneString(phone)),
	})
}
func newPool(ctx context.Context, public_id string, report_id int32) (Pool, *ErrorWithCode) {
	row, err := models.FindPublicreportPool(ctx, db.PGInstance.BobDB, report_id)
	if err != nil {
		log.Error().Err(err).Msg("Failed to find pool report")
		return Pool{}, newInternalError(err, "Failed to find pool report %d: %w", public_id, err)
	}
	return Pool{
		id:             report_id,
		publicReportID: public_id,
		row:            row,
	}, nil
}
