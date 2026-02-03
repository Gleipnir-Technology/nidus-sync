package rmo

import (
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog/log"
)

func getScssDebug(w http.ResponseWriter, r *http.Request) {
	path := chi.URLParam(r, "*")
	full_path := "scss/" + path
	//log.Debug().Str("path", path).Str("full_path", full_path).Msg("working on SCSS debug")
	file, err := os.Open(full_path)
	if err != nil {
		respondError(w, "failed to open file", err, http.StatusInternalServerError)
		return
	}
	defer file.Close()
	fileInfo, err := file.Stat()
	if err != nil {
		respondError(w, "failed to stat file", err, http.StatusInternalServerError)
		return
	}
	// Set appropriate headers
	w.Header().Set("Content-Type", "text/scss")
	w.Header().Set("Content-Length", fmt.Sprintf("%d", fileInfo.Size()))
	// Copy file contents to response writer
	_, err = io.Copy(w, file)
	if err != nil {
		// Note: At this point, we've already started writing the response,
		// so we can't change the status code anymore. The best we can do
		// is log the error and abandon the connection.
		log.Warn().Str("path", path).Msg("Failed to write scss file to output")
	}
}
