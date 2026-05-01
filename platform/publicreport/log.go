package publicreport

import (
	"context"
	"fmt"

	"github.com/Gleipnir-Technology/bob"
	"github.com/Gleipnir-Technology/bob/dialect/psql"
	"github.com/Gleipnir-Technology/bob/dialect/psql/sm"
	"github.com/Gleipnir-Technology/nidus-sync/db"
	"github.com/Gleipnir-Technology/nidus-sync/db/enums"
	"github.com/Gleipnir-Technology/nidus-sync/db/models"
	"github.com/Gleipnir-Technology/nidus-sync/platform/types"
	"github.com/rs/zerolog/log"
	"github.com/stephenafamo/scan"
)

func logEntriesByReportID(ctx context.Context, report_ids []int32, is_public bool) (map[int32][]*types.LogEntry, error) {
	results := make(map[int32][]*types.LogEntry, len(report_ids))
	for _, report_id := range report_ids {
		results[report_id] = make([]*types.LogEntry, 0)
	}

	rows, err := bob.All(ctx, db.PGInstance.BobDB, psql.Select(
		sm.Columns(
			"l.created",
			"l.id",
			"COALESCE(t.content, '') AS message",
			"l.report_id",
			"l.type_",
			"l.user_id",
		),
		sm.From("publicreport.report_log").As("l"),
		sm.LeftJoin("comms.email_log").As("e").OnEQ(
			psql.Quote("l", "email_log_id"),
			psql.Quote("e", "id"),
		),
		sm.LeftJoin("comms.text_log").As("t").OnEQ(
			psql.Quote("l", "text_log_id"),
			psql.Quote("t", "id"),
		),
		sm.Where(psql.Quote("l", "report_id").EQ(psql.Any(report_ids))),
		sm.OrderBy(psql.Quote("l", "created")),
	), scan.StructMapper[types.LogEntry]())
	if err != nil {
		return results, fmt.Errorf("query created: %w", err)
	}
	log.Debug().Int("len(report_ids)", len(report_ids)).Int("len(rows)", len(rows)).Msg("getting log entries")
	for _, row := range rows {
		results[row.ReportID] = append(results[row.ReportID], &row)
	}
	if !is_public {
		logs_from_texts, err := logEntriesFromTexts(ctx, report_ids)
		if err != nil {
			return results, fmt.Errorf("log from texts: %w", err)
		}
		for report_id, logs := range logs_from_texts {
			cur_logs, ok := results[report_id]
			if !ok {
				return results, fmt.Errorf("no text logs for %d", report_id)
			}
			cur_logs = append(cur_logs, logs...)
			results[report_id] = cur_logs
		}
	}
	return results, nil
}

func logEntriesFromTexts(ctx context.Context, report_ids []int32) (map[int32][]*types.LogEntry, error) {
	results := make(map[int32][]*types.LogEntry, len(report_ids))
	for _, report_id := range report_ids {
		results[report_id] = make([]*types.LogEntry, 0)
	}

	type _Row1 struct {
		ReportID      int32  `db:"report_id"`
		ReporterPhone string `db:"reporter_phone"`
	}
	rows, err := bob.All(ctx, db.PGInstance.BobDB, psql.Select(
		sm.Columns(
			"r.id AS report_id",
			"r.reporter_phone AS reporter_phone",
		),
		sm.From("publicreport.report").As("r"),
		sm.Where(psql.Quote("r", "id").EQ(psql.Any(report_ids))),
	), scan.StructMapper[_Row1]())
	if err != nil {
		return results, fmt.Errorf("query reporter_phone: %w", err)
	}

	phone_number_to_report_id := make(map[string]int32, len(rows))
	phone_numbers := make([]string, 0)
	for _, row := range rows {
		if row.ReporterPhone != "" {
			phone_numbers = append(phone_numbers, row.ReporterPhone)
		}
		phone_number_to_report_id[row.ReporterPhone] = row.ReportID
	}
	rows2, err := bob.All(ctx, db.PGInstance.BobDB, psql.Select(
		sm.Columns(
			models.CommsTextLogs.Columns.Content,
			models.CommsTextLogs.Columns.Created,
			models.CommsTextLogs.Columns.Destination,
			models.CommsTextLogs.Columns.ID,
			models.CommsTextLogs.Columns.IsVisibleToLLM,
			models.CommsTextLogs.Columns.IsWelcome,
			models.CommsTextLogs.Columns.Origin,
			models.CommsTextLogs.Columns.Source,
			models.CommsTextLogs.Columns.TwilioSid,
			models.CommsTextLogs.Columns.TwilioStatus,
		),
		sm.From(models.CommsTextLogs.NameAs()),
		sm.Where(
			psql.Or(
				models.CommsTextLogs.Columns.Destination.EQ(psql.Any(phone_numbers)),
				models.CommsTextLogs.Columns.Source.EQ(psql.Any(phone_numbers)),
			),
		),
		sm.OrderBy(
			models.CommsTextLogs.Columns.Created,
		),
	), scan.StructMapper[models.CommsTextLog]())
	if err != nil {
		return results, fmt.Errorf("query text logs: %w", err)
	}

	report_texts, err := models.ReportTexts.Query(
		models.SelectWhere.ReportTexts.ReportID.In(report_ids...),
	).All(ctx, db.PGInstance.BobDB)
	if err != nil {
		return results, fmt.Errorf("query report texts: %w", err)
	}
	report_text_id_to_user_id := make(map[int32]int32, len(report_texts))
	for _, rt := range report_texts {
		report_text_id_to_user_id[rt.TextLogID] = rt.CreatorID
	}
	for _, row := range rows2 {
		// Either the source or destination will be our mapping to our report ID, we just don't
		// know which one it'll be.
		var report_id int32
		var ok bool
		report_id, ok = phone_number_to_report_id[row.Source]
		if !ok {
			report_id, ok = phone_number_to_report_id[row.Destination]
			if !ok {
				return results, fmt.Errorf("can't map %s or %s to a row ID", row.Source, row.Destination)
			}
		}
		logs, ok := results[report_id]
		if !ok {
			return results, fmt.Errorf("Report %d is not in the mapping", report_id)
		}
		var user_id_ptr *int32 = nil
		var user_id int32 = 0
		user_id, ok = report_text_id_to_user_id[row.ID]
		if !ok {
			user_id_ptr = nil
		} else {
			user_id_ptr = &user_id
		}
		type_ := "message-text-outgoing"
		if row.Origin == enums.CommsTextoriginCustomer {
			type_ = "message-text-incoming"
		}
		logs = append(logs, &types.LogEntry{
			Created:  row.Created,
			ID:       row.ID,
			Message:  row.Content,
			ReportID: report_id,
			Type:     type_,
			UserID:   user_id_ptr,
		})
		results[report_id] = logs
	}
	return results, err
}
