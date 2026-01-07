package main

import (
	"bytes"
	"context"
	"embed"
	"errors"
	"fmt"
	"html/template"
	"io"
	"math"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/Gleipnir-Technology/nidus-sync/db"
	"github.com/Gleipnir-Technology/nidus-sync/db/models"
	"github.com/Gleipnir-Technology/nidus-sync/db/sql"
	"github.com/aarondl/opt/null"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"github.com/stephenafamo/bob"
	"github.com/stephenafamo/bob/dialect/psql"
	"github.com/stephenafamo/bob/dialect/psql/sm"
	"github.com/uber/h3-go/v4"
)

//go:embed templates/*
var embeddedFiles embed.FS

// Authenticated pages
var (
	cell        = newBuiltTemplate("cell", "authenticated")
	dashboard   = newBuiltTemplate("dashboard", "authenticated")
	oauthPrompt = newBuiltTemplate("oauth-prompt", "authenticated")
	settings    = newBuiltTemplate("settings", "authenticated")
	source      = newBuiltTemplate("source", "authenticated")
)

// Unauthenticated pages
var (
	admin                           = newBuiltTemplate("admin", "base")
	dataEntry                       = newBuiltTemplate("data-entry", "base")
	dataEntryGood                   = newBuiltTemplate("data-entry-good", "base")
	dataEntryBad                    = newBuiltTemplate("data-entry-bad", "base")
	dispatch                        = newBuiltTemplate("dispatch", "base")
	dispatchResults                 = newBuiltTemplate("dispatch-results", "base")
	mockRoot                        = newBuiltTemplate("mock-root", "base")
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
	settingRoot                     = newBuiltTemplate("setting-mock", "base")
	settingIntegration              = newBuiltTemplate("setting-integration", "base")
	settingPesticide                = newBuiltTemplate("setting-pesticide", "base")
	settingPesticideAdd             = newBuiltTemplate("setting-pesticide-add", "base")
	settingUsers                    = newBuiltTemplate("setting-user", "base")
	settingUsersAdd                 = newBuiltTemplate("setting-user-add", "base")
	signin                          = newBuiltTemplate("signin", "base")
	signup                          = newBuiltTemplate("signup", "base")
)
var components = [...]string{"header", "map"}
var templatesByFilename = make(map[string]BuiltTemplate, 0)

type BreedingSourceSummary struct {
	ID            uuid.UUID
	Type          string
	LastInspected *time.Time
	LastTreated   *time.Time
}

type BuiltTemplate struct {
	files    []string
	name     string
	template *template.Template
}

type MapMarker struct {
	LatLng h3.LatLng
}
type ComponentMap struct {
	Center      h3.LatLng
	GeoJSON     interface{}
	MapboxToken string
	Markers     []MapMarker
	Zoom        int
}
type ContentAuthenticatedPlaceholder struct {
	User User
}
type ContentCell struct {
	BreedingSources []BreedingSourceSummary
	CellBoundary    h3.CellBoundary
	Inspections     []Inspection
	MapData         ComponentMap
	Treatments      []Treatment
	User            User
}
type ContentMockURLs struct {
	Dispatch            string
	DispatchResults     string
	ReportConfirmation  string
	ReportDetail        string
	ReportContribute    string
	ReportEvidence      string
	ReportSchedule      string
	ReportUpdate        string
	Root                string
	Setting             string
	SettingIntegration  string
	SettingPesticide    string
	SettingPesticideAdd string
	SettingUser         string
	SettingUserAdd      string
}
type ContentMock struct {
	DistrictName string
	URLs         ContentMockURLs
}
type ContentReportDetail struct {
	NextURL   string
	UpdateURL string
}
type ContentReportDiagnostic struct {
}
type ContentDashboard struct {
	CountInspections     int
	CountMosquitoSources int
	CountServiceRequests int
	Geo                  template.JS
	IsSyncOngoing        bool
	LastSync             *time.Time
	MapData              ComponentMap
	Org                  string
	RecentRequests       []ServiceRequestSummary
	User                 User
}

