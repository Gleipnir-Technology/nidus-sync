package publicreports

import (
	"embed"
	"fmt"

	"github.com/Gleipnir-Technology/nidus-sync/htmlpage"
)

//go:embed template/*
var embeddedFiles embed.FS

//go:embed static/*
var EmbeddedStaticFS embed.FS

type RootContext struct{}

var (
	Root = buildTemplate("root", "base")
)

var components = [...]string{}

func buildTemplate(files ...string) *htmlpage.BuiltTemplate {
	subdir := "htmlpage/public-reports"
	full_files := make([]string, 0)
	for _, f := range files {
		full_files = append(full_files, fmt.Sprintf("%s/template/%s.html", subdir, f))
	}
	for _, c := range components {
		full_files = append(full_files, fmt.Sprintf("%s/template/components/%s.html", subdir, c))
	}
	return htmlpage.NewBuiltTemplate(embeddedFiles, "htmlpage/public-reports/", full_files...)
}
