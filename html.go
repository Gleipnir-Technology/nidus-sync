package main

import (
	"context"
	"errors"
	"fmt"
	"html/template"
	"io"
	"log/slog"
	"os"
	"strings"
	"time"

	"github.com/Gleipnir-Technology/nidus-sync/models"
	"github.com/aarondl/opt/null"
	//"github.com/riverqueue/river/rivershared/util/slogutil"
	"github.com/stephenafamo/bob/dialect/psql/sm"
)

var (
	dashboard                       = newBuiltTemplate("dashboard", "authenticated")
	oauthPrompt                     = newBuiltTemplate("oauth-prompt", "authenticated")
	report                          = newBuiltTemplate("report", "base")
	reportConfirmation              = newBuiltTemplate("report-confirmation", "base")
	reportContribute                = newBuiltTemplate("report-contribute", "base")
	reportDetail                    = newBuiltTemplate("report-detail", "base")
	reportEvidence                  = newBuiltTemplate("report-evidence", "base")
	reportSchedule                  = newBuiltTemplate("report-schedule", "base")
	reportUpdate                    = newBuiltTemplate("report-update", "base")
	serviceRequest                  = newBuiltTemplate("service-request", "base")
	serviceRequestDetail            = newBuiltTemplate("service-request-detail", "base")
	serviceRequestLocation          = newBuiltTemplate("service-request-location", "base")
	serviceRequestMosquito          = newBuiltTemplate("service-request-mosquito", "base")
	serviceRequestPool              = newBuiltTemplate("service-request-pool", "base")
	serviceRequestQuick             = newBuiltTemplate("service-request-quick", "base")
	serviceRequestQuickConfirmation = newBuiltTemplate("service-request-quick-confirmation", "base")
	serviceRequestUpdates           = newBuiltTemplate("service-request-updates", "base")
	signin                          = newBuiltTemplate("signin", "base")
	signup                          = newBuiltTemplate("signup", "base")
)
var components = [...]string{"header"}

type BuiltTemplate struct {
	files    []string
	template *template.Template
}

