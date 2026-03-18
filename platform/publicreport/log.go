package publicreport

import (
	"context"
	"fmt"
	"time"

	"github.com/Gleipnir-Technology/bob"
	"github.com/Gleipnir-Technology/bob/dialect/psql"
	"github.com/Gleipnir-Technology/bob/dialect/psql/sm"
	"github.com/Gleipnir-Technology/nidus-sync/db"
	"github.com/rs/zerolog/log"
	"github.com/stephenafamo/scan"
)

type LogEntry struct {
	Created  time.Time `db:"created" json:"created"`
	ID       int32     `db:"id" json:"-"`
	Message  string    `db:"message" json:"message"`
	ReportID int32     `db:"report_id" json:"-"`
	Type     string    `db:"type_" json:"type"`
	UserID   *int32    `db:"user_id" json:"user_id"`
}

func logEntriesByReportID(ctx context.Context, report_ids []int32) (map[int32][]LogEntry, error) {
	results := make(map[int32][]LogEntry, len(report_ids))
	for _, report_id := range report_ids {
		results[report_id] = make([]LogEntry, 0)
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
	), scan.StructMapper[LogEntry]())
	if err != nil {
		return results, fmt.Errorf("query created: %w", err)
	}
	log.Debug().Int("len(report_ids)", len(report_ids)).Int("len(rows)", len(rows)).Msg("getting log entries")
	for _, row := range rows {
		results[row.ReportID] = append(results[row.ReportID], row)
	}
	return results, nil
}
