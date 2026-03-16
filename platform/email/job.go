package email

import (
	"context"

	"github.com/Gleipnir-Technology/bob"
	//"github.com/rs/zerolog/log"
)

func Job(ctx context.Context, txn bob.Executor, email_id int32) error {
	return sendEmailComplete(ctx, txn, email_id)
}