type Link struct {
	Href  string
	Title string
}
type ContentReportDetail struct {
	NextURL   string
	UpdateURL string
}
type ContentReportDiagnostic struct {
	URL string
}
type ContentDashboard struct {
	CountInspections     int
	CountMosquitoSources int
	CountServiceRequests int
	LastSync             time.Time
	Org                  string
	RecentRequests       []ServiceRequestSummary
	User                 User
}
type ContentPlaceholder struct {
}
type ContentSignin struct {
	InvalidCredentials bool
}
type ContentSignup struct{}
type ServiceRequestSummary struct {
	Date     time.Time
	Location string
	Status   string
}
type User struct {
	DisplayName string
	Initials    string
	Username    string
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

func htmlDashboard(ctx context.Context, w io.Writer, user *models.User) error {
	org, err := user.Organization().One(ctx, PGInstance.BobDB)
	if err != nil {
		return fmt.Errorf("Failed to get org: %v", err)
	}
	var lastSync time.Time
	sync, err := org.FieldseekerSyncs(sm.OrderBy("created")).One(ctx, PGInstance.BobDB)
	if err != nil {
		return fmt.Errorf("Failed to get sync: %v", err)
	} else {
		lastSync = sync.Created
	}
	inspectionCount, err := org.FSMosquitoinspections().Count(ctx, PGInstance.BobDB)
	if err != nil {
		return fmt.Errorf("Failed to get inspection count: %v", err)
	}
	sourceCount, err := org.FSPointlocations().Count(ctx, PGInstance.BobDB)
	if err != nil {
		return fmt.Errorf("Failed to get inspection count: %v", err)
	}
	serviceCount, err := org.FSServicerequests().Count(ctx, PGInstance.BobDB)
	if err != nil {
		return fmt.Errorf("Failed to get service count: %v", err)
	}
	recentRequests, err := org.FSServicerequests(sm.OrderBy("creationdate").Desc(), sm.Limit(10)).All(ctx, PGInstance.BobDB)
	if err != nil {
		return fmt.Errorf("Failed to get recent service: %v", err)
	}

	requests := make([]ServiceRequestSummary, 0)
	for _, r := range recentRequests {
		requests = append(requests, ServiceRequestSummary{
			Date:     time.UnixMilli(r.Creationdate.MustGet()),
			Location: r.Reqaddr1.MustGet(),
			Status:   "Completed",
		})
	}
	data := ContentDashboard{
		CountInspections:     int(inspectionCount),
		CountMosquitoSources: int(sourceCount),
		CountServiceRequests: int(serviceCount),
		LastSync:             lastSync,
		Org:                  org.Name.MustGet(),
		RecentRequests:       requests,
		User: User{
			DisplayName: user.DisplayName,
			Initials:    extractInitials(user.DisplayName),
			Username:    user.Username,
		},
	}
	return dashboard.ExecuteTemplate(w, data)
}

func htmlOauthPrompt(w io.Writer, user *models.User) error {
	data := ContentDashboard{
		User: User{
			DisplayName: user.DisplayName,
			Initials:    extractInitials(user.DisplayName),
			Username:    user.Username,
		},
	}
	return oauthPrompt.ExecuteTemplate(w, data)
}

func htmlReport(w io.Writer) error {
	url := BaseURL + "/report/t78fd3"
	data := ContentReportDiagnostic{
		URL: url,
	}
	return report.ExecuteTemplate(w, data)
}

func htmlReportConfirmation(w io.Writer, code string) error {
	url := BaseURL + "/report/" + code + "/history"
	data := ContentReportDiagnostic{
		URL: url,
	}
	return reportConfirmation.ExecuteTemplate(w, data)
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
		NextURL:   nextURL,
		UpdateURL: BaseURL + "/report/" + code + "/update",
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

func htmlReportUpdate(w io.Writer, code string) error {
	nextURL := BaseURL + "/report/" + code + "/evidence"
	data := ContentReportDetail{
		NextURL: nextURL,
	}
	return reportUpdate.ExecuteTemplate(w, data)
}

func htmlServiceRequest(w io.Writer) error {
	data := ContentPlaceholder{}
	return serviceRequest.ExecuteTemplate(w, data)
}

func htmlServiceRequestDetail(w io.Writer, code string) error {
	data := ContentPlaceholder{}
	return serviceRequestDetail.ExecuteTemplate(w, data)
}

func htmlServiceRequestLocation(w io.Writer) error {
	data := ContentPlaceholder{}
	return serviceRequestLocation.ExecuteTemplate(w, data)
}

func htmlServiceRequestMosquito(w io.Writer) error {
	data := ContentPlaceholder{}
	return serviceRequestMosquito.ExecuteTemplate(w, data)
}

func htmlServiceRequestPool(w io.Writer) error {
	data := ContentPlaceholder{}
	return serviceRequestPool.ExecuteTemplate(w, data)
}

func htmlServiceRequestQuick(w io.Writer) error {
	data := ContentPlaceholder{}
	return serviceRequestQuick.ExecuteTemplate(w, data)
}

func htmlServiceRequestQuickConfirmation(w io.Writer) error {
	data := ContentPlaceholder{}
	return serviceRequestQuickConfirmation.ExecuteTemplate(w, data)
}

func htmlServiceRequestUpdates(w io.Writer) error {
	data := ContentPlaceholder{}
	return serviceRequestUpdates.ExecuteTemplate(w, data)
}

func htmlSignin(w io.Writer, errorCode string) error {
	data := ContentSignin{
		InvalidCredentials: errorCode == "invalid-credentials",
	}
	return signin.ExecuteTemplate(w, data)
}

func htmlSignup(w io.Writer, path string) error {
	data := ContentSignup{}
	return signup.ExecuteTemplate(w, data)
}

func makeFuncMap() template.FuncMap {
	funcMap := template.FuncMap{
		"timeElapsed": timeElapsed,
		"timeSince":   timeSince,
	}
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
	//slog.Info("Rendering templates from disk", slog.Any("paths", slogutil.SliceString(paths)))
	templ, err := template.New(name).Funcs(funcMap).ParseFiles(paths...)
	if err != nil {
		return nil, fmt.Errorf("Failed to parse %s: %v", paths, err)
	}
	return templ, nil
}
func timeElapsed(seconds null.Val[float32]) string {
	if !seconds.IsValue() {
		return "none"
	}
	s := int(seconds.MustGet())
	hours := s / 3600
	remainder := s - (hours * 3600)
	minutes := remainder / 60
	remainder = remainder - (minutes * 60)
	if hours > 0 {
		return fmt.Sprintf("%02d:%02d:%02d", hours, minutes, remainder)
	} else if minutes > 0 {
		return fmt.Sprintf("%02d:%02d", minutes, remainder)
	} else {
		return fmt.Sprintf("%d seconds", remainder)
	}
}

func timeSince(t time.Time) string {
	now := time.Now()
	diff := now.Sub(t)

	hours := diff.Hours()
	slog.Info("time since", slog.String("t", t.String()), slog.Float64("hours", hours))
	if hours < 1 {
		minutes := diff.Minutes()
		return fmt.Sprintf("%d minutes ago", int(minutes))
	} else if hours < 24 {
		return fmt.Sprintf("%d hours ago", int(hours))
	} else {
		days := hours / 24
		return fmt.Sprintf("%d days ago", int(days))
	}
}