type ContentDashboardLoading struct {
	User User
}

type ContentPlaceholder struct {
}
type ContentSignin struct {
	InvalidCredentials bool
}
type ContentSignup struct{}
type ContentSource struct {
	Inspections []Inspection
	MapData     ComponentMap
	Source      *BreedingSourceDetail
	Traps       []TrapNearby
	Treatments  []Treatment
	//TreatmentCadence TreatmentCadence
	TreatmentModels []TreatmentModel
	User            User
}
type Inspection struct {
	Action     string
	Date       *time.Time
	Notes      string
	Location   string
	LocationID uuid.UUID
}
type Link struct {
	Href  string
	Title string
}
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
			return fmt.Errorf("Failed to parse template file: %w", err)
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

func bigNumber(n int) string {
	// Convert the number to a string
	numStr := strconv.FormatInt(int64(n), 10)

	// Add commas every three digits from the right
	var result strings.Builder
	for i, char := range numStr {
		if i > 0 && (len(numStr)-i)%3 == 0 {
			result.WriteByte(',')
		}
		result.WriteRune(char)
	}

	return result.String()
}

func contentForUser(ctx context.Context, user *models.User) (User, error) {
	notifications, err := notificationsForUser(ctx, user)
	if err != nil {
		return User{}, err
	}
	return User{
		DisplayName:   user.DisplayName,
		Initials:      extractInitials(user.DisplayName),
		Notifications: notifications,
		Username:      user.Username,
	}, nil

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

func htmlCell(ctx context.Context, w http.ResponseWriter, user *models.User, c int64) {
	org, err := user.Organization().One(ctx, db.PGInstance.BobDB)
	if err != nil {
		respondError(w, "Failed to get org", err, http.StatusInternalServerError)
		return
	}
	userContent, err := contentForUser(ctx, user)
	if err != nil {
		respondError(w, "Failed to get user", err, http.StatusInternalServerError)
		return
	}
	center, err := h3.Cell(c).LatLng()
	if err != nil {
		respondError(w, "Failed to get center", err, http.StatusInternalServerError)
		return
	}
	boundary, err := h3.Cell(c).Boundary()
	if err != nil {
		respondError(w, "Failed to get boundary", err, http.StatusInternalServerError)
		return
	}
	inspections, err := inspectionsByCell(ctx, org, h3.Cell(c))
	if err != nil {
		respondError(w, "Failed to get inspections by cell", err, http.StatusInternalServerError)
		return
	}
	geojson, err := h3ToGeoJSON([]h3.Cell{h3.Cell(c)})
	if err != nil {
		respondError(w, "Failed to get boundaries", err, http.StatusInternalServerError)
		return
	}
	resolution := h3.Cell(c).Resolution()
	sources, err := breedingSourcesByCell(ctx, org, h3.Cell(c))
	if err != nil {
		respondError(w, "Failed to get sources", err, http.StatusInternalServerError)
		return
	}
	treatments, err := treatmentsByCell(ctx, org, h3.Cell(c))
	if err != nil {
		respondError(w, "Failed to get treatments", err, http.StatusInternalServerError)
		return
	}
	data := ContentCell{
		BreedingSources: sources,
		CellBoundary:    boundary,
		Inspections:     inspections,
		MapData: ComponentMap{
			Center: h3.LatLng{
				Lat: center.Lat,
				Lng: center.Lng,
			},
			GeoJSON:     geojson,
			MapboxToken: MapboxToken,
			Zoom:        resolution + 5,
		},
		Treatments: treatments,
		User:       userContent,
	}
	renderOrError(w, cell, &data)
}

func htmlDashboard(ctx context.Context, w http.ResponseWriter, user *models.User) {
	org, err := user.Organization().One(ctx, db.PGInstance.BobDB)
	if err != nil {
		respondError(w, "Failed to get org", err, http.StatusInternalServerError)
		return
	}
	var lastSync *time.Time
	sync, err := org.FieldseekerSyncs(sm.OrderBy("created").Desc()).One(ctx, db.PGInstance.BobDB)
	if err != nil {
		if err.Error() != "sql: no rows in result set" {
			respondError(w, "Failed to get syncs", err, http.StatusInternalServerError)
			return
		}
	} else {
		lastSync = &sync.Created
	}
	is_syncing := isSyncOngoing(org.ID)
	inspectionCount, err := org.Mosquitoinspections().Count(ctx, db.PGInstance.BobDB)
	if err != nil {
		respondError(w, "Failed to get inspection count", err, http.StatusInternalServerError)
		return
	}
	sourceCount, err := org.Pointlocations().Count(ctx, db.PGInstance.BobDB)
	if err != nil {
		respondError(w, "Failed to get source count", err, http.StatusInternalServerError)
		return
	}
	serviceCount, err := org.Servicerequests().Count(ctx, db.PGInstance.BobDB)
	if err != nil {
		respondError(w, "Failed to get service count", err, http.StatusInternalServerError)
		return
	}
	recentRequests, err := org.Servicerequests(sm.OrderBy("creationdate").Desc(), sm.Limit(10)).All(ctx, db.PGInstance.BobDB)
	if err != nil {
		respondError(w, "Failed to get recent service", err, http.StatusInternalServerError)
		return
	}

	requests := make([]ServiceRequestSummary, 0)
	for _, r := range recentRequests {
		requests = append(requests, ServiceRequestSummary{
			Date:     r.Creationdate.MustGet(),
			Location: r.Reqaddr1.MustGet(),
			Status:   "Completed",
		})
	}
	userContent, err := contentForUser(ctx, user)
	if err != nil {
		respondError(w, "Failed to get user context", err, http.StatusInternalServerError)
		return
	}
	data := ContentDashboard{
		CountInspections:     int(inspectionCount),
		CountMosquitoSources: int(sourceCount),
		CountServiceRequests: int(serviceCount),
		IsSyncOngoing:        is_syncing,
		LastSync:             lastSync,
		MapData: ComponentMap{
			MapboxToken: MapboxToken,
		},
		Org:            org.Name.MustGet(),
		RecentRequests: requests,
		User:           userContent,
	}
	renderOrError(w, dashboard, data)
}

func htmlMock(t string, w http.ResponseWriter, code string) {
	data := ContentMock{
		DistrictName: "Delta MVCD",
		URLs: ContentMockURLs{
			Dispatch:            "/mock/dispatch",
			DispatchResults:     "/mock/dispatch-results",
			ReportConfirmation:  fmt.Sprintf("/mock/report/%s/confirm", code),
			ReportDetail:        fmt.Sprintf("/mock/report/%s", code),
			ReportContribute:    fmt.Sprintf("/mock/report/%s/contribute", code),
			ReportEvidence:      fmt.Sprintf("/mock/report/%s/evidence", code),
			ReportSchedule:      fmt.Sprintf("/mock/report/%s/schedule", code),
			ReportUpdate:        fmt.Sprintf("/mock/report/%s/update", code),
			Root:                "/mock",
			Setting:             "/mock/setting",
			SettingIntegration:  "/mock/setting/integration",
			SettingPesticide:    "/mock/setting/pesticide",
			SettingPesticideAdd: "/mock/setting/pesticide/add",
			SettingUser:         "/mock/setting/user",
			SettingUserAdd:      "/mock/setting/user/add",
		},
	}
	template, ok := templatesByFilename[t]
	if !ok {
		log.Error().Str("template", t).Msg("Failed to find template")
		respondError(w, "Failed to render template", nil, http.StatusInternalServerError)
		return
	}
	renderOrError(w, &template, data)
}

func htmlOauthPrompt(w http.ResponseWriter, user *models.User) {
	dp := user.DisplayName
	data := ContentDashboard{
		User: User{
			DisplayName: dp,
			Initials:    extractInitials(dp),
			Username:    user.Username,
		},
	}
	renderOrError(w, oauthPrompt, data)
}

func htmlSettings(w http.ResponseWriter, r *http.Request, user *models.User) {
	userContent, err := contentForUser(r.Context(), user)
	if err != nil {
		respondError(w, "Failed to get user content", err, http.StatusInternalServerError)
		return
	}
	data := ContentAuthenticatedPlaceholder{
		User: userContent,
	}
	renderOrError(w, settings, data)
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

func htmlSource(w http.ResponseWriter, r *http.Request, user *models.User, id uuid.UUID) {
	org, err := user.Organization().One(r.Context(), db.PGInstance.BobDB)
	if err != nil {
		respondError(w, "Failed to get org", err, http.StatusInternalServerError)
		return
	}
	userContent, err := contentForUser(r.Context(), user)
	if err != nil {
		respondError(w, "Failed to get user content", err, http.StatusInternalServerError)
		return
	}
	s, err := sourceByGlobalId(r.Context(), org, id)
	if err != nil {
		respondError(w, "Failed to get source", err, http.StatusInternalServerError)
		return
	}
	inspections, err := inspectionsBySource(r.Context(), org, id)
	if err != nil {
		respondError(w, "Failed to get inspections", err, http.StatusInternalServerError)
		return
	}
	traps, err := trapsBySource(r.Context(), org, id)
	if err != nil {
		respondError(w, "Failed to get traps", err, http.StatusInternalServerError)
		return
	}

	treatments, err := treatmentsBySource(r.Context(), org, id)
	if err != nil {
		respondError(w, "Failed to get treatments", err, http.StatusInternalServerError)
		return
	}
	treatment_models := modelTreatment(treatments)
	latlng, err := s.H3Cell.LatLng()
	if err != nil {
		respondError(w, "Failed to get latlng", err, http.StatusInternalServerError)
		return
	}
	data := ContentSource{
		Inspections: inspections,
		MapData: ComponentMap{
			Center: latlng,
			//GeoJSON:
			MapboxToken: MapboxToken,
			Markers: []MapMarker{
				MapMarker{
					LatLng: latlng,
				},
			},
			Zoom: 13,
		},
		Source:          s,
		Traps:           traps,
		Treatments:      treatments,
		TreatmentModels: treatment_models,
		User:            userContent,
	}

	renderOrError(w, source, data)
}

func gisStatement(cb h3.CellBoundary) string {
	var content strings.Builder
	for i, p := range cb {
		if i != 0 {
			content.WriteString(", ")
		}
		content.WriteString(fmt.Sprintf("%f %f", p.Lng, p.Lat))
	}
	// Repeat the first coordinate to close the polygon
	content.WriteString(fmt.Sprintf(", %f %f", cb[0].Lng, cb[0].Lat))
	return fmt.Sprintf("ST_GeomFromText('POLYGON((%s))', 3857)", content.String())
}

func latLngDisplay(ll h3.LatLng) string {
	latDir := "N"
	latVal := ll.Lat
	if ll.Lat < 0 {
		latDir = "S"
		latVal = -ll.Lat
	}

	lngDir := "E"
	lngVal := ll.Lng
	if ll.Lng < 0 {
		lngDir = "W"
		lngVal = -ll.Lng
	}

	return fmt.Sprintf("%.4f° %s, %.4f° %s", latVal, latDir, lngVal, lngDir)
}

func makeFuncMap() template.FuncMap {
	funcMap := template.FuncMap{
		"bigNumber":          bigNumber,
		"GISStatement":       gisStatement,
		"latLngDisplay":      latLngDisplay,
		"timeAsRelativeDate": timeAsRelativeDate,
		"timeDelta":          timeDelta,
		"timeElapsed":        timeElapsed,
		"timeInterval":       timeInterval,
		"timeSince":          timeSince,
		"uuidShort":          uuidShort,
	}
	return funcMap
}
func newBuiltTemplate(name string, files ...string) *BuiltTemplate {
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
	var result BuiltTemplate
	if files_on_disk {
		result = BuiltTemplate{
			files:    all_files,
			name:     name,
			template: nil,
		}
	} else {
		result = BuiltTemplate{
			files:    all_files,
			name:     name,
			template: parseEmbedded(all_files),
		}
	}
	templatesByFilename[name] = result
	return &result
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
	templ, err := template.New(name).Funcs(funcMap).ParseFiles(paths...)
	if err != nil {
		return nil, fmt.Errorf("Failed to parse %s: %w", paths, err)
	}
	return templ, nil
}

func timeAsRelativeDate(d time.Time) string {
	return d.Format("01-02")
}

// FormatTimeDuration returns a human-readable string representing a time.Duration
// as "X units early" or "X units late"
func timeDelta(d time.Duration) string {
	suffix := "late"
	if d < 0 {
		suffix = "early"
		d = -d // Make duration positive for calculations
	}

	const (
		day  = 24 * time.Hour
		week = 7 * day
	)

	log.Info().Int64("delta", int64(d)).Str("suffix", suffix).Msg("Time delta")
	switch {
	case d >= week:
		weeks := d / week
		if weeks == 1 {
			return "1 week " + suffix
		}
		return fmt.Sprintf("%d weeks %s", weeks, suffix)

	case d >= day:
		days := d / day
		if days == 1 {
			return "1 day " + suffix
		}
		return fmt.Sprintf("%d days %s", days, suffix)

	case d >= time.Hour:
		hours := d / time.Hour
		if hours == 1 {
			return "1 hour " + suffix
		}
		return fmt.Sprintf("%d hours %s", hours, suffix)

	case d >= time.Minute:
		minutes := d / time.Minute
		if minutes == 1 {
			return "1 minute " + suffix
		}
		return fmt.Sprintf("%d minutes %s", minutes, suffix)

	default:
		seconds := d / time.Second
		if seconds == 1 {
			return "1 second " + suffix
		}
		return fmt.Sprintf("%d seconds %s", seconds, suffix)
	}
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

func timeInterval(d time.Duration) string {
	seconds := d.Seconds()

	// Less than 120 seconds -> show in seconds
	if seconds < 120 {
		return fmt.Sprintf("every %d seconds", int(math.Round(seconds)))
	}

	minutes := d.Minutes()
	// Less than 120 minutes -> show in minutes
	if minutes < 120 {
		return fmt.Sprintf("every %d minutes", int(math.Round(minutes)))
	}

	hours := d.Hours()
	// Less than 48 hours -> show in hours
	if hours < 48 {
		return fmt.Sprintf("every %d hours", int(math.Round(hours)))
	}

	days := hours / 24
	// Less than 14 days -> show in days
	if days < 14 {
		return fmt.Sprintf("every %d days", int(math.Round(days)))
	}

	weeks := days / 7
	// Less than 8 weeks -> show in weeks
	if weeks < 8 {
		return fmt.Sprintf("every %d weeks", int(math.Round(weeks)))
	}

	months := days / 30
	// Less than 24 months -> show in months
	if months < 24 {
		return fmt.Sprintf("every %d months", int(math.Round(months)))
	}

	years := days / 365
	return fmt.Sprintf("every %d years", int(math.Round(years)))
}
func timeSince(t *time.Time) string {
	if t == nil {
		return "never"
	}
	now := time.Now()
	diff := now.Sub(*t)

	hours := diff.Hours()
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

func trapsBySource(ctx context.Context, org *models.Organization, sourceID uuid.UUID) ([]TrapNearby, error) {
	locations, err := sql.TrapLocationBySourceID(org.ID, sourceID).All(ctx, db.PGInstance.BobDB)
	if err != nil {
		return nil, fmt.Errorf("Failed to query rows: %w", err)
	}

	location_ids := make([]uuid.UUID, 0)
	var args []bob.Expression
	for _, location := range locations {
		location_ids = append(location_ids, location.TrapLocationGlobalid)
		args = append(args, psql.Arg(location.TrapLocationGlobalid))
	}
	/*
		trap_data, err := org.FSTrapdata(
			sm.Where(
				models.FSTrapdata.Columns.LocID.In(args...),
			),
			sm.OrderBy("enddatetime"),
		).All(ctx, db.PGInstance.BobDB)
	*/

	/*
		query := org.FSTrapdata(
			sm.From(
				psql.Select(
					sm.From(psql.F("ROW_NUMBER")(
						fm.Over(
							wm.PartitionBy(models.FSTrapdata.Columns.LocID),
							wm.OrderBy(models.FSTrapdata.Columns.Enddatetime).Desc(),
						),
					)).As("row_num"),
				sm.Where(models.FSTrapdata.Columns.LocID.In(args...))),
			),
			sm.Where(psql.Quote("row_num").LTE(psql.Arg(10))),
			sm.OrderBy(models.FSTrapdata.Columns.LocID),
			sm.OrderBy(models.FSTrapdata.Columns.Enddatetime).Desc(),
		)
	*/
	/*
		query := psql.Select(
			sm.From(
				psql.Select(
					sm.Columns(
						models.FSTrapdata.Columns.Globalid,
						psql.F("ROW_NUMBER")(
						fm.Over(
							wm.PartitionBy(models.FSTrapdata.Columns.LocID),
							wm.OrderBy(models.FSTrapdata.Columns.Enddatetime).Desc(),
						),
					).As("row_num"),
					sm.From(models.FSTrapdata.Name()),
				),
				sm.Where(models.FSTrapdata.Columns.LocID.In(args...))),
			),
			sm.Where(psql.Quote("row_num").LTE(psql.Arg(10))),
			sm.OrderBy(models.FSTrapdata.Columns.LocID),
			sm.OrderBy(models.FSTrapdata.Columns.Enddatetime).Desc(),
		)
		log.Info().Str("trapdata", queryToString(query)).Msg("Getting trap data")
		trap_data, err := query.Exec(ctx, db.PGInstance.BobDB)
	*/

	trap_data, err := sql.TrapDataByLocationIDRecent(org.ID, location_ids).All(ctx, db.PGInstance.BobDB)
	if err != nil {
		return nil, fmt.Errorf("Failed to query trap data: %w", err)
	}

	counts, err := sql.TrapCountByLocationID(org.ID, location_ids).All(ctx, db.PGInstance.BobDB)
	if err != nil {
		return nil, fmt.Errorf("Failed to query trap counts: %w", err)
	}

	traps, err := toTemplateTraps(locations, trap_data, counts)
	if err != nil {
		return nil, fmt.Errorf("Failed to convert trap data: %w", err)
	}
	return traps, nil
}

func renderOrError(w http.ResponseWriter, template *BuiltTemplate, context interface{}) {
	buf := &bytes.Buffer{}
	err := template.ExecuteTemplate(buf, context)
	if err != nil {
		log.Error().Err(err).Str("template", template.name).Msg("Failed to render template")
		respondError(w, "Failed to render template", err, http.StatusInternalServerError)
		return
	}
	buf.WriteTo(w)
}

func treatmentsBySource(ctx context.Context, org *models.Organization, sourceID uuid.UUID) ([]Treatment, error) {
	var results []Treatment
	rows, err := org.Treatments(
		sm.Where(
			models.FieldseekerTreatments.Columns.Pointlocid.EQ(psql.Arg(sourceID)),
		),
		sm.OrderBy("enddatetime").Desc(),
	).All(ctx, db.PGInstance.BobDB)
	if err != nil {
		return results, fmt.Errorf("Failed to query rows: %w", err)
	}
	//log.Info().Int("row count", len(rows)).Msg("Getting treatments")
	return toTemplateTreatment(rows)
}

func treatmentsByCell(ctx context.Context, org *models.Organization, c h3.Cell) ([]Treatment, error) {
	var results []Treatment
	boundary, err := c.Boundary()
	if err != nil {
		return results, fmt.Errorf("Failed to get cell boundary: %w", err)
	}
	geom_query := gisStatement(boundary)
	rows, err := org.Treatments(
		sm.Where(
			psql.F("ST_Within", "geospatial", geom_query),
		),
		sm.OrderBy("pointlocid"),
		sm.OrderBy("enddatetime"),
	).All(ctx, db.PGInstance.BobDB)
	if err != nil {
		return results, fmt.Errorf("Failed to query rows: %w", err)
	}
	return toTemplateTreatment(rows)
}
func inspectionsByCell(ctx context.Context, org *models.Organization, c h3.Cell) ([]Inspection, error) {
	var results []Inspection

	boundary, err := c.Boundary()
	if err != nil {
		return results, fmt.Errorf("Failed to get cell boundary: %w", err)
	}
	geom_query := gisStatement(boundary)
	rows, err := org.Mosquitoinspections(
		sm.Where(
			psql.F("ST_Within", "geospatial", geom_query),
		),
		sm.OrderBy("pointlocid"),
		sm.OrderBy("enddatetime"),
	).All(ctx, db.PGInstance.BobDB)
	if err != nil {
		return results, fmt.Errorf("Failed to query rows: %w", err)
	}
	return toTemplateInspection(rows)
}
func inspectionsBySource(ctx context.Context, org *models.Organization, sourceID uuid.UUID) ([]Inspection, error) {
	var results []Inspection

	rows, err := org.Mosquitoinspections(
		sm.Where(
			models.FieldseekerMosquitoinspections.Columns.Pointlocid.EQ(psql.Arg(sourceID)),
		),
		sm.OrderBy("enddatetime").Desc(),
	).All(ctx, db.PGInstance.BobDB)
	if err != nil {
		return results, fmt.Errorf("Failed to query rows: %w", err)
	}
	return toTemplateInspection(rows)
}
func breedingSourcesByCell(ctx context.Context, org *models.Organization, c h3.Cell) ([]BreedingSourceSummary, error) {
	var results []BreedingSourceSummary

	boundary, err := c.Boundary()
	if err != nil {
		return results, fmt.Errorf("Failed to get cell boundary: %w", err)
	}
	geom_query := gisStatement(boundary)
	rows, err := org.Pointlocations(
		sm.Where(
			psql.F("ST_Within", "geospatial", geom_query),
		),
		sm.OrderBy("lasttreatdate"),
	).All(ctx, db.PGInstance.BobDB)
	if err != nil {
		return results, fmt.Errorf("Failed to query rows: %w", err)
	}
	for _, r := range rows {
		var last_inspected *time.Time
		if !r.Lastinspectdate.IsNull() {
			l := r.Lastinspectdate.MustGet()
			last_inspected = &l
		}
		var last_treat_date *time.Time
		if !r.Lasttreatdate.IsNull() {
			l := r.Lasttreatdate.MustGet()
			last_treat_date = &l
		}
		results = append(results, BreedingSourceSummary{
			ID:            r.Globalid,
			LastInspected: last_inspected,
			LastTreated:   last_treat_date,
			Type:          r.Habitat.GetOr("none"),
		})
	}
	return results, nil
}

func uuidShort(uuid string) string {
	if len(uuid) < 7 {
		return uuid // Return as is if too short
	}

	return uuid[:3] + "..." + uuid[len(uuid)-4:]
}

func sourceByGlobalId(ctx context.Context, org *models.Organization, id uuid.UUID) (*BreedingSourceDetail, error) {
	row, err := org.Pointlocations(
		sm.Where(models.FieldseekerPointlocations.Columns.Globalid.EQ(psql.Arg(id))),
	).One(ctx, db.PGInstance.BobDB)
	if err != nil {
		return nil, fmt.Errorf("Failed to get point location: %w", err)
	}
	return toTemplateBreedingSource(row), nil
}
