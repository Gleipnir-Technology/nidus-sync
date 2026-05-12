package platform

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/Gleipnir-Technology/nidus-sync/db"
	modelpublic "github.com/Gleipnir-Technology/nidus-sync/db/gen/nidus-sync/public/model"
	modelpublicreport "github.com/Gleipnir-Technology/nidus-sync/db/gen/nidus-sync/publicreport/model"
	querycomms "github.com/Gleipnir-Technology/nidus-sync/db/query/comms"
	querypublic "github.com/Gleipnir-Technology/nidus-sync/db/query/public"
	querypublicreport "github.com/Gleipnir-Technology/nidus-sync/db/query/publicreport"
	"github.com/Gleipnir-Technology/nidus-sync/lint"
	"github.com/Gleipnir-Technology/nidus-sync/platform/event"
	"github.com/rs/zerolog/log"
)

type RelatedRecordType int

const (
	RelatedRecordTypeUnknown RelatedRecordType = iota
	RelatedRecordTypeEmail
	RelatedRecordTypeReportCompliance
	RelatedRecordTypeReportNuisance
	RelatedRecordTypeReportWater
	RelatedRecordTypeText
)

func recordTypeFromReportType(t modelpublicreport.Reporttype) RelatedRecordType {
	switch t {
	case modelpublicreport.Reporttype_Compliance:
		return RelatedRecordTypeReportCompliance
	case modelpublicreport.Reporttype_Nuisance:
		return RelatedRecordTypeReportNuisance
	case modelpublicreport.Reporttype_Water:
		return RelatedRecordTypeReportWater
	default:
		return RelatedRecordTypeUnknown
	}
}

type RelatedRecord struct {
	Created time.Time
	ID      string
	Type    RelatedRecordType
}

func CommunicationRelatedRecords(ctx context.Context, user User, comm *modelpublic.Communication) ([]RelatedRecord, error) {
	// Gather associated records
	//  * address
	//  * phone number
	//  * email
	//  * name
	result := make([]RelatedRecord, 0)
	if comm.SourceEmailLogID != nil {
		email_log, err := querycomms.EmailLogFromID(ctx, int64(*comm.SourceEmailLogID))
		if err != nil {
			return result, fmt.Errorf("email log from ID: %w", err)
		}
		email_logs, err := querycomms.EmailLogsFromAddress(ctx, email_log.Source)
		if err != nil {
			return result, fmt.Errorf("email log from ID: %w", err)
		}
		for _, log := range email_logs {
			result = append(result, RelatedRecord{
				Created: log.Created,
				ID:      strconv.Itoa(int(log.ID)),
				Type:    RelatedRecordTypeEmail,
			})
		}
	} else if comm.SourceTextLogID != nil {
		text_log, err := querycomms.TextLogFromID(ctx, int64(*comm.SourceTextLogID))
		if err != nil {
			return result, fmt.Errorf("text log from ID: %w", err)
		}
		text_logs, err := querycomms.EmailLogsFromAddress(ctx, text_log.Source)
		if err != nil {
			return result, fmt.Errorf("text log from ID: %w", err)
		}
		for _, log := range text_logs {
			result = append(result, RelatedRecord{
				Created: log.Created,
				ID:      strconv.Itoa(int(log.ID)),
				Type:    RelatedRecordTypeText,
			})
		}
	} else if comm.SourceReportID != nil {
		report, err := querypublicreport.ReportFromID(ctx, int64(*comm.SourceReportID))
		if err != nil {
			return result, fmt.Errorf("report from ID: %w", err)
		}
		if report.ReporterName != "" {
			reports_by_name, err := querypublicreport.ReportsFromReporterName(ctx, db.PGInstance.PGXPool, int64(user.Organization.ID), report.ReporterName)
			if err != nil {
				return result, fmt.Errorf("reports from reporter name '%s': %w", report.ReporterName, err)
			}
			for _, r := range reports_by_name {
				record_type := recordTypeFromReportType(r.ReportType)
				result = append(result, RelatedRecord{
					Created: r.Created,
					ID:      r.PublicID,
					Type:    record_type,
				})
			}
		}
		if report.AddressID != nil {
			reports_by_address, err := querypublicreport.ReportsFromAddressID(ctx, db.PGInstance.PGXPool, int64(user.Organization.ID), int64(*report.AddressID))
			if err != nil {
				return result, fmt.Errorf("reports from reporter name '%s': %w", report.ReporterName, err)
			}
			for _, r := range reports_by_address {
				record_type := recordTypeFromReportType(r.ReportType)
				result = append(result, RelatedRecord{
					Created: r.Created,
					ID:      r.PublicID,
					Type:    record_type,
				})
			}
		}
	}
	return result, nil
}
func CommunicationsForOrganization(ctx context.Context, org_id int64) ([]modelpublic.Communication, error) {
	return querypublic.CommunicationsFromOrganization(ctx, org_id)
}
func CommunicationFromID(ctx context.Context, user User, comm_id int64) (*modelpublic.Communication, error) {
	comm, err := querypublic.CommunicationFromID(ctx, comm_id)
	if err != nil {
		return nil, err
	}
	if comm.OrganizationID != user.Organization.ID {
		return nil, nil
	}
	return &comm, nil
}
func CommunicationMarkInvalid(ctx context.Context, user User, comm_id int32) error {
	return communicationMark(ctx, user, comm_id, modelpublic.Communicationstatus_Invalid, modelpublic.Communicationlogentry_StatusInvalidated)
}
func CommunicationMarkPendingResponse(ctx context.Context, user User, comm_id int32) error {
	return communicationMark(ctx, user, comm_id, modelpublic.Communicationstatus_Pending, modelpublic.Communicationlogentry_StatusPending)
}
func CommunicationMarkPossibleIssue(ctx context.Context, user User, comm_id int32) error {
	return communicationMark(ctx, user, comm_id, modelpublic.Communicationstatus_PossibleIssue, modelpublic.Communicationlogentry_StatusPossibleIssue)
}
func CommunicationMarkPossibleResolved(ctx context.Context, user User, comm_id int32) error {
	return communicationMark(ctx, user, comm_id, modelpublic.Communicationstatus_PossibleResolved, modelpublic.Communicationlogentry_StatusPossibleResolved)
}

func communicationMark(ctx context.Context, user User, comm_id int32, status modelpublic.Communicationstatus, log_type modelpublic.Communicationlogentry) error {
	txn, err := db.BeginTxn(ctx)
	if err != nil {
		return fmt.Errorf("begin txn: %w", err)
	}
	defer lint.LogOnErrRollback(txn.Rollback, ctx, "rollback")
	err = querypublic.CommunicationSetStatus(ctx, txn, int64(user.Organization.ID), int64(comm_id), status)
	if err != nil {
		return fmt.Errorf("mark: %w", err)
	}
	user_id := int32(user.ID)
	log_entry := modelpublic.CommunicationLogEntry{
		CommunicationID: comm_id,
		Created:         time.Now(),
		Type:            log_type,
		User:            &user_id,
	}
	_, err = querypublic.CommunicationLogEntryInsert(ctx, txn, log_entry)
	if err != nil {
		return fmt.Errorf("insert communication log entry: %w", err)
	}
	if err := txn.Commit(ctx); err != nil {
		return fmt.Errorf("commit: %w", err)
	}
	log.Info().Int32("communication", comm_id).Str("status", status.String()).Msg("Marked communication")

	event.Updated(event.TypeCommunication, user.Organization.ID, strconv.Itoa(int(comm_id)))
	return nil
}
