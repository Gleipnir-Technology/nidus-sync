package resource

import (
	"context"
	"net/http"
	"slices"
	"time"

	"github.com/Gleipnir-Technology/nidus-sync/config"
	nhttp "github.com/Gleipnir-Technology/nidus-sync/http"
	"github.com/Gleipnir-Technology/nidus-sync/platform"
	"github.com/google/uuid"
	//"github.com/gorilla/mux"
	//"github.com/rs/zerolog/log"
)

type communicationR struct {
	router *router
}

func Communication(r *router) *communicationR {
	return &communicationR{
		router: r,
	}
}

type communication struct {
	Created      time.Time `json:"created"`
	ID           string    `json:"id"`
	PublicReport string    `json:"public_report"`
	Type         string    `json:"type"`
}
type communicationList struct {
	Communications []communication `json:"communications"`
}

func toImageURLs(m map[string][]uuid.UUID, id string) []string {
	uuids, ok := m[id]
	if !ok {
		return []string{}
	}
	urls := make([]string, len(uuids))
	for i, u := range uuids {
		urls[i] = config.MakeURLNidus("/api/image/%s/content", u.String())
	}
	return urls
}
func (res *communicationR) List(ctx context.Context, r *http.Request, user platform.User, query QueryParams) (*communicationList, *nhttp.ErrorWithStatus) {
	reports, err := platform.PublicReportsForOrganization(ctx, user.Organization.ID, false)
	if err != nil {
		return nil, nhttp.NewError("nuisance report query: %w", err)
	}
	comms := make([]communication, len(reports))
	for i, report := range reports {
		populateDistrictURI(report, res.router)
		populateReportURI(report, res.router, false)
		comms[i] = communication{
			Created:      report.Created,
			ID:           report.PublicID,
			PublicReport: report.URI,
			Type:         "publicreport." + string(report.Type),
		}
	}
	_by_created := func(a, b communication) int {
		if a.Created == b.Created {
			return 0
		} else if a.Created.Before(b.Created) {
			return 1
		} else {
			return -1
		}
	}
	slices.SortFunc(comms, _by_created)
	return &communicationList{
		Communications: comms,
	}, nil
}
