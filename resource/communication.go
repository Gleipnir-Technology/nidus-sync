package resource

import (
	"context"
	"net/http"
	"slices"
	"time"

	"github.com/Gleipnir-Technology/nidus-sync/config"
	"github.com/Gleipnir-Technology/nidus-sync/db/gen/nidus-sync/publicreport/model"
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
	Created time.Time `json:"created"`
	ID      int32     `json:"id"`
	Source  string    `json:"source"`
	Type    string    `json:"type"`
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
func (res *communicationR) List(ctx context.Context, r *http.Request, user platform.User, query QueryParams) ([]*communication, *nhttp.ErrorWithStatus) {
	comms, err := platform.CommunicationsForOrganization(ctx, int64(user.Organization.ID))
	if err != nil {
		return nil, nhttp.NewError("nuisance report query: %w", err)
	}
	report_ids := make([]int64, 0)
	for _, comm := range comms {
		if comm.SourceReportID != nil {
			report_ids = append(report_ids, int64(*comm.SourceReportID))
		}
	}
	public_reports, err := platform.PublicReportsFromIDs(ctx, report_ids)
	if err != nil {
		return nil, nhttp.NewError("public reports from IDs: %w", err)
	}
	public_report_id_to_report := make(map[int32]*model.Report, 0)
	for _, pr := range public_reports {
		public_report_id_to_report[pr.ID] = pr
	}
	result := make([]*communication, len(comms))
	for i, comm := range comms {
		source_uri := "unknown"
		type_ := "unknown"
		if comm.SourceReportID != nil {
			public_report, ok := public_report_id_to_report[*comm.SourceReportID]
			if !ok {
				return nil, nhttp.NewError("lookup report id %d failed", comm.SourceReportID)
			}
			source_uri, err = reportURI(res.router, "", public_report.PublicID)
			if err != nil {
				return nil, nhttp.NewError("gen report URI: %w", err)
			}
			type_ = "publicreport." + public_report.ReportType.String()
		} else if comm.SourceEmailLogID != nil {
			source_uri, err = emailURI(res.router, *comm.SourceEmailLogID)
			if err != nil {
				return nil, nhttp.NewError("gen email URI: %w", err)
			}
			type_ = "email"
		} else if comm.SourceTextLogID != nil {
			source_uri, err = textURI(res.router, *comm.SourceTextLogID)
			if err != nil {
				return nil, nhttp.NewError("gen email URI: %w", err)
			}
			source_uri = "text"
		}
		result[i] = &communication{
			Created: comm.Created,
			ID:      comm.ID,
			Source:  source_uri,
			Type:    type_,
		}
	}
	_by_created := func(a, b *communication) int {
		if a.Created.Equal(b.Created) {
			return 0
		} else if a.Created.Before(b.Created) {
			return 1
		} else {
			return -1
		}
	}
	slices.SortFunc(result, _by_created)
	return result, nil
}

func emailURI(r *router, id int32) (string, error) {
	return "fake email uri", nil
}
func textURI(r *router, id int32) (string, error) {
	return "fake text uri", nil
}
