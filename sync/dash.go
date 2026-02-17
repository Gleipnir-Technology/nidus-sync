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
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

// Authenticated pages
var ()

type Config struct {
}

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
type ContentTrap struct {
	MapData ComponentMap
	Trap    Trap
	User    User
}
type ContentDashboard struct {
	CountTraps           int
	CountMosquitoSources int
	CountServiceRequests int
	Geo                  template.JS
	IsSyncOngoing        bool
	LastSync             *time.Time
	MapData              ComponentMap
	Organization         *models.Organization
	RecentRequests       []ServiceRequestSummary
	URL                  ContentURL
	User                 User
}

type ContentLayoutTest struct {
	User User
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

func getLayoutTest(w http.ResponseWriter, r *http.Request, u *models.User) {
	userContent, err := contentForUser(r.Context(), u)
	if err != nil {
		respondError(w, "Failed to get user", err, http.StatusInternalServerError)
		return
	}
	html.RenderOrError(w, "sync/layout-test.html", &ContentLayoutTest{User: userContent})
}

func getRoot(w http.ResponseWriter, r *http.Request) {
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
		has, err := background.HasFieldseekerConnection(r.Context(), user)
		if err != nil {
			respondError(w, "Failed to check for ArcGIS connection", err, http.StatusInternalServerError)
			return
		}
		if has {
			dashboard(r.Context(), w, user)
			return
		} else {
			oauthPrompt(w, r, user)
			return
		}
	}
	if err != nil {
		respondError(w, "Failed to render root", err, http.StatusInternalServerError)
	}
}

func getSource(w http.ResponseWriter, r *http.Request, u *models.User) {
	globalid_s := chi.URLParam(r, "globalid")
	if globalid_s == "" {
		respondError(w, "No globalid provided", nil, http.StatusBadRequest)
		return
	}
	globalid, err := uuid.Parse(globalid_s)
	if err != nil {
		respondError(w, "globalid is not a UUID", nil, http.StatusBadRequest)
		return
	}
	source(w, r, u, globalid)
}

func getStadia(w http.ResponseWriter, r *http.Request, u *models.User) {
	userContent, err := contentForUser(r.Context(), u)
	if err != nil {
		respondError(w, "Failed to get user content", err, http.StatusInternalServerError)
		return
	}
	data := ContentDashboard{
		MapData: ComponentMap{
			MapboxToken: config.MapboxToken,
		},
		URL:  newContentURL(),
		User: userContent,
	}
	html.RenderOrError(w, "sync/stadia.html", data)

}
func getTemplateTest(w http.ResponseWriter, r *http.Request) {
	html.RenderOrError(w, "sync/template-test.html", nil)
}
func getTrap(w http.ResponseWriter, r *http.Request, u *models.User) {
	globalid_s := chi.URLParam(r, "globalid")
	if globalid_s == "" {
		respondError(w, "No globalid provided", nil, http.StatusBadRequest)
		return
	}
	globalid, err := uuid.Parse(globalid_s)
	if err != nil {
		respondError(w, "globalid is not a UUID", nil, http.StatusBadRequest)
		return
	}
	trap(w, r, u, globalid)
}

func dashboard(ctx context.Context, w http.ResponseWriter, user *models.User) {
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
	userContent, err := contentForUser(ctx, user)
	if err != nil {
		respondError(w, "Failed to get user context", err, http.StatusInternalServerError)
		return
	}
	data := ContentDashboard{
		CountTraps:           int(trapCount),
		CountMosquitoSources: int(sourceCount),
		CountServiceRequests: int(serviceCount),
		IsSyncOngoing:        is_syncing,
		LastSync:             lastSync,
		MapData: ComponentMap{
			MapboxToken: config.MapboxToken,
		},
		Organization:   org,
		RecentRequests: requests,
		URL:            newContentURL(),
		User:           userContent,
	}
	html.RenderOrError(w, "sync/dashboard.html", data)
}

func source(w http.ResponseWriter, r *http.Request, user *models.User, id uuid.UUID) {
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

	html.RenderOrError(w, "sync/source.html", data)
}

func trap(w http.ResponseWriter, r *http.Request, user *models.User, id uuid.UUID) {
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
	t, err := trapByGlobalId(r.Context(), org, id)
	if err != nil {
		respondError(w, "Failed to get trap", err, http.StatusInternalServerError)
		return
	}
	latlng, err := t.H3Cell.LatLng()
	if err != nil {
		respondError(w, "Failed to get latlng", err, http.StatusInternalServerError)
		return
	}
	data := ContentTrap{
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

	html.RenderOrError(w, "sync/trap.html", data)
}
