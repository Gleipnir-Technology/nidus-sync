package sync

import (
	"embed"
	"fmt"
	"net/http"

	"github.com/Gleipnir-Technology/nidus-sync/html"
	"github.com/rs/zerolog/log"
)

//go:embed template/*
var embeddedFiles embed.FS

var components = [...]string{"header", "icons", "map", "sidebar"}

func buildTemplate(files ...string) *html.BuiltTemplate {
	subdir := "sync"
	full_files := make([]string, 0)
	for _, f := range files {
		full_files = append(full_files, fmt.Sprintf("%s/template/%s.html", subdir, f))
	}
	for _, c := range components {
		full_files = append(full_files, fmt.Sprintf("%s/template/components/%s.html", subdir, c))
	}
	return html.NewBuiltTemplate(embeddedFiles, "sync/", []string{}, full_files...)
}

// Respond with an error that is visible to the user
func respondError(w http.ResponseWriter, m string, e error, s int) {
	log.Warn().Int("status", s).Err(e).Str("user message", m).Msg("Responding with an error from sync pages")
	http.Error(w, m, s)
}
