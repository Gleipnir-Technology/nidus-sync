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

type ContextNuisance struct{}
type ContextNuisanceSubmitComplete struct {
	ReportID string
}
type ContextPool struct{
	MapboxToken string
}
type ContextQuick struct{}
type ContextQuickSubmitComplete struct {
	ReportID string
}
type ContextRegisterNotificationsComplete struct {
	ReportID string
}
type ContextRoot struct{}
type ContextStatus struct{}

var (
	Nuisance                      = buildTemplate("nuisance", "base")
	NuisanceSubmitComplete        = buildTemplate("nuisance-submit-complete", "base")
	Pool                          = buildTemplate("pool", "base")
	Quick                         = buildTemplate("quick", "base")
	QuickSubmitComplete           = buildTemplate("quick-submit-complete", "base")
	RegisterNotificationsComplete = buildTemplate("register-notifications-complete", "base")
	Root                          = buildTemplate("root", "base")
	Status                        = buildTemplate("status", "base")
)

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
