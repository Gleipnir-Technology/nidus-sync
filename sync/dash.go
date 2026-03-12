package sync

import (
	"context"
	"errors"
	"html/template"
	"net/http"
	"time"

	"github.com/Gleipnir-Technology/nidus-sync/auth"
	"github.com/Gleipnir-Technology/nidus-sync/html"
	nhttp "github.com/Gleipnir-Technology/nidus-sync/http"
	"github.com/Gleipnir-Technology/nidus-sync/platform"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type contentSource struct {
	Inspections []platform.Inspection
	MapData     ComponentMap
	Source      *platform.BreedingSourceDetail
	Traps       []platform.TrapNearby
	Treatments  []platform.Treatment
	//TreatmentCadence TreatmentCadence
	TreatmentModels []platform.TreatmentModel
	User            platform.User
}
type contentTrap struct {
	MapData ComponentMap
	Trap    platform.Trap
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
}

func getDistrict(w http.ResponseWriter, r *http.Request) {
	context := ContentDistrict{}
	html.RenderOrError(w, "sync/district.html", &context)
}

func getLayoutTest(ctx context.Context, r *http.Request, user platform.User) (*html.Response[contentLayoutTest], *nhttp.ErrorWithStatus) {
	return html.NewResponse("sync/layout-test.html", contentLayoutTest{}), nil
}

func getRoot(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	user, err := auth.GetAuthenticatedUser(r)
	if err != nil {
		// No credentials or user not found: go to login
		if errors.Is(err, &auth.NoCredentialsError{}) || errors.Is(err, &platform.NoUserError{}) {
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
		dashboard(ctx, w, *user)
		return
	}
}

func getSource(ctx context.Context, r *http.Request, user platform.User) (*html.Response[contentSource], *nhttp.ErrorWithStatus) {
	globalid_s := chi.URLParam(r, "globalid")
	if globalid_s == "" {
		return nil, nhttp.NewError("No globalid provided: %w", nil)
	}
	globalid, err := uuid.Parse(globalid_s)
	if err != nil {
		return nil, nhttp.NewError("globalid is not a UUID: %w", nil)
	}
	s, err := platform.SourceByGlobalID(ctx, user.Organization, globalid)
	if err != nil {
		return nil, nhttp.NewError("Failed to get source: %w", err)
	}
	inspections, err := platform.InspectionsBySource(ctx, user.Organization, globalid)
	if err != nil {
		return nil, nhttp.NewError("Failed to get inspections: %w", err)
	}
	traps, err := platform.TrapsBySource(ctx, user.Organization, globalid)
	if err != nil {
		return nil, nhttp.NewError("Failed to get traps: %w", err)
	}

	treatments, err := platform.TreatmentsBySource(ctx, user.Organization, globalid)
	if err != nil {
		return nil, nhttp.NewError("Failed to get treatments: %w", err)
	}
	treatment_models := platform.ModelTreatment(treatments)
	latlng, err := s.H3Cell.LatLng()
	if err != nil {
		return nil, nhttp.NewError("Failed to get latlng: %w", err)
	}
	data := contentSource{
		Inspections: inspections,
		MapData: ComponentMap{
			Center: latlng,
			//GeoJSON:
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
		User:            user,
	}

	return html.NewResponse("sync/source.html", data), nil
}

func getStadia(ctx context.Context, r *http.Request, u platform.User) (*html.Response[contentDashboard], *nhttp.ErrorWithStatus) {
	data := contentDashboard{
		MapData: ComponentMap{},
	}
	return html.NewResponse("sync/stadia.html", data), nil
}
func getTemplateTest(w http.ResponseWriter, r *http.Request) {
	html.RenderOrError(w, "sync/template-test.html", nil)
}
func getTrap(ctx context.Context, r *http.Request, user platform.User) (*html.Response[contentTrap], *nhttp.ErrorWithStatus) {
	globalid_s := chi.URLParam(r, "globalid")
	if globalid_s == "" {
		return nil, nhttp.NewError("No globalid provided: %w", nil)
	}
	globalid, err := uuid.Parse(globalid_s)
	if err != nil {
		return nil, nhttp.NewError("globalid is not a UUID: %w", nil)
	}
	t, err := platform.TrapByGlobalId(ctx, user.Organization, globalid)
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
			Markers: []MapMarker{
				MapMarker{
					LatLng: latlng,
				},
			},
			Zoom: 13,
		},
		Trap: *t,
		User: user,
	}
	return html.NewResponse("sync/trap.html", data), nil
}

func dashboard(ctx context.Context, w http.ResponseWriter, user platform.User) {
	var lastSync *time.Time
	sync, err := user.Organization.FieldseekerSyncLatest(ctx)
	if err != nil {
		respondError(w, "Failed to get syncs", err, http.StatusInternalServerError)
	} else if sync != nil {
		lastSync = &sync.Created
	}
	is_syncing := user.Organization.IsSyncOngoing()
	count_trap, err := user.Organization.CountTrap(ctx)
	if err != nil {
		respondError(w, "Failed to get trap count", err, http.StatusInternalServerError)
		return
	}
	count_source, err := user.Organization.CountTrap(ctx)
	if err != nil {
		respondError(w, "Failed to get source count", err, http.StatusInternalServerError)
		return
	}
	count_service, err := user.Organization.CountServiceRequest(ctx)
	if err != nil {
		respondError(w, "Failed to get service count", err, http.StatusInternalServerError)
		return
	}
	service_request_recent, err := user.Organization.ServiceRequestRecent(ctx)
	if err != nil {
		respondError(w, "Failed to get recent service", err, http.StatusInternalServerError)
		return
	}

	requests := make([]ServiceRequestSummary, 0)
	for _, r := range service_request_recent {
		requests = append(requests, ServiceRequestSummary{
			Date:     r.Creationdate.MustGet(),
			Location: r.Reqaddr1.MustGet(),
			Status:   "Completed",
		})
	}
	content := contentDashboard{
		CountTraps:           int(count_trap),
		CountMosquitoSources: int(count_source),
		CountServiceRequests: int(count_service),
		IsSyncOngoing:        is_syncing,
		LastSync:             lastSync,
		MapData:              ComponentMap{},
		RecentRequests:       requests,
	}
	html.RenderOrError(w, "sync/dashboard.html", contentAuthenticated[contentDashboard]{
		C:            content,
		Config:       html.NewContentConfig(),
		Organization: user.Organization,
		URL:          html.NewContentURL(),
		User:         user,
	})
}

func source(w http.ResponseWriter, r *http.Request, user platform.User, id uuid.UUID) {
}

func trap(w http.ResponseWriter, r *http.Request, user platform.User, id uuid.UUID) {
}
