package main

import (
	"errors"
	"reflect"

	"github.com/rs/zerolog/log"
)

func LogErrorTypeInfo(err error) {
	if err == nil {
		log.Error().Msg("Error is nil")
		return
	}

	// Log current error type
	errType := reflect.TypeOf(err)
	log.Warn().Err(err).Str("type", errType.String()).Str("pkgPath", errType.PkgPath()).Msg("Error type info")

	// Recursively log wrapped errors
	wrappedErr := errors.Unwrap(err)
	if wrappedErr != nil {
		log.Info().Msg("Contains wrapped error")
		LogErrorTypeInfo(wrappedErr)
	}
}
