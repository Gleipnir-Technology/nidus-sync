package email

import (
	"context"
	//"github.com/rs/zerolog/log"
)

func Job(ctx context.Context, email_id int32) error {
	return sendEmailComplete(ctx, email_id)
}
