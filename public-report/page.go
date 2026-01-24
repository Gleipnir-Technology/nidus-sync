package publicreport

import (
	"embed"
	"fmt"

	"github.com/Gleipnir-Technology/nidus-sync/htmlpage"
)

//go:embed template/*
var embeddedFiles embed.FS

var components = [...]string{"footer", "header", "photo-upload", "photo-upload-header"}
var svgs = [...]string{"check-report", "mosquito", "pond"}

func buildTemplate(files ...string) *htmlpage.BuiltTemplate {
	subdir := "public-report"
	full_files := make([]string, 0)
	for _, f := range files {
		full_files = append(full_files, fmt.Sprintf("%s/template/%s.html", subdir, f))
	}
	for _, c := range components {
		full_files = append(full_files, fmt.Sprintf("%s/template/component/%s.html", subdir, c))
	}
	for _, c := range svgs {
		full_files = append(full_files, fmt.Sprintf("%s/template/svg/%s.svg", subdir, c))
	}
	return htmlpage.NewBuiltTemplate(embeddedFiles, "public-report/", full_files...)
}
