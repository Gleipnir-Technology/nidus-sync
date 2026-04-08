package platform

import (
	"context"
	"fmt"

	"github.com/Gleipnir-Technology/nidus-sync/db"
	"github.com/Gleipnir-Technology/nidus-sync/db/models"
	"github.com/Gleipnir-Technology/nidus-sync/platform/report"
	"github.com/Gleipnir-Technology/nidus-sync/platform/types"
	//"github.com/rs/zerolog/log"
)

type PublicreportNotification struct {
	Consent      bool
	Email        string
	Name         string
	Notification bool
	Phone        *types.E164
	ReportID     string
	Subscription bool
}

func PublicreportNotificationCreate(ctx context.Context, pn PublicreportNotification) error {
	txn, err := db.PGInstance.BobDB.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin txn: %w", err)
	}
	defer txn.Rollback(ctx)
	rep, err := models.PublicreportReports.Query(
		models.SelectWhere.PublicreportReports.PublicID.EQ(pn.ReportID),
	).One(ctx, db.PGInstance.BobDB)
	if err != nil {
		return fmt.Errorf("find report '%s': %w", pn.ReportID, err)
	}

	err = report.SaveReporter(ctx, txn, pn.ReportID, pn.Name, pn.Email, pn.Phone, pn.Consent)
	if err != nil {
		return fmt.Errorf("save reporter: %w", err)
	}
	if pn.Email != "" {
		if pn.Subscription {
			err = report.RegisterSubscriptionEmail(ctx, txn, pn.Email)
			if err != nil {
				return fmt.Errorf("register subscription email: %w", err)
			}
		}
		if pn.Notification {
			err = report.RegisterNotificationEmail(ctx, txn, pn.ReportID, pn.Email)
			if err != nil {
				return fmt.Errorf("register notification email: %w", err)
			}
		}
	}
	if pn.Phone != nil {
		if pn.Subscription {
			err = report.RegisterSubscriptionPhone(ctx, txn, *pn.Phone)
			if err != nil {
				return fmt.Errorf("register subscription phone: %w", err)
			}
		}
		if pn.Notification {
			err = report.RegisterNotificationPhone(ctx, txn, pn.ReportID, *pn.Phone)
			if err != nil {
				return fmt.Errorf("register notification phone: %w", err)
			}
		}
	}
	txn.Commit(ctx)
	PublicReportReporterUpdated(ctx, rep.OrganizationID, pn.ReportID)
	return nil
}
