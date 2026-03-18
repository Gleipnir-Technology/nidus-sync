package text

import (
	"context"

	"fmt"
	"github.com/rs/zerolog/log"

	"github.com/Gleipnir-Technology/bob"
	"github.com/Gleipnir-Technology/bob/dialect/psql"
	"github.com/Gleipnir-Technology/bob/dialect/psql/um"
	"github.com/Gleipnir-Technology/nidus-sync/config"
	"github.com/Gleipnir-Technology/nidus-sync/db"
	//"github.com/Gleipnir-Technology/nidus-sync/db/enums"
	"github.com/Gleipnir-Technology/nidus-sync/db/models"
	"github.com/Gleipnir-Technology/nidus-sync/db/sql"
	"github.com/Gleipnir-Technology/nidus-sync/llm"
	"github.com/Gleipnir-Technology/nidus-sync/platform/types"
	//"github.com/rs/zerolog/log"
)

func SendTextFromLLM(content string) {
	log.Info().Str("content", content).Msg("Pretend I sent a message")
}
func generateNextMessage(ctx context.Context, history []llm.Message, customer_phone types.E164) (llm.Message, error) {
	_handle_report_status := func() (string, error) {
		return "Report: ABCD-1234-5678, District: Delta MVCD, Status: scheduled, Appointment: Wednesday 3:30pm", nil
	}
	_handle_contact_district := func(reason string) {
		log.Warn().Str("reason", reason).Msg("Contacting district")
	}
	_handle_contact_supervisor := func(reason string) {
		log.Warn().Str("reason", reason).Msg("Contacting supervisor")
	}
	return llm.GenerateNextMessage(ctx, history, _handle_report_status, _handle_contact_district, _handle_contact_supervisor)
}
func handleResetConversation(ctx context.Context, txn bob.Executor, src types.E164) error {
	err := wipeLLMMemory(ctx, src)
	sublog := log.With().Str("src", src.PhoneString()).Logger()
	if err != nil {
		return fmt.Errorf("wipe memory: %w")
	}
	content := "LLM memory wiped"
	err = sendTextCommandResponse(ctx, txn, src, content)
	if err != nil {
		return fmt.Errorf("Failed to indicated memory wiped: %w", err)
	}
	sublog.Info().Err(err).Msg("Wiped LLM memory")
	return nil
}

func loadPreviousMessagesForLLM(ctx context.Context, src types.E164) ([]llm.Message, error) {
	messages, err := sql.TextsBySenders(config.PhoneNumberReportStr, src.PhoneString()).All(ctx, db.PGInstance.BobDB)
	results := make([]llm.Message, 0)
	if err != nil {
		return results, fmt.Errorf("Failed to get message history for %s and %s: %w", config.PhoneNumberReportStr, src, err)
	}
	for _, m := range messages {
		if m.IsVisibleToLLM {
			is_from_customer := (m.Source == src.PhoneString())
			results = append(results, llm.Message{
				IsFromCustomer: is_from_customer,
				Content:        m.Content,
			})
		}
	}
	return results, nil
}
func wipeLLMMemory(ctx context.Context, src types.E164) error {
	destination := config.PhoneNumberReportStr
	rows, err := sql.TextsBySenders(destination, src.PhoneString()).All(ctx, db.PGInstance.BobDB)
	if err != nil {
		return fmt.Errorf("Failed to query for texts: %w", err)
	}
	ids := make([]int32, 0)
	for _, r := range rows {
		ids = append(ids, r.ID)
	}
	_, err = models.CommsTextLogs.Update(
		um.Where(
			models.CommsTextLogs.Columns.ID.EQ(psql.Any(ids)),
		),
		um.SetCol("is_visible_to_llm").ToArg(false),
	).Exec(ctx, db.PGInstance.BobDB)
	if err != nil {
		return fmt.Errorf("Failed to update texts: %w", err)
	}

	return nil
}
