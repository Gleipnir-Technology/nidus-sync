package sync

import (
	"embed"

	"github.com/Gleipnir-Technology/nidus-sync/htmlpage"

	//"bytes"
	"context"
	//"errors"
	"fmt"
	//"html/template"
	//"io"
	//"math"
	"net/http"
	//"os"
	//"strconv"
	//"strings"
	"time"

	"github.com/Gleipnir-Technology/nidus-sync/background"
	"github.com/Gleipnir-Technology/nidus-sync/config"
	"github.com/Gleipnir-Technology/nidus-sync/db"
	"github.com/Gleipnir-Technology/nidus-sync/db/models"
	"github.com/Gleipnir-Technology/nidus-sync/h3utils"
	//"github.com/aarondl/opt/null"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"github.com/stephenafamo/bob/dialect/psql/sm"
	"github.com/uber/h3-go/v4"
)

//go:embed template/*
var embeddedFiles embed.FS

//go:embed static/*
var EmbeddedStaticFS embed.FS

// Authenticated pages
var (
	cell        = buildTemplate("cell", "authenticated")
	dashboard   = buildTemplate("dashboard", "authenticated")
	oauthPrompt = buildTemplate("oauth-prompt", "authenticated")
	settings    = buildTemplate("settings", "authenticated")
	source      = buildTemplate("source", "authenticated")
)

// Unauthenticated pages
var (
	admin                           = buildTemplate("admin", "base")
	dataEntry                       = buildTemplate("data-entry", "base")
	dataEntryGood                   = buildTemplate("data-entry-good", "base")
	dataEntryBad                    = buildTemplate("data-entry-bad", "base")
	dispatch                        = buildTemplate("dispatch", "base")
	dispatchResults                 = buildTemplate("dispatch-results", "base")
	mockRoot                        = buildTemplate("mock-root", "base")
	reportPage                      = buildTemplate("report", "base")
	reportConfirmation              = buildTemplate("report-confirmation", "base")
	reportContribute                = buildTemplate("report-contribute", "base")
	reportDetail                    = buildTemplate("report-detail", "base")
	reportEvidence                  = buildTemplate("report-evidence", "base")
	reportSchedule                  = buildTemplate("report-schedule", "base")
	reportUpdate                    = buildTemplate("report-update", "base")
	serviceRequest                  = buildTemplate("service-request", "base")
	serviceRequestDetail            = buildTemplate("service-request-detail", "base")
	serviceRequestLocation          = buildTemplate("service-request-location", "base")
	serviceRequestMosquito          = buildTemplate("service-request-mosquito", "base")
	serviceRequestPool              = buildTemplate("service-request-pool", "base")
	serviceRequestQuick             = buildTemplate("service-request-quick", "base")
	serviceRequestQuickConfirmation = buildTemplate("service-request-quick-confirmation", "base")
	serviceRequestUpdates           = buildTemplate("service-request-updates", "base")
	settingRoot                     = buildTemplate("setting-mock", "base")
	settingIntegration              = buildTemplate("setting-integration", "base")
	settingPesticide                = buildTemplate("setting-pesticide", "base")
	settingPesticideAdd             = buildTemplate("setting-pesticide-add", "base")
	settingUsers                    = buildTemplate("setting-user", "base")
	settingUsersAdd                 = buildTemplate("setting-user-add", "base")
	signin                          = buildTemplate("signin", "base")
	signup                          = buildTemplate("signup", "base")
)

var components = [...]string{"header", "map"}

func buildTemplate(files ...string) *htmlpage.BuiltTemplate {
	subdir := "sync"
	full_files := make([]string, 0)
	for _, f := range files {
		full_files = append(full_files, fmt.Sprintf("%s/template/%s.html", subdir, f))
	}
	for _, c := range components {
		full_files = append(full_files, fmt.Sprintf("%s/template/components/%s.html", subdir, c))
	}
	return htmlpage.NewBuiltTemplate(embeddedFiles, "sync/", full_files...)
}

func Cell(ctx context.Context, w http.ResponseWriter, user *models.User, c int64) {
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
	geojson, err := h3utils.H3ToGeoJSON([]h3.Cell{h3.Cell(c)})
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
			MapboxToken: config.MapboxToken,
			Zoom:        resolution + 5,
		},
		Treatments: treatments,
		User:       userContent,
	}
	htmlpage.RenderOrError(w, cell, &data)
}

func Dashboard(ctx context.Context, w http.ResponseWriter, user *models.User) {
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
	is_syncing := background.IsSyncOngoing(org.ID)
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
			MapboxToken: config.MapboxToken,
		},
		Org:            org.Name.MustGet(),
		RecentRequests: requests,
		User:           userContent,
	}
	htmlpage.RenderOrError(w, dashboard, data)
}

func Mock(t string, w http.ResponseWriter, code string) {
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
	template, ok := htmlpage.TemplatesByFilename[t+".html"]
	if !ok {
		log.Error().Str("template", t).Msg("Failed to find template")
		respondError(w, "Failed to render template", nil, http.StatusInternalServerError)
		return
	}
	htmlpage.RenderOrError(w, &template, data)
}

func OauthPrompt(w http.ResponseWriter, user *models.User) {
	dp := user.DisplayName
	data := ContentDashboard{
		User: User{
			DisplayName: dp,
			Initials:    extractInitials(dp),
			Username:    user.Username,
		},
	}
	htmlpage.RenderOrError(w, oauthPrompt, data)
}

func Settings(w http.ResponseWriter, r *http.Request, user *models.User) {
	userContent, err := contentForUser(r.Context(), user)
	if err != nil {
		respondError(w, "Failed to get user content", err, http.StatusInternalServerError)
		return
	}
	data := ContentAuthenticatedPlaceholder{
		User: userContent,
	}
	htmlpage.RenderOrError(w, settings, data)
}

func Signin(w http.ResponseWriter, errorCode string) {
	data := ContentSignin{
		InvalidCredentials: errorCode == "invalid-credentials",
	}
	htmlpage.RenderOrError(w, signin, data)
}

func Signup(w http.ResponseWriter, path string) {
	data := ContentSignup{}
	htmlpage.RenderOrError(w, signup, data)
}

func Source(w http.ResponseWriter, r *http.Request, user *models.User, id uuid.UUID) {
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
			MapboxToken: config.MapboxToken,
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

	htmlpage.RenderOrError(w, source, data)
}

// Respond with an error that is visible to the user
func respondError(w http.ResponseWriter, m string, e error, s int) {
	log.Warn().Int("status", s).Err(e).Str("user message", m).Msg("Responding with an error from sync pages")
	http.Error(w, m, s)
}
