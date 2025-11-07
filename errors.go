package main

import (
	"errors"
	"log/slog"
	"reflect"
)

func LogErrorTypeInfo(err error) {
	if err == nil {
		slog.Info("Error is nil")
		return
	}

	// Log current error type
	errType := reflect.TypeOf(err)
	slog.Info("Error type info",
		"type", errType.String(),
		"pkgPath", errType.PkgPath(),
		"error", err.Error())

	// Recursively log wrapped errors
	wrappedErr := errors.Unwrap(err)
	if wrappedErr != nil {
		slog.Info("Contains wrapped error")
		LogErrorTypeInfo(wrappedErr)
	}
}
