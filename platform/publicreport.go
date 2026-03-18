package platform

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/aarondl/opt/omit"
	"github.com/aarondl/opt/omitnull"
	//"github.com/Gleipnir-Technology/bob"
	"github.com/Gleipnir-Technology/nidus-sync/db"
	"github.com/Gleipnir-Technology/nidus-sync/db/enums"
	"github.com/Gleipnir-Technology/nidus-sync/db/models"
	//"github.com/Gleipnir-Technology/nidus-sync/platform/background"
	"github.com/Gleipnir-Technology/nidus-sync/platform/email"
	"github.com/Gleipnir-Technology/nidus-sync/platform/event"
	"github.com/Gleipnir-Technology/nidus-sync/platform/text"
	"github.com/rs/zerolog/log"
)

func PublicreportInvalid(ctx context.Context, user User, report_id string) error {
	report, err := reportFromID(ctx, user, report_id)
	if err != nil {
		return fmt.Errorf("query report existence: %w", err)
	}

	err = report.Update(ctx, db.PGInstance.BobDB, &models.PublicreportReportSetter{
		Reviewed:   omitnull.From(time.Now()),
		ReviewerID: omitnull.From(int32(user.ID)),
		Status:     omit.From(enums.PublicreportReportstatustypeInvalidated),
	})

	log.Info().Int32("id", report.ID).Msg("Report marked as invalid")
	event.Updated(event.TypeRMOReport, user.Organization.ID(), report_id)
	return nil
}

func PublicReportMessageCreate(ctx context.Context, user User, report_id, message string) (message_id *int32, err error) {
	txn, err := db.PGInstance.BobDB.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("create txn: %w", err)
	}
	defer txn.Rollback(ctx)

	report, err := reportFromID(ctx, user, report_id)
	if err != nil {
		return nil, fmt.Errorf("query report existence: %w", err)
	}
	if report.ReporterPhone != "" {
		log.Debug().Str("report_id", report_id).Msg("contacting via phone")
		p, err := text.ParsePhoneNumber(report.ReporterPhone)
		if err != nil {
			return nil, fmt.Errorf("parse phone: %w", err)
		}
		msg_id, err := text.ReportMessage(ctx, txn, int32(user.ID), int32(report.ID), *p, message)
		if err != nil {
			return nil, fmt.Errorf("send text: %w", err)
		}
		txn.Commit(ctx)
		log.Debug().Int32("msg_id", *msg_id).Msg("Created text.ReportMessage")
		return msg_id, nil
	} else if report.ReporterEmail != "" {
		msg_id, err := email.ReportMessage(ctx, int32(user.ID), report_id, report.ReporterEmail, message)
		if err != nil {
			return nil, fmt.Errorf("send email: %w", err)
		}
		txn.Commit(ctx)
		return msg_id, nil
	} else {
		log.Debug().Str("report_id", report_id).Msg("contacting via email")
		return nil, errors.New("no contact methods available")
	}
}
func PublicReportReporterUpdated(ctx context.Context, org_id int32, report_id string, tablename string) {
	event.Updated(event.TypeRMOReport, org_id, report_id)
}
func reportFromID(ctx context.Context, user User, report_id string) (*models.PublicreportReport, error) {
	report, err := models.PublicreportReports.Query(
		models.SelectWhere.PublicreportReports.PublicID.EQ(report_id),
		models.SelectWhere.PublicreportReports.OrganizationID.EQ(user.Organization.ID()),
	).One(ctx, db.PGInstance.BobDB)
	if err != nil {
		return nil, err
	}
	return report, nil
}
