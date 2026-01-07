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

type ContextNuisance struct{}
type ContextPool struct{}
type ContextQuick struct{}
type ContextRoot struct{}
type ContextStatus struct{}

var (
	Nuisance = buildTemplate("nuisance", "base")
	Pool = buildTemplate("pool", "base")
	Quick = buildTemplate("quick", "base")
	Root = buildTemplate("root", "base")
	Status = buildTemplate("status", "base")
)

var components = [...]string{"footer"}

func buildTemplate(files ...string) *htmlpage.BuiltTemplate {
	subdir := "htmlpage/public-reports"
	full_files := make([]string, 0)
	for _, f := range files {
		full_files = append(full_files, fmt.Sprintf("%s/template/%s.html", subdir, f))
	}
	for _, c := range components {
		full_files = append(full_files, fmt.Sprintf("%s/template/component/%s.html", subdir, c))
	}
	return htmlpage.NewBuiltTemplate(embeddedFiles, "htmlpage/public-reports/", full_files...)
}
