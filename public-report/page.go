package publicreport

import (
	"embed"
	"fmt"

	"github.com/Gleipnir-Technology/nidus-sync/htmlpage"
)

//go:embed template/*
var embeddedFiles embed.FS

//go:embed static/*
var EmbeddedStaticFS embed.FS

var components = [...]string{"footer", "location-geocode", "location-geocode-header", "photo-upload", "photo-upload-header"}

func buildTemplate(files ...string) *htmlpage.BuiltTemplate {
	subdir := "public-report"
	full_files := make([]string, 0)
	for _, f := range files {
		full_files = append(full_files, fmt.Sprintf("%s/template/%s.html", subdir, f))
	}
	for _, c := range components {
		full_files = append(full_files, fmt.Sprintf("%s/template/component/%s.html", subdir, c))
	}
	return htmlpage.NewBuiltTemplate(embeddedFiles, "public-report/", full_files...)
}
