package rmo

import (
	"embed"
	"fmt"

	"github.com/Gleipnir-Technology/nidus-sync/html"
)

//go:embed template/*
var embeddedFiles embed.FS

var components = [...]string{"footer", "header-district", "header-rmo", "photo-upload", "photo-upload-header"}

func buildTemplate(files ...string) *html.BuiltTemplate {
	subdir := "rmo"
	full_files := make([]string, 0)
	for _, f := range files {
		full_files = append(full_files, fmt.Sprintf("%s/template/%s.html", subdir, f))
	}
	for _, c := range components {
		full_files = append(full_files, fmt.Sprintf("%s/template/component/%s.html", subdir, c))
	}
	return html.NewBuiltTemplate(embeddedFiles, "rmo/", full_files...)
}
