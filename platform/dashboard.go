package platform

import (
	"context"
	"time"

	nhttp "github.com/Gleipnir-Technology/nidus-sync/http"
)
type ServiceRequestSummary struct {
	Date     time.Time
	Location string
	Status   string
}
type contentDashboard struct {
	CountTraps           int
	CountMosquitoSources int
	CountServiceRequests int
	IsSyncOngoing        bool
	LastSync             *time.Time
	RecentRequests       []ServiceRequestSummary
}

func getDashboardData(ctx context.Context, user User) (*contentDashboard, *nhttp.ErrorWithStatus) {
	var lastSync *time.Time
	sync, err := user.Organization.FieldseekerSyncLatest(ctx)
	if err != nil {
		return nil, nhttp.NewError("Failed to get syncs: %w", err)
	} else if sync != nil {
		lastSync = &sync.Created
	}
	is_syncing := user.Organization.IsSyncOngoing()
	count_trap, err := user.Organization.CountTrap(ctx)
	if err != nil {
		return nil, nhttp.NewError("Failed to get trap count: %w", err)
	}
	count_source, err := user.Organization.CountTrap(ctx)
	if err != nil {
		return nil, nhttp.NewError("Failed to get source count: %w", err)
	}
	count_service, err := user.Organization.CountServiceRequest(ctx)
	if err != nil {
		return nil, nhttp.NewError("Failed to get service count: %w", err)
	}
	service_request_recent, err := user.Organization.ServiceRequestRecent(ctx)
	if err != nil {
		return nil, nhttp.NewError("Failed to get recent service: %w", err)
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
		RecentRequests:       requests,
	}
	return &content, nil
}
