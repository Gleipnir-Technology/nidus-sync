package lint

import (
	"context"

	"github.com/rs/zerolog/log"
)

type Errorable = func() error

func LogOnErr(f Errorable, msg string) {
	e := f()
	if e != nil {
		log.Error().Err(e).Msg(msg)
	}
}

type ErrorableCtx = func(context.Context) error

func LogOnErrCtx(f ErrorableCtx, ctx context.Context, msg string) {
	e := f(ctx)
	if e != nil {
		log.Error().Err(e).Msg(msg)
	}
}
func LogOnErrRollback(f ErrorableCtx, ctx context.Context, msg string) {
	e := f(ctx)
	if e != nil {
		// We're fine with rollbacks that are already properly closed
		if e.Error() == "sql: transaction has already been committed or rolled back" {
			return
		}
		log.Error().Err(e).Msg(msg)
	}
}
