package resource

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	modelpublic "github.com/Gleipnir-Technology/nidus-sync/db/gen/nidus-sync/public/model"
	modelpublicreport "github.com/Gleipnir-Technology/nidus-sync/db/gen/nidus-sync/publicreport/model"
	nhttp "github.com/Gleipnir-Technology/nidus-sync/http"
	"github.com/Gleipnir-Technology/nidus-sync/platform"
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

type communicationLogType string
type communicationType string

const (
	communicationTypeUnknown          communicationType = "unknown"
	communicationTypeReportCompliance                   = "publicreport.compliance"
	communicationTypeReportNuisance                     = "publicreport.nuisance"
	communicationTypeReportWater                        = "publicreport.water"
)

type communicationLog struct {
	Created time.Time            `json:"created"`
	ID      string               `json:"id"`
	Type    communicationLogType `json:"type"`
	User    string               `json:"user"`
}
type communication struct {
	Created  time.Time          `json:"created"`
	ID       string             `json:"id"`
	Log      []communicationLog `json:"log"`
	Related  []resourceStub     `json:"related"`
	Response string             `json:"response"`
	Source   string             `json:"source"`
	Status   string             `json:"status"`
	Type     communicationType  `json:"type"`
	URI      string             `json:"uri"`
}
type resourceType string

const (
	resourceTypeUnknown          resourceType = "unknown"
	resourceTypeEmail                         = "email"
	resourceTypeReportCompliance              = "publicreport.compliance"
	resourceTypeReportNuisance                = "publicreport.nuisance"
	resourceTypeReportWater                   = "publicreport.water"
	resourceTypeText                          = "text"
)

type resourceStub struct {
	Created time.Time    `json:"created"`
	Type    resourceType `json:"type"`
	URI     string       `json:"uri"`
}

func (res *communicationR) Get(ctx context.Context, r *http.Request, user platform.User, query QueryParams) (communication, *nhttp.ErrorWithStatus) {
	comm_id, error_with_status := res.router.IDFromMux(r)
	if error_with_status != nil {
		return communication{}, error_with_status
	}
	return res.hydratedCommunicationFromID(ctx, user, int32(comm_id))
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
	public_report_id_to_report := make(map[int32]modelpublicreport.Report, 0)
	for _, pr := range public_reports {
		public_report_id_to_report[pr.ID] = pr
	}
	result := make([]communication, len(comms))
	for i, comm := range comms {
		public_report, ok := public_report_id_to_report[*comm.SourceReportID]
		if !ok {
			return nil, nhttp.NewError("lookup report id %d failed", comm.SourceReportID)
		}
		related_records, err := platform.CommunicationRelatedRecords(ctx, user, &comm)
		if err != nil {
			return nil, nhttp.NewError("related records: %w", err)
		}
		c, e := res.hydrateCommunication(comm, related_records, &public_report)
		if e != nil {
			return nil, e
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
func (res *communicationR) hydrateCommunication(comm modelpublic.Communication, related_records []platform.RelatedRecord, public_report *modelpublicreport.Report) (communication, *nhttp.ErrorWithStatus) {
	var err error
	source_uri := "unknown"
	type_ := communicationTypeUnknown
	if comm.SourceReportID != nil && public_report != nil {
		source_uri, err = reportURI(res.router, "", public_report.PublicID)
		if err != nil {
			return communication{}, nhttp.NewError("gen report URI: %w", err)
		}
		switch public_report.ReportType {
		case modelpublicreport.Reporttype_Compliance:
			type_ = communicationTypeReportCompliance
		case modelpublicreport.Reporttype_Nuisance:
			type_ = communicationTypeReportNuisance
		case modelpublicreport.Reporttype_Water:
			type_ = communicationTypeReportWater
		default:
			type_ = communicationTypeUnknown
		}
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
	uri, err := res.router.IDToURI("communication.ByIDGet", int(comm.ID))
	if err != nil {
		return communication{}, nhttp.NewError("gen comm uri: %w", err)
	}
	related := make([]resourceStub, len(related_records))
	for i, rr := range related_records {
		var uri string
		var r_type resourceType
		switch rr.Type {
		case platform.RelatedRecordTypeEmail:
			r_type = resourceTypeEmail
			uri, err = res.router.IDStrToURI("email.ByIDGet", rr.ID)
		case platform.RelatedRecordTypeReportCompliance:
			r_type = resourceTypeReportCompliance
			uri, err = res.router.IDStrToURI("publicreport.compliance.ByIDGet", rr.ID)
		case platform.RelatedRecordTypeReportNuisance:
			r_type = resourceTypeReportNuisance
			uri, err = res.router.IDStrToURI("publicreport.nuisance.ByIDGet", rr.ID)
		case platform.RelatedRecordTypeReportWater:
			r_type = resourceTypeReportWater
			uri, err = res.router.IDStrToURI("publicreport.water.ByIDGet", rr.ID)
		case platform.RelatedRecordTypeText:
			r_type = resourceTypeText
			uri, err = res.router.IDStrToURI("text.ByIDGet", rr.ID)
		default:
			r_type = resourceTypeUnknown
			err = fmt.Errorf("unrecognized related record type '%s'", rr.Type)
		}
		if err != nil {
			return communication{}, nhttp.NewError("related record hydration: %w", err)
		}
		related[i] = resourceStub{
			Created: rr.Created,
			Type:    r_type,
			URI:     uri,
		}
	}
	response, err := responseURI(*res.router, comm)
	if err != nil {
		return communication{}, nhttp.NewError("gen response URI: %w", err)
	}
	return communication{
		Created:  comm.Created,
		ID:       strconv.Itoa(int(comm.ID)),
		Related:  related,
		Response: response,
		Source:   source_uri,
		Status:   comm.Status.String(),
		Type:     type_,
		URI:      uri,
	}, nil
}

func (res *communicationR) hydratedCommunicationFromID(ctx context.Context, user platform.User, comm_id int32) (communication, *nhttp.ErrorWithStatus) {
	result, err := platform.CommunicationFromID(ctx, user, int64(comm_id))
	if result == nil {
		return communication{}, nhttp.NewUnauthorized("you are not authorized to modify communication %d", comm_id)
	}
	comm, err := platform.CommunicationFromID(ctx, user, int64(comm_id))
	if err != nil {
		return communication{}, nhttp.NewError("comm from ID: %w", err)
	}
	var public_report modelpublicreport.Report
	if comm.SourceReportID != nil {
		comm_ids := []int64{int64(*comm.SourceReportID)}
		public_reports, err := platform.PublicReportsFromIDs(ctx, comm_ids)
		if err != nil {
			return communication{}, nhttp.NewError("Get report %d: %w", *comm.SourceReportID, err)
		}
		public_report = public_reports[0]
	}
	related_records, err := platform.CommunicationRelatedRecords(ctx, user, comm)
	if err != nil {
		return communication{}, nhttp.NewError("related records: %w", err)
	}

	return res.hydrateCommunication(*comm, related_records, &public_report)
}

type markFunc = func(context.Context, platform.User, int32) error

func (res *communicationR) markCommunication(ctx context.Context, r *http.Request, user platform.User, status string, m markFunc) (communication, *nhttp.ErrorWithStatus) {
	comm_id, err_with_status := res.router.IDFromMux(r)
	if err_with_status != nil {
		return communication{}, err_with_status
	}
	err := m(ctx, user, int32(comm_id))
	if err != nil {
		return communication{}, nhttp.NewError("failed to mark %d: %w", comm_id, err)
	}
	return res.hydratedCommunicationFromID(ctx, user, int32(comm_id))
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
