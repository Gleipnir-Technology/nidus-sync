package main

import (
	"errors"
	"fmt"
	"html/template"
	"io"
	"log/slog"
	"os"
	"strings"

	"github.com/riverqueue/river/rivershared/util/slogutil"
	"github.com/Gleipnir-Technology/nidus-sync/models"
)

var (
	dashboard = newBuiltTemplate("dashboard", "authenticated")
	report = newBuiltTemplate("report", "base")
	reportContribute = newBuiltTemplate("report-contribute", "base")
	reportDetail = newBuiltTemplate("report-detail", "base")
	reportEvidence = newBuiltTemplate("report-evidence", "base")
	reportSchedule = newBuiltTemplate("report-schedule", "base")
	signin      = newBuiltTemplate("signin", "base")
	signup    = newBuiltTemplate("signup", "base")
)
var components = [ ... ]string{"header"}

type BuiltTemplate struct {
	files    []string
	template *template.Template
}

type Link struct {
	Href  string
	Title string
}
type ContentReportDetail struct {
	NextURL string
}
type ContentReportDiagnostic struct {
	URL string
}
type ContentDashboard struct {
	User User
}
type ContentPlaceholder struct {
}
type ContentSignin struct {
	InvalidCredentials bool
}
type ContentSignup struct { }
type User struct {
	DisplayName string
	Initials string
	Username string
}

func (bt *BuiltTemplate) ExecuteTemplate(w io.Writer, data any) error {
	name := bt.files[0] + ".html"
	if bt.template == nil {
		templ, err := parseFromDisk(bt.files)
		if err != nil {
			return fmt.Errorf("Failed to parse template file: %v", err)
		}
		if templ == nil {
			w.Write([]byte("Failed to read from disk: "))
			return errors.New("Template parsing failed")
		}
		return templ.ExecuteTemplate(w, name, data)
	} else {
		return bt.template.ExecuteTemplate(w, name, data)
	}
}

func extractInitials(name string) string {
	parts := strings.Fields(name)
    var initials strings.Builder
    
    for _, part := range parts {
        if len(part) > 0 {
            initials.WriteString(strings.ToUpper(string(part[0])))
        }
    }
    
    return initials.String()
}

func htmlDashboard(w io.Writer, user *models.User) error {
	data := ContentDashboard{
		User:    User{
			DisplayName: user.DisplayName,
			Initials: extractInitials(user.DisplayName),
			Username: user.Username,
		},
	}
	return dashboard.ExecuteTemplate(w, data)
}

func htmlReport(w io.Writer) error {
	url := BaseURL + "/report/t78fd3"
	data := ContentReportDiagnostic{
		URL: url,
	}
	return report.ExecuteTemplate(w, data)
}

func htmlReportContribute(w io.Writer, code string) error {
	nextURL := BaseURL + "/report/" + code + "/schedule"
	data := ContentReportDetail{
		NextURL: nextURL,
	}
	return reportContribute.ExecuteTemplate(w, data)
}

func htmlReportDetail(w io.Writer, code string) error {
	nextURL := BaseURL + "/report/" + code + "/evidence"
	data := ContentReportDetail{
		NextURL: nextURL,
	}
	return reportDetail.ExecuteTemplate(w, data)
}

func htmlReportEvidence(w io.Writer, code string) error {
	nextURL := BaseURL + "/report/" + code + "/contribute"
	data := ContentReportDetail{
		NextURL: nextURL,
	}
	return reportEvidence.ExecuteTemplate(w, data)
}

func htmlReportSchedule(w io.Writer, code string) error {
	nextURL := BaseURL + "/report/" + code + "/confirm"
	data := ContentReportDetail{
		NextURL: nextURL,
	}
	return reportSchedule.ExecuteTemplate(w, data)
}

func htmlSignin(w io.Writer, errorCode string) error {
	data := ContentSignin{
		InvalidCredentials: errorCode == "invalid-credentials",
	}
	return signin.ExecuteTemplate(w, data)
}

func htmlSignup(w io.Writer, path string) error {
	data := ContentSignup{
	}
	return signup.ExecuteTemplate(w, data)
}

func makeFuncMap() template.FuncMap {
	funcMap := template.FuncMap{}
	return funcMap
}
func newBuiltTemplate(files ...string) BuiltTemplate {
	files_on_disk := true
	for _, f := range files {
		full_path := "templates/" + f + ".html"
		_, err := os.Stat(full_path)
		if err != nil {
			files_on_disk = false
			break
		}
	}
	if files_on_disk {
		return BuiltTemplate{
			files:    files,
			template: nil,
		}
	}
	return BuiltTemplate{
		files:    files,
		template: parseEmbedded(files),
	}
}

func parseEmbedded(files []string) *template.Template {
	return nil
}

func parseFromDisk(files []string) (*template.Template, error) {
	funcMap := makeFuncMap()
	paths := make([]string, 0)
	for _, f := range files {
		paths = append(paths, "templates/"+f+".html")
	}
	name := files[0] + ".html"
	for _, f := range components {
		paths = append(paths, "templates/components/"+f+".html")
	}
	slog.Info("Rendering templates from disk", slog.Any("paths", slogutil.SliceString(paths)))
	templ, err := template.New(name).Funcs(funcMap).ParseFiles(paths...)
	if err != nil {
		return nil, fmt.Errorf("Failed to parse %s: %v", paths, err)
	}
	return templ, nil
}
