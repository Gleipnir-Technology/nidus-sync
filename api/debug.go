package api

import (
	"io"
	"net/http"
	"os"

	"github.com/rs/zerolog/log"
)

func debugSaveRequest(r *http.Request) {
	tmpFile, err := os.CreateTemp("/tmp", "request-*.data")
	if err != nil {
		log.Error().Err(err).Msg("failed to create temp file for debugSaveRequest")
		return
	}
	defer tmpFile.Close()

	_, err = io.Copy(tmpFile, r.Body)
	if err != nil {
		log.Error().Err(err).Msg("failed to copy request body in debugSaveRequest")
		return
	}
	log.Info().Str("filename", tmpFile.Name()).Msg("Saved request body")
}
