package sync

import (
	"context"
	"net/http"

	"github.com/Gleipnir-Technology/nidus-sync/html"
	nhttp "github.com/Gleipnir-Technology/nidus-sync/http"
	"github.com/Gleipnir-Technology/nidus-sync/platform"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type contentSource struct {
	Inspections []platform.Inspection
	Source      *platform.BreedingSourceDetail
	Traps       []platform.TrapNearby
	Treatments  []platform.Treatment
	//TreatmentCadence TreatmentCadence
	TreatmentModels []platform.TreatmentModel
	User            platform.User
}
type contentTrap struct {
	Trap platform.Trap
	User platform.User
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
	html.RenderOrError(w, "static/gen/main.html", struct{}{})
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
	data := contentSource{
		Inspections:     inspections,
		Source:          s,
		Traps:           traps,
		Treatments:      treatments,
		TreatmentModels: treatment_models,
		User:            user,
	}

	return html.NewResponse("sync/source.html", data), nil
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
	/*
		latlng, err := t.H3Cell.LatLng()
		if err != nil {
			return nil, nhttp.NewError("Failed to get latlng: %w", err)
		}
	*/
	data := contentTrap{
		Trap: *t,
		User: user,
	}
	return html.NewResponse("sync/trap.html", data), nil
}

func source(w http.ResponseWriter, r *http.Request, user platform.User, id uuid.UUID) {
}

func trap(w http.ResponseWriter, r *http.Request, user platform.User, id uuid.UUID) {
}
