package main

import (
	"bytes"
	"context"
	"embed"
	"errors"
	"fmt"
	"html/template"
	"io"
	"log/slog"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/Gleipnir-Technology/nidus-sync/models"
	"github.com/aarondl/opt/null"
	//"github.com/riverqueue/river/rivershared/util/slogutil"
	"github.com/stephenafamo/bob/dialect/psql/sm"
)

//go:embed templates/*
var embeddedFiles embed.FS

var (
	dashboard                       = newBuiltTemplate("dashboard", "authenticated")
	oauthPrompt                     = newBuiltTemplate("oauth-prompt", "authenticated")
	phoneCall                       = newBuiltTemplate("phone-call", "base")
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
	name     string
	template *template.Template
}

type Link struct {
	Href  string
	Title string
}
type ContentPhoneCall struct {
	DistrictName string
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
	DisplayName   string
	Initials      string
	Notifications []Notification
	Username      string
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

func htmlDashboard(ctx context.Context, w http.ResponseWriter, user *models.User) {
	org, err := user.Organization().One(ctx, PGInstance.BobDB)
	if err != nil {
		respondError(w, "Failed to get org", err, http.StatusInternalServerError)
		return
	}
	var lastSync time.Time
	sync, err := org.FieldseekerSyncs(sm.OrderBy("created")).One(ctx, PGInstance.BobDB)
	if err != nil {
		respondError(w, "Failed to get syncs", err, http.StatusInternalServerError)
		return
	} else {
		lastSync = sync.Created
	}
	inspectionCount, err := org.FSMosquitoinspections().Count(ctx, PGInstance.BobDB)
	if err != nil {
		respondError(w, "Failed to get inspection count", err, http.StatusInternalServerError)
		return
	}
	sourceCount, err := org.FSPointlocations().Count(ctx, PGInstance.BobDB)
	if err != nil {
		respondError(w, "Failed to get source count", err, http.StatusInternalServerError)
		return
	}
	serviceCount, err := org.FSServicerequests().Count(ctx, PGInstance.BobDB)
	if err != nil {
		respondError(w, "Failed to get service count", err, http.StatusInternalServerError)
		return
	}
	recentRequests, err := org.FSServicerequests(sm.OrderBy("creationdate").Desc(), sm.Limit(10)).All(ctx, PGInstance.BobDB)
	if err != nil {
		respondError(w, "Failed to get recent service", err, http.StatusInternalServerError)
		return
	}

	requests := make([]ServiceRequestSummary, 0)
	for _, r := range recentRequests {
		requests = append(requests, ServiceRequestSummary{
			Date:     time.UnixMilli(r.Creationdate.MustGet()),
			Location: r.Reqaddr1.MustGet(),
			Status:   "Completed",
		})
	}
	notifications, err := notificationsForUser(user)
	if err != nil {
		respondError(w, "Failed to get notifications", err, http.StatusInternalServerError)
		return
	}
	data := ContentDashboard{
		CountInspections:     int(inspectionCount),
		CountMosquitoSources: int(sourceCount),
		CountServiceRequests: int(serviceCount),
		LastSync:             lastSync,
		Org:                  org.Name.MustGet(),
		RecentRequests:       requests,
		User: User{
			DisplayName:   user.DisplayName,
			Initials:      extractInitials(user.DisplayName),
			Notifications: notifications,
			Username:      user.Username,
		},
	}
	renderOrError(w, dashboard, data)
}

func htmlOauthPrompt(w http.ResponseWriter, user *models.User) {
	data := ContentDashboard{
		User: User{
			DisplayName: user.DisplayName,
			Initials:    extractInitials(user.DisplayName),
			Username:    user.Username,
		},
	}
	renderOrError(w, oauthPrompt, data)
}

func htmlPhoneCall(w http.ResponseWriter) {
	data := ContentPhoneCall{
		DistrictName: "[District Name]",
	}
	renderOrError(w, phoneCall, data)
}

func htmlReport(w http.ResponseWriter) {
	url := BaseURL + "/report/t78fd3"
	data := ContentReportDiagnostic{
		URL: url,
	}
	renderOrError(w, report, data)
}

func htmlReportConfirmation(w http.ResponseWriter, code string) {
	url := BaseURL + "/report/" + code + "/history"
	data := ContentReportDiagnostic{
		URL: url,
	}
	renderOrError(w, reportConfirmation, data)
}

func htmlReportContribute(w http.ResponseWriter, code string) {
	nextURL := BaseURL + "/report/" + code + "/schedule"
	data := ContentReportDetail{
		NextURL: nextURL,
	}
	renderOrError(w, reportContribute, data)
}

func htmlReportDetail(w http.ResponseWriter, code string) {
	nextURL := BaseURL + "/report/" + code + "/evidence"
	data := ContentReportDetail{
		NextURL:   nextURL,
		UpdateURL: BaseURL + "/report/" + code + "/update",
	}
	renderOrError(w, reportDetail, data)
}

func htmlReportEvidence(w http.ResponseWriter, code string) {
	nextURL := BaseURL + "/report/" + code + "/contribute"
	data := ContentReportDetail{
		NextURL: nextURL,
	}
	renderOrError(w, reportEvidence, data)
}

func htmlReportSchedule(w http.ResponseWriter, code string) {
	nextURL := BaseURL + "/report/" + code + "/confirm"
	data := ContentReportDetail{
		NextURL: nextURL,
	}
	renderOrError(w, reportSchedule, data)
}

func htmlReportUpdate(w http.ResponseWriter, code string) {
	nextURL := BaseURL + "/report/" + code + "/evidence"
	data := ContentReportDetail{
		NextURL: nextURL,
	}
	renderOrError(w, reportUpdate, data)
}

func htmlServiceRequest(w http.ResponseWriter) {
	data := ContentPlaceholder{}
	renderOrError(w, serviceRequest, data)
}

func htmlServiceRequestDetail(w http.ResponseWriter, code string) {
	data := ContentPlaceholder{}
	renderOrError(w, serviceRequestDetail, data)
}

func htmlServiceRequestLocation(w http.ResponseWriter) {
	data := ContentPlaceholder{}
	renderOrError(w, serviceRequestLocation, data)
}

func htmlServiceRequestMosquito(w http.ResponseWriter) {
	data := ContentPlaceholder{}
	renderOrError(w, serviceRequestMosquito, data)
}

func htmlServiceRequestPool(w http.ResponseWriter) {
	data := ContentPlaceholder{}
	renderOrError(w, serviceRequestPool, data)
}

func htmlServiceRequestQuick(w http.ResponseWriter) {
	data := ContentPlaceholder{}
	renderOrError(w, serviceRequestQuick, data)
}

func htmlServiceRequestQuickConfirmation(w http.ResponseWriter) {
	data := ContentPlaceholder{}
	renderOrError(w, serviceRequestQuickConfirmation, data)
}

func htmlServiceRequestUpdates(w http.ResponseWriter) {
	data := ContentPlaceholder{}
	renderOrError(w, serviceRequestUpdates, data)
}

func htmlSignin(w http.ResponseWriter, errorCode string) {
	data := ContentSignin{
		InvalidCredentials: errorCode == "invalid-credentials",
	}
	renderOrError(w, signin, data)
}

func htmlSignup(w http.ResponseWriter, path string) {
	data := ContentSignup{}
	renderOrError(w, signup, data)
}

func makeFuncMap() template.FuncMap {
	funcMap := template.FuncMap{
		"timeElapsed": timeElapsed,
		"timeSince":   timeSince,
	}
	return funcMap
}
func newBuiltTemplate(name string, files ...string) BuiltTemplate {
	files_on_disk := true
	all_files := append([]string{name}, files...)
	for _, f := range all_files {
		full_path := "templates/" + f + ".html"
		_, err := os.Stat(full_path)
		if err != nil {
			files_on_disk = false
			break
		}
	}
	if files_on_disk {
		return BuiltTemplate{
			files:    all_files,
			name:     name,
			template: nil,
		}
	}
	return BuiltTemplate{
		files:    all_files,
		name:     name,
		template: parseEmbedded(all_files),
	}
}

func parseEmbedded(files []string) *template.Template {
	funcMap := makeFuncMap()
	// Remap the file names to embedded paths
	paths := make([]string, 0)
	for _, f := range files {
		paths = append(paths, "templates/"+f+".html")
	}
	for _, f := range components {
		paths = append(paths, "templates/components/"+f+".html")
	}
	name := files[0]
	return template.Must(
		template.New(name).Funcs(funcMap).ParseFS(embeddedFiles, paths...))
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

type Notification struct {
	Link    string
	Message string
	Time    time.Time
	Type    string
}

func notificationsForUser(u *models.User) ([]Notification, error) {
	return []Notification{Notification{
		Link:    "/foo/bar",
		Message: "hey, your oauth is broken.",
		Time:    time.Now(),
		Type:    "alert",
	}}, nil
}

func renderOrError(w http.ResponseWriter, template BuiltTemplate, context interface{}) {
	buf := &bytes.Buffer{}
	err := template.ExecuteTemplate(buf, context)
	if err != nil {
		slog.Error("Failed to render template", slog.String("err", err.Error()), slog.String("template", template.name))
		respondError(w, "Failed to render template", err, http.StatusInternalServerError)
		return
	}
	buf.WriteTo(w)
}
