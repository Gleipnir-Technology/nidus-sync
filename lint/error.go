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
