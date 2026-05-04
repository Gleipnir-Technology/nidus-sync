package resource

import (
	"context"
	"net/http"
	"slices"
	"strconv"
	"time"

	"github.com/Gleipnir-Technology/nidus-sync/config"
	modelpublic "github.com/Gleipnir-Technology/nidus-sync/db/gen/nidus-sync/public/model"
	modelpublicreport "github.com/Gleipnir-Technology/nidus-sync/db/gen/nidus-sync/publicreport/model"
	nhttp "github.com/Gleipnir-Technology/nidus-sync/http"
	"github.com/Gleipnir-Technology/nidus-sync/platform"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"
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
	Closed   *time.Time `json:"closed"`
	ClosedBy string     `json:"closed_by"`
	Created  time.Time  `json:"created"`
	ID       string     `json:"id"`
	Opened   *time.Time `json:"opened"`
	OpenedBy string     `json:"opened_by"`
	Response string     `json:"response"`
	Source   string     `json:"source"`
	Type     string     `json:"type"`
	URI      string     `json:"uri"`
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
func (res *communicationR) Get(ctx context.Context, r *http.Request, user platform.User, query QueryParams) (*communication, *nhttp.ErrorWithStatus) {
	return nil, nil
}
func (res *communicationR) List(ctx context.Context, r *http.Request, user platform.User, query QueryParams) ([]communication, *nhttp.ErrorWithStatus) {
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
	public_report_id_to_report := make(map[int32]*modelpublicreport.Report, 0)
	for _, pr := range public_reports {
		public_report_id_to_report[pr.ID] = pr
	}
	result := make([]communication, len(comms))
	for i, comm := range comms {
		public_report, ok := public_report_id_to_report[*comm.SourceReportID]
		if !ok {
			return nil, nhttp.NewError("lookup report id %d failed", comm.SourceReportID)
		}
		c, err := res.hydrateCommunication(*comm, public_report)
		if err != nil {
			return nil, err
		}
		result[i] = c
	}
	_by_created := func(a, b communication) int {
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

type communicationMarkRequest struct{}

func (res *communicationR) MarkInvalid(ctx context.Context, r *http.Request, user platform.User, cmr communicationMarkRequest) (communication, *nhttp.ErrorWithStatus) {
	return res.markCommunication(ctx, r, user, "invalid", platform.CommunicationMarkInvalid)
}
func (res *communicationR) MarkPendingResponse(ctx context.Context, r *http.Request, user platform.User, cmr communicationMarkRequest) (communication, *nhttp.ErrorWithStatus) {
	return res.markCommunication(ctx, r, user, "pending-response", platform.CommunicationMarkPendingResponse)
}
func (res *communicationR) MarkPossibleIssue(ctx context.Context, r *http.Request, user platform.User, cmr communicationMarkRequest) (communication, *nhttp.ErrorWithStatus) {
	return res.markCommunication(ctx, r, user, "possible-issue", platform.CommunicationMarkPossibleIssue)
}
func (res *communicationR) MarkPossibleResolved(ctx context.Context, r *http.Request, user platform.User, cmr communicationMarkRequest) (communication, *nhttp.ErrorWithStatus) {
	return res.markCommunication(ctx, r, user, "possible-resolved", platform.CommunicationMarkPossibleResolved)
}
func (res *communicationR) hydrateCommunication(comm modelpublic.Communication, public_report *modelpublicreport.Report) (communication, *nhttp.ErrorWithStatus) {
	var err error
	source_uri := "unknown"
	type_ := "unknown"
	if comm.SourceReportID != nil && public_report != nil {
		source_uri, err = reportURI(res.router, "", public_report.PublicID)
		if err != nil {
			return communication{}, nhttp.NewError("gen report URI: %w", err)
		}
		type_ = "publicreport." + public_report.ReportType.String()
	} else if comm.SourceEmailLogID != nil {
		source_uri, err = emailURI(*res.router, *comm.SourceEmailLogID)
		if err != nil {
			return communication{}, nhttp.NewError("gen email URI: %w", err)
		}
		type_ = "email"
	} else if comm.SourceTextLogID != nil {
		source_uri, err = textURI(*res.router, *comm.SourceTextLogID)
		if err != nil {
			return communication{}, nhttp.NewError("gen email URI: %w", err)
		}
		source_uri = "text"
	}
	closed_by, err := userURI(res.router, comm.ClosedBy)
	if err != nil {
		return communication{}, nhttp.NewError("gen closed_by URI: %w", err)
	}
	opened_by, err := userURI(res.router, comm.OpenedBy)
	if err != nil {
		return communication{}, nhttp.NewError("gen opened_by URI: %w", err)
	}
	response, err := responseURI(*res.router, comm)
	if err != nil {
		return communication{}, nhttp.NewError("gen response URI: %w", err)
	}
	uri, err := res.router.IDToURI("communication.ByIDGet", int(comm.ID))
	if err != nil {
		return communication{}, nhttp.NewError("gen comm uri: %w", err)
	}
	return communication{
		Closed:   comm.Closed,
		ClosedBy: closed_by,
		Created:  comm.Created,
		ID:       strconv.Itoa(int(comm.ID)),
		Opened:   comm.Opened,
		OpenedBy: opened_by,
		Response: response,
		Source:   source_uri,
		Type:     type_,
		URI:      uri,
	}, nil
}

type markFunc = func(context.Context, platform.User, int64) error

func (res *communicationR) markCommunication(ctx context.Context, r *http.Request, user platform.User, status string, m markFunc) (communication, *nhttp.ErrorWithStatus) {
	vars := mux.Vars(r)
	comm_id_str := vars["id"]
	if comm_id_str == "" {
		return communication{}, nhttp.NewBadRequest("no id provided")
	}
	comm_id, err := strconv.Atoi(comm_id_str)
	if err != nil {
		return communication{}, nhttp.NewBadRequest("can't turn report ID into an int: %w", err)
	}
	m(ctx, user, int64(comm_id))
	result, err := platform.CommunicationFromID(ctx, user, int64(comm_id))
	if result == nil {
		return communication{}, nhttp.NewUnauthorized("you are not authorized to modify communication %d", comm_id)
	}
	var public_report *modelpublicreport.Report
	if result.SourceReportID != nil {
		comm_ids := []int64{int64(*result.SourceReportID)}
		public_reports, err := platform.PublicReportsFromIDs(ctx, comm_ids)
		if err != nil {
			return communication{}, nhttp.NewError("Get report %d: %w", *result.SourceReportID, err)
		}
		public_report = public_reports[0]
	}

	log.Info().Int("communication", comm_id).Str("status", status).Msg("Marked communication")
	return res.hydrateCommunication(*result, public_report)
}
func responseURI(r router, comm modelpublic.Communication) (string, error) {
	if comm.ResponseEmailLogID != nil {
		return emailURI(r, *comm.ResponseEmailLogID)
	} else if comm.ResponseTextLogID != nil {
		return textURI(r, *comm.ResponseTextLogID)
	} else {
		return "", nil
	}
}
func emailURI(r router, id int32) (string, error) {
	return "fake email uri", nil
}

func textURI(r router, id int32) (string, error) {
	return "fake text uri", nil
}
func userURI(r *router, id *int32) (string, error) {
	if id == nil {
		return "", nil
	}
	return r.IDToURI("user.ByIDGet", int(*id))
}
