package api

import (
	"context"
	"net/http"

	nhttp "github.com/Gleipnir-Technology/nidus-sync/http"
	"github.com/Gleipnir-Technology/nidus-sync/platform"
)

type createLead struct {
	PoolLocations map[int]platform.Location `json:"pool_locations"`
	SignalIDs     []int                     `json:"signal_ids"`
}
type createdLead struct {
	ID int32 `json:"id"`
}
type contentListLead struct {
	Leads []lead `json:"leads"`
}
type lead struct {
	ID int32 `json:"id"`
}

func listLead(ctx context.Context, r *http.Request, user platform.User, query queryParams) (*contentListLead, *nhttp.ErrorWithStatus) {
	return &contentListLead{
		Leads: make([]lead, 0),
	}, nil
}
func postLeads(ctx context.Context, r *http.Request, user platform.User, req createLead) (*createdLead, *nhttp.ErrorWithStatus) {
	if len(req.SignalIDs) == 0 {
		return nil, nhttp.NewErrorStatus(http.StatusBadRequest, "can't make a lead with no signals")
	}
	if len(req.SignalIDs) > 1 {
		return nil, nhttp.NewErrorStatus(http.StatusBadRequest, "can't make a lead with multiple signals yet")
	}
	signal_id := req.SignalIDs[0]
	var pool_location *platform.Location
	l, ok := req.PoolLocations[signal_id]
	if ok {
		pool_location = &l
	}
	site_id, err := platform.SiteFromSignal(ctx, user, int32(signal_id))
	if err != nil || site_id == nil {
		return nil, nhttp.NewError("site from signal: %w", err)
	}
	lead_id, err := platform.LeadCreate(ctx, user, int32(signal_id), *site_id, pool_location)
	if err != nil || lead_id == nil {
		return nil, nhttp.NewError("lead create: %w", err)
	}

	return &createdLead{
		ID: *lead_id,
	}, nil
}
