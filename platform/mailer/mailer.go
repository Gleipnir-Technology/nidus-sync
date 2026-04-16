package mailer

import (
	"context"
	"fmt"

	"github.com/Gleipnir-Technology/bob"
	"github.com/Gleipnir-Technology/nidus-sync/db/models"
	"github.com/rs/zerolog/log"
)

func ComplianceSend(ctx context.Context, txn bob.Executor, row_id int32) error {
	compliance_req, err := models.FindComplianceReportRequest(ctx, txn, row_id)
	if err != nil {
		return fmt.Errorf("find compliance report: %w", err)
	}
	log.Debug().Int32("id", row_id).Str("public_id", compliance_req.PublicID).Msg("working on mailer")
	return nil
}
