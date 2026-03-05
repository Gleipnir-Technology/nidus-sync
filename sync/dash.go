package sync

import (
	"context"
	"errors"
	"html/template"
	"net/http"
	"time"

	"github.com/Gleipnir-Technology/bob/dialect/psql/sm"
	"github.com/Gleipnir-Technology/nidus-sync/auth"
	"github.com/Gleipnir-Technology/nidus-sync/background"
	"github.com/Gleipnir-Technology/nidus-sync/config"
	"github.com/Gleipnir-Technology/nidus-sync/db"
	"github.com/Gleipnir-Technology/nidus-sync/db/models"
	"github.com/Gleipnir-Technology/nidus-sync/html"
	nhttp "github.com/Gleipnir-Technology/nidus-sync/http"
	"github.com/Gleipnir-Technology/nidus-sync/platform"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

// Authenticated pages
var ()

type contentSource struct {
	Inspections []Inspection
	MapData     ComponentMap
	Source      *BreedingSourceDetail
	Traps       []TrapNearby
	Treatments  []Treatment
	//TreatmentCadence TreatmentCadence
	TreatmentModels []TreatmentModel
	User            platform.User
}
type contentTrap struct {
	MapData ComponentMap
	Trap    Trap
	User    platform.User
}
type contentDashboard struct {
	CountTraps           int
	CountMosquitoSources int
	CountServiceRequests int
	Geo                  template.JS
	IsSyncOngoing        bool
	LastSync             *time.Time
	MapData              ComponentMap
	RecentRequests       []ServiceRequestSummary
}

type contentLayoutTest struct {
	User platform.User
}
type ContentDistrict struct {
	MapboxToken string
}

func getDistrict(w http.ResponseWriter, r *http.Request) {
	context := ContentDistrict{
		MapboxToken: config.MapboxToken,
	}
	html.RenderOrError(w, "sync/district.html", &context)
}

func getLayoutTest(ctx context.Context, r *http.Request, org *models.Organization, user *models.User) (*html.Response[contentLayoutTest], *nhttp.ErrorWithStatus) {
	return html.NewResponse("sync/layout-test.html", contentLayoutTest{}), nil
}

func getRoot(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	user, err := auth.GetAuthenticatedUser(r)
	if err != nil {
		// No credentials or user not found: go to login
		if errors.Is(err, &auth.NoCredentialsError{}) || errors.Is(err, &auth.NoUserError{}) {
			http.Redirect(w, r, "/signin", http.StatusFound)
			return
		} else {
			respondError(w, "Failed to get root", err, http.StatusInternalServerError)
			return
		}
	}
	if user == nil {
		errorCode := r.URL.Query().Get("error")
		signin(w, errorCode, "/")
		return
	} else {
		org, err := user.Organization().One(ctx, db.PGInstance.BobDB)
		if err != nil {
			respondError(w, "Failed to get organization", err, http.StatusInternalServerError)
			return
		}
		dashboard(ctx, w, org, user)
		return
	}
}

func getSource(ctx context.Context, r *http.Request, org *models.Organization, user *models.User) (*html.Response[contentSource], *nhttp.ErrorWithStatus) {
	globalid_s := chi.URLParam(r, "globalid")
	if globalid_s == "" {
		return nil, nhttp.NewError("No globalid provided: %w", nil)
	}
	globalid, err := uuid.Parse(globalid_s)
	if err != nil {
		return nil, nhttp.NewError("globalid is not a UUID: %w", nil)
	}
	userContent, err := auth.ContentForUser(r.Context(), user)
	if err != nil {
		return nil, nhttp.NewError("Failed to get user content: %w", err)
	}
	s, err := sourceByGlobalId(r.Context(), org, globalid)
	if err != nil {
		return nil, nhttp.NewError("Failed to get source: %w", err)
	}
	inspections, err := inspectionsBySource(r.Context(), org, globalid)
	if err != nil {
		return nil, nhttp.NewError("Failed to get inspections: %w", err)
	}
	traps, err := trapsBySource(r.Context(), org, globalid)
	if err != nil {
		return nil, nhttp.NewError("Failed to get traps: %w", err)
	}

	treatments, err := treatmentsBySource(r.Context(), org, globalid)
	if err != nil {
		return nil, nhttp.NewError("Failed to get treatments: %w", err)
	}
	treatment_models := modelTreatment(treatments)
	latlng, err := s.H3Cell.LatLng()
	if err != nil {
		return nil, nhttp.NewError("Failed to get latlng: %w", err)
	}
	data := contentSource{
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

	return html.NewResponse("sync/source.html", data), nil
}

func getStadia(ctx context.Context, r *http.Request, org *models.Organization, u *models.User) (*html.Response[contentDashboard], *nhttp.ErrorWithStatus) {
	data := contentDashboard{
		MapData: ComponentMap{
			MapboxToken: config.MapboxToken,
		},
	}
	return html.NewResponse("sync/stadia.html", data), nil
}
func getTemplateTest(w http.ResponseWriter, r *http.Request) {
	html.RenderOrError(w, "sync/template-test.html", nil)
}
func getTrap(ctx context.Context, r *http.Request, org *models.Organization, user *models.User) (*html.Response[contentTrap], *nhttp.ErrorWithStatus) {
	globalid_s := chi.URLParam(r, "globalid")
	if globalid_s == "" {
		return nil, nhttp.NewError("No globalid provided: %w", nil)
	}
	globalid, err := uuid.Parse(globalid_s)
	if err != nil {
		return nil, nhttp.NewError("globalid is not a UUID: %w", nil)
	}
	userContent, err := auth.ContentForUser(r.Context(), user)
	if err != nil {
		return nil, nhttp.NewError("Failed to get user content: %w", err)
	}
	t, err := trapByGlobalId(r.Context(), org, globalid)
	if err != nil {
		return nil, nhttp.NewError("Failed to get trap: %w", err)
	}
	latlng, err := t.H3Cell.LatLng()
	if err != nil {
		return nil, nhttp.NewError("Failed to get latlng: %w", err)
	}
	data := contentTrap{
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
		Trap: t,
		User: userContent,
	}
	return html.NewResponse("sync/trap.html", data), nil
}

func dashboard(ctx context.Context, w http.ResponseWriter, org *models.Organization, user *models.User) {
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
	trapCount, err := org.Traplocations().Count(ctx, db.PGInstance.BobDB)
	if err != nil {
		respondError(w, "Failed to get trap count", err, http.StatusInternalServerError)
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
	content := contentDashboard{
		CountTraps:           int(trapCount),
		CountMosquitoSources: int(sourceCount),
		CountServiceRequests: int(serviceCount),
		IsSyncOngoing:        is_syncing,
		LastSync:             lastSync,
		MapData: ComponentMap{
			MapboxToken: config.MapboxToken,
		},
		RecentRequests: requests,
	}
	userContent, err := auth.ContentForUser(ctx, user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	html.RenderOrError(w, "sync/dashboard.html", contentAuthenticated[contentDashboard]{
		C:            content,
		Organization: org,
		URL:          html.NewContentURL(),
		User:         userContent,
	})
}

func source(w http.ResponseWriter, r *http.Request, org *models.Organization, user *models.User, id uuid.UUID) {
}

func trap(w http.ResponseWriter, r *http.Request, org *models.Organization, user *models.User, id uuid.UUID) {
}
