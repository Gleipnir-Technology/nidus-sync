package platform

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/Gleipnir-Technology/nidus-sync/db"
	"github.com/Gleipnir-Technology/nidus-sync/db/gen/nidus-sync/public/model"
	querycomms "github.com/Gleipnir-Technology/nidus-sync/db/query/comms"
	querypublic "github.com/Gleipnir-Technology/nidus-sync/db/query/public"
	querypublicreport "github.com/Gleipnir-Technology/nidus-sync/db/query/publicreport"
	"github.com/Gleipnir-Technology/nidus-sync/lint"
	"github.com/Gleipnir-Technology/nidus-sync/platform/event"
	"github.com/rs/zerolog/log"
)

type RelatedRecord struct {
	ID   int32
	Type string
}

func CommunicationRelatedRecords(ctx context.Context, user User, comm *model.Communication) ([]RelatedRecord, error) {
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
		log.Debug().Int32("id", email_log.ID).Send()
	}
	if comm.SourceTextLogID != nil {
		text_log, err := querycomms.TextLogFromID(ctx, int64(*comm.SourceTextLogID))
		if err != nil {
			return result, fmt.Errorf("text log from ID: %w", err)
		}
		log.Debug().Int32("id", text_log.ID).Send()
	}
	if comm.SourceReportID != nil {
		report, err := querypublicreport.ReportFromID(ctx, int64(*comm.SourceReportID))
		if err != nil {
			return result, fmt.Errorf("report from ID: %w", err)
		}
		log.Debug().Int32("id", report.ID).Send()
	}
	return result, nil
}
func CommunicationsForOrganization(ctx context.Context, org_id int64) ([]model.Communication, error) {
	return querypublic.CommunicationsFromOrganization(ctx, org_id)
}
func CommunicationFromID(ctx context.Context, user User, comm_id int64) (*model.Communication, error) {
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
	return communicationMark(ctx, user, comm_id, model.Communicationstatus_Invalid, model.Communicationlogentry_StatusInvalidated)
}
func CommunicationMarkPendingResponse(ctx context.Context, user User, comm_id int32) error {
	return communicationMark(ctx, user, comm_id, model.Communicationstatus_Pending, model.Communicationlogentry_StatusPending)
}
func CommunicationMarkPossibleIssue(ctx context.Context, user User, comm_id int32) error {
	return communicationMark(ctx, user, comm_id, model.Communicationstatus_PossibleIssue, model.Communicationlogentry_StatusPossibleIssue)
}
func CommunicationMarkPossibleResolved(ctx context.Context, user User, comm_id int32) error {
	return communicationMark(ctx, user, comm_id, model.Communicationstatus_PossibleResolved, model.Communicationlogentry_StatusPossibleResolved)
}

func communicationMark(ctx context.Context, user User, comm_id int32, status model.Communicationstatus, log_type model.Communicationlogentry) error {
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
	log_entry := model.CommunicationLogEntry{
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
