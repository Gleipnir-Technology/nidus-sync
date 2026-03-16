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
	"github.com/Gleipnir-Technology/nidus-sync/platform/types"
	"github.com/aarondl/opt/omit"
	"github.com/aarondl/opt/omitnull"
	"github.com/rs/zerolog/log"
	"github.com/stephenafamo/scan"
)

type Water struct {
	id             int32
	publicReportID string
	row            *models.PublicreportWater
}

func (sr Water) PublicReportID() string {
	return sr.publicReportID
}
func (sr Water) addNotificationEmail(ctx context.Context, txn bob.Executor, email string) *ErrorWithCode {
	setter := models.PublicreportNotifyEmailWaterSetter{
		Created:      omit.From(time.Now()),
		Deleted:      omitnull.FromPtr[time.Time](nil),
		EmailAddress: omit.From(email),
		WaterID:      omit.From(sr.id),
	}
	_, err := models.PublicreportNotifyEmailWaters.Insert(&setter).Exec(ctx, txn)
	if err != nil {
		log.Error().Err(err).Msg("Failed to save new notification email row")
		return newInternalError(err, "Failed to save new notification email row")
	}
	return nil
}
func (sr Water) addNotificationPhone(ctx context.Context, txn bob.Executor, phone types.E164) *ErrorWithCode {
	setter := models.PublicreportNotifyPhoneWaterSetter{
		Created:   omit.From(time.Now()),
		Deleted:   omitnull.FromPtr[time.Time](nil),
		PhoneE164: omit.From(phone.PhoneString()),
		WaterID:   omit.From(sr.id),
	}
	_, err := models.PublicreportNotifyPhoneWaters.Insert(&setter).Exec(ctx, txn)
	if err != nil {
		log.Error().Err(err).Msg("Failed to save new notification phone row")
		return newInternalError(err, "Failed to save new notification phone row")
	}
	return nil
}
func (sr Water) districtID(ctx context.Context) *int32 {
	type _Row struct {
		OrganizationID *int32
	}

	row, err := bob.One(ctx, db.PGInstance.BobDB, psql.Select(
		sm.From("publicreport.water"),
		sm.Columns("organization_id"),
		sm.Where(psql.Quote("public_id").EQ(psql.Arg(sr.publicReportID))),
	), scan.StructMapper[_Row]())
	if err != nil {
		log.Warn().Err(err).Msg("Failed to query for organization_id")
		return nil
	}
	return row.OrganizationID
}
func (sr Water) reportID() int32 {
	return sr.id
}
func (sr Water) updateReporterConsent(ctx context.Context, txn bob.Executor, has_consent bool) *ErrorWithCode {
	return sr.updateReportCol(ctx, txn, &models.PublicreportWaterSetter{
		ReporterContactConsent: omitnull.From(has_consent),
	})
}
func (sr Water) updateReporterEmail(ctx context.Context, txn bob.Executor, email string) *ErrorWithCode {
	return sr.updateReportCol(ctx, txn, &models.PublicreportWaterSetter{
		ReporterEmail: omit.From(email),
	})
}
func (sr Water) updateReporterName(ctx context.Context, txn bob.Executor, name string) *ErrorWithCode {
	return sr.updateReportCol(ctx, txn, &models.PublicreportWaterSetter{
		ReporterName: omit.From(name),
	})
}
func (sr Water) updateReportCol(ctx context.Context, txn bob.Executor, setter *models.PublicreportWaterSetter) *ErrorWithCode {
	err := sr.row.Update(ctx, txn, setter)
	if err != nil {
		log.Error().Err(err).Str("public_id", sr.publicReportID).Int32("report_id", sr.id).Msg("Failed to update report")
		return newInternalError(err, "Failed to update water report in the database")
	}
	return nil
}
func (sr Water) updateReporterPhone(ctx context.Context, txn bob.Executor, phone types.E164) *ErrorWithCode {
	return sr.updateReportCol(ctx, txn, &models.PublicreportWaterSetter{
		ReporterPhone: omit.From(phone.PhoneString()),
	})
}
func newWater(ctx context.Context, public_id string, report_id int32) (Water, *ErrorWithCode) {
	row, err := models.FindPublicreportWater(ctx, db.PGInstance.BobDB, report_id)
	if err != nil {
		log.Error().Err(err).Msg("Failed to find water report")
		return Water{}, newInternalError(err, "Failed to find water report %d: %w", public_id, err)
	}
	return Water{
		id:             report_id,
		publicReportID: public_id,
		row:            row,
	}, nil
}
