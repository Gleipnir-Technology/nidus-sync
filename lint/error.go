package lint

import (
	"context"
	"fmt"
	"io"

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
		if e.Error() == "sql: transaction has already been committed or rolled back" || e.Error() == "tx is closed" {
			return
		}
		log.Error().Err(e).Msg(msg)
	}
}

// Fprintf writes a formatted string to w, logging any error.
func Fprintf(w io.Writer, format string, args ...any) {
	_, err := fmt.Fprintf(w, format, args...)
	if err != nil {
		log.Error().Err(err).Msg("fprintf failed")
	}
}

// Write writes p to w, logging any error.
func Write(w io.Writer, p []byte) {
	_, err := w.Write(p)
	if err != nil {
		log.Error().Err(err).Msg("write failed")
	}
}

// Fprint writes a string to w, logging any error.
func Fprint(w io.Writer, a ...any) {
	_, err := fmt.Fprint(w, a...)
	if err != nil {
		log.Error().Err(err).Msg("fprint failed")
	}
}
