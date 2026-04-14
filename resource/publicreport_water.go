package resource

import (
	"context"
	"net/http"
	"time"

	"github.com/Gleipnir-Technology/nidus-sync/db/enums"
	"github.com/Gleipnir-Technology/nidus-sync/db/models"
	"github.com/Gleipnir-Technology/nidus-sync/html"
	nhttp "github.com/Gleipnir-Technology/nidus-sync/http"
	"github.com/Gleipnir-Technology/nidus-sync/platform"
	"github.com/Gleipnir-Technology/nidus-sync/platform/types"
	"github.com/aarondl/opt/omit"
	"github.com/aarondl/opt/omitnull"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

func Water(r *router) *waterR {
	return &waterR{
		router: r,
	}
}

type waterR struct {
	router *router
}
type water struct {
	District string `json:"district"`
	PublicID string `json:"public_id"`
	URI      string `json:"uri"`
}
type waterForm struct {
	AccessComments         string         `schema:"access-comments"`
	AccessDog              bool           `schema:"access-dog"`
	AccessFence            bool           `schema:"access-fence"`
	AccessGate             bool           `schema:"access-gate"`
	AccessLocked           bool           `schema:"access-locked"`
	AccessOther            bool           `schema:"access-other"`
	Address                types.Address  `schema:"address"`
	AddressGID             string         `schema:"address-gid"`
	ClientID               uuid.UUID      `schema:"client_id" json:"client_id"`
	Comments               string         `schema:"comments"`
	HasAdult               bool           `schema:"has-adult"`
	HasBackyardPermission  bool           `schema:"backyard-permission"`
	HasLarvae              bool           `schema:"has-larvae"`
	HasPupae               bool           `schema:"has-pupae"`
	IsReporterConfidential bool           `schema:"reporter-confidential"`
	IsReporter_owner       bool           `schema:"property-ownership"`
	Location               types.Location `schema:"location"`
	OwnerEmail             string         `schema:"owner-email"`
	OwnerName              string         `schema:"owner-name"`
	OwnerPhone             string         `schema:"owner-phone"`
}

func (res *waterR) Create(ctx context.Context, r *http.Request, w waterForm) (*water, *nhttp.ErrorWithStatus) {
	user_agent := r.Header.Get("User-Agent")
	err := platform.EnsureClient(ctx, w.ClientID, user_agent)
	if err != nil {
		return nil, nhttp.NewError("Failed to ensure client: %w", err)
	}

	uploads, err := html.ExtractImageUploads(r)
	log.Info().Int("len", len(uploads)).Msg("extracted water uploads")
	if err != nil {
		return nil, nhttp.NewError("Failed to extract image uploads: %w", err)
	}

	accuracy := float32(0.0)
	if w.Location.Accuracy != nil {
		accuracy = *w.Location.Accuracy
	}
	setter_report := models.PublicreportReportSetter{
		AddressGid: omit.From(w.Address.GID),
		AddressRaw: omit.From(w.Address.Raw),
		ClientUUID: omitnull.From(w.ClientID),
		Created:    omit.From(time.Now()),
		//H3cell:       omitnull.From(geospatial.Cell.String()),
		LatlngAccuracyType:  omit.From(enums.PublicreportAccuracytypeBrowser),
		LatlngAccuracyValue: omit.From(accuracy),
		//Location: add later
		Location: omitnull.FromPtr[string](nil),
		MapZoom:  omit.From(float32(0.0)),
		//OrganizationID: omitnull.FromPtr(organization_id),
		//PublicID:       omit.From(public_id),
		ReporterEmail:       omit.From(""),
		ReporterName:        omit.From(""),
		ReporterPhone:       omit.From(""),
		ReporterPhoneCanSMS: omit.From(true),
		ReportType:          omit.From(enums.PublicreportReporttypeWater),
		Status:              omit.From(enums.PublicreportReportstatustypeReported),
	}
	setter_water := models.PublicreportWaterSetter{
		AccessComments:         omit.From(w.AccessComments),
		AccessDog:              omit.From(w.AccessDog),
		AccessFence:            omit.From(w.AccessFence),
		AccessGate:             omit.From(w.AccessGate),
		AccessLocked:           omit.From(w.AccessLocked),
		AccessOther:            omit.From(w.AccessOther),
		Comments:               omit.From(w.Comments),
		HasAdult:               omit.From(w.HasAdult),
		HasBackyardPermission:  omit.From(w.HasBackyardPermission),
		HasLarvae:              omit.From(w.HasLarvae),
		HasPupae:               omit.From(w.HasPupae),
		IsReporterConfidential: omit.From(w.IsReporterConfidential),
		IsReporterOwner:        omit.From(w.IsReporter_owner),
		OwnerEmail:             omit.From(w.OwnerEmail),
		OwnerName:              omit.From(w.OwnerName),
		OwnerPhone:             omit.From(w.OwnerPhone),
		//ReportID               omit.Val[int32]
	}
	report, err := platform.PublicReportWaterCreate(ctx, setter_report, setter_water, w.Location, w.Address, uploads)
	if err != nil {
		return nil, nhttp.NewError("Failed to save new report: %w", err)
	}
	uri, err := res.router.IDStrToURI("publicreport.ByIDGet", report.PublicID)
	if err != nil {
		return nil, nhttp.NewError("generate uri: %w", err)
	}
	district_uri, err := res.router.IDToURI("district.ByIDGet", int(report.OrganizationID))
	if err != nil {
		return nil, nhttp.NewError("generate district uri: %w", err)
	}
	return &water{
		District: district_uri,
		PublicID: report.PublicID,
		URI:      uri,
	}, nil
}
