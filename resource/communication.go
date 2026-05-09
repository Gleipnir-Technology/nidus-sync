package resource

import (
	"context"
	"net/http"
	"strconv"
	"time"

	modelpublic "github.com/Gleipnir-Technology/nidus-sync/db/gen/nidus-sync/public/model"
	modelpublicreport "github.com/Gleipnir-Technology/nidus-sync/db/gen/nidus-sync/publicreport/model"
	nhttp "github.com/Gleipnir-Technology/nidus-sync/http"
	"github.com/Gleipnir-Technology/nidus-sync/platform"
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

type communicationLog struct {
	Created time.Time `json:"created"`
	ID      string    `json:"id"`
	Type    string    `json:"type"`
	User    string    `json:"user"`
}
type communication struct {
	Context  []resourceStub     `json:"context"`
	Created  time.Time          `json:"created"`
	ID       string             `json:"id"`
	Log      []communicationLog `json:"log"`
	Response string             `json:"response"`
	Source   string             `json:"source"`
	Status   string             `json:"status"`
	Type     string             `json:"type"`
	URI      string             `json:"uri"`
}
type communicationStub struct {
	Created time.Time `json:"created"`
	ID      string    `json:"id"`
	Source  string    `json:"source"`
	Status  string    `json:"status"`
	Type    string    `json:"type"`
	URI     string    `json:"uri"`
}
type resourceStub struct {
	Created time.Time `json:"created"`
	Type    string    `json:"type"`
	URI     string    `json:"uri"`
}

func (res *communicationR) Get(ctx context.Context, r *http.Request, user platform.User, query QueryParams) (*communication, *nhttp.ErrorWithStatus) {
	return nil, nil
}
func (res *communicationR) List(ctx context.Context, r *http.Request, user platform.User, query QueryParams) ([]communicationStub, *nhttp.ErrorWithStatus) {
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
	public_report_id_to_report := make(map[int32]modelpublicreport.Report, 0)
	for _, pr := range public_reports {
		public_report_id_to_report[pr.ID] = pr
	}
	result := make([]communicationStub, len(comms))
	for i, comm := range comms {
		public_report, ok := public_report_id_to_report[*comm.SourceReportID]
		if !ok {
			return nil, nhttp.NewError("lookup report id %d failed", comm.SourceReportID)
		}
		c, err := res.hydrateCommunicationStub(comm, &public_report)
		if err != nil {
			return nil, err
		}
		result[i] = c
	}
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
	stub, err := res.hydrateCommunicationStub(comm, public_report)
	if err != nil {
		return communication{}, nhttp.NewError("hydrate stub: %w", err)
	}
	response, err := responseURI(*res.router, comm)
	if err != nil {
		return communication{}, nhttp.NewError("gen response URI: %w", err)
	}
	return communication{
		Created:  stub.Created,
		ID:       stub.ID,
		Response: response,
		Source:   stub.Source,
		Status:   stub.Status,
		Type:     stub.Type,
		URI:      stub.URI,
	}, nil
}
func (res *communicationR) hydrateCommunicationStub(comm modelpublic.Communication, public_report *modelpublicreport.Report) (communicationStub, *nhttp.ErrorWithStatus) {
	var err error
	source_uri := "unknown"
	type_ := "unknown"
	if comm.SourceReportID != nil && public_report != nil {
		source_uri, err = reportURI(res.router, "", public_report.PublicID)
		if err != nil {
			return communicationStub{}, nhttp.NewError("gen report URI: %w", err)
		}
		type_ = "publicreport." + public_report.ReportType.String()
	} else if comm.SourceEmailLogID != nil {
		source_uri, err = emailURI(*res.router, *comm.SourceEmailLogID)
		if err != nil {
			return communicationStub{}, nhttp.NewError("gen email URI: %w", err)
		}
		type_ = "email"
	} else if comm.SourceTextLogID != nil {
		source_uri, err = textURI(*res.router, *comm.SourceTextLogID)
		if err != nil {
			return communicationStub{}, nhttp.NewError("gen email URI: %w", err)
		}
		source_uri = "text"
	}
	/*
		response, err := responseURI(*res.router, comm)
		if err != nil {
			return communicationStub{}, nhttp.NewError("gen response URI: %w", err)
		}
	*/
	uri, err := res.router.IDToURI("communication.ByIDGet", int(comm.ID))
	if err != nil {
		return communicationStub{}, nhttp.NewError("gen comm uri: %w", err)
	}
	return communicationStub{
		Created: comm.Created,
		ID:      strconv.Itoa(int(comm.ID)),
		Source:  source_uri,
		Status:  comm.Status.String(),
		Type:    type_,
		URI:     uri,
	}, nil
}

type markFunc = func(context.Context, platform.User, int32) error

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
	if err := m(ctx, user, int32(comm_id)); err != nil {
		return communication{}, nhttp.NewError("mark communication: %w", err)
	}
	result, err := platform.CommunicationFromID(ctx, user, int64(comm_id))
	if result == nil {
		return communication{}, nhttp.NewUnauthorized("you are not authorized to modify communication %d", comm_id)
	}
	var public_report modelpublicreport.Report
	if result.SourceReportID != nil {
		comm_ids := []int64{int64(*result.SourceReportID)}
		public_reports, err := platform.PublicReportsFromIDs(ctx, comm_ids)
		if err != nil {
			return communication{}, nhttp.NewError("Get report %d: %w", *result.SourceReportID, err)
		}
		public_report = public_reports[0]
	}

	log.Info().Int("communication", comm_id).Str("status", status).Msg("Marked communication")
	return res.hydrateCommunication(*result, &public_report)
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

