package rmo

import (
	"embed"
	"fmt"

	"github.com/Gleipnir-Technology/nidus-sync/html"
)

//go:embed template/*
var embeddedFiles embed.FS

var components = [...]string{"footer", "header", "photo-upload", "photo-upload-header"}
var svgs = [...]string{"check-report", "mosquito", "pond"}

func buildTemplate(files ...string) *html.BuiltTemplate {
	subdir := "rmo"
	full_files := make([]string, 0)
	for _, f := range files {
		full_files = append(full_files, fmt.Sprintf("%s/template/%s.html", subdir, f))
	}
	for _, c := range components {
		full_files = append(full_files, fmt.Sprintf("%s/template/component/%s.html", subdir, c))
	}
	full_svgs := make([]string, 0)
	for _, c := range svgs {
		full_svgs = append(full_svgs, fmt.Sprintf("%s/template/svg/%s.svg", subdir, c))
	}
	return html.NewBuiltTemplate(embeddedFiles, "rmo/", full_svgs, full_files...)
}
