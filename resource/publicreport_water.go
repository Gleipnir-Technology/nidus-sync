package resource

import (
	"context"
	"net/http"
	"time"

	modelpublicreport "github.com/Gleipnir-Technology/nidus-sync/db/gen/nidus-sync/publicreport/model"
	//tablepublicreport "github.com/Gleipnir-Technology/nidus-sync/db/gen/nidus-sync/publicreport/table"
	//querypublicreport "github.com/Gleipnir-Technology/nidus-sync/db/query/publicreport"
	"github.com/Gleipnir-Technology/nidus-sync/html"
	nhttp "github.com/Gleipnir-Technology/nidus-sync/http"
	"github.com/Gleipnir-Technology/nidus-sync/platform"
	"github.com/Gleipnir-Technology/nidus-sync/platform/types"
	"github.com/aarondl/opt/omit"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
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
	AccessComments         string                                           `schema:"access-comments"`
	AccessDog              bool                                             `schema:"access-dog"`
	AccessFence            bool                                             `schema:"access-fence"`
	AccessGate             bool                                             `schema:"access-gate"`
	AccessLocked           bool                                             `schema:"access-locked"`
	AccessOther            bool                                             `schema:"access-other"`
	Address                types.Address                                    `schema:"address"`
	AddressGID             string                                           `schema:"address-gid"`
	ClientID               uuid.UUID                                        `schema:"client_id" json:"client_id"`
	Comments               string                                           `schema:"comments"`
	Duration               omit.Val[modelpublicreport.Nuisancedurationtype] `schema:"duration"`
	HasAdult               bool                                             `schema:"has-adult"`
	HasBackyardPermission  bool                                             `schema:"backyard-permission"`
	HasLarvae              bool                                             `schema:"has-larvae"`
	HasPupae               bool                                             `schema:"has-pupae"`
	IsReporterConfidential bool                                             `schema:"reporter-confidential"`
	IsReporter_owner       bool                                             `schema:"property-ownership"`
	Location               types.Location                                   `schema:"location"`
	OwnerEmail             string                                           `schema:"owner-email"`
	OwnerName              string                                           `schema:"owner-name"`
	OwnerPhone             string                                           `schema:"owner-phone"`
}

func (res *waterR) ByID(ctx context.Context, r *http.Request, u platform.User, query QueryParams) (*types.PublicReportWater, *nhttp.ErrorWithStatus) {
	return res.byID(ctx, r, false)
}
func (res *waterR) ByIDPublic(ctx context.Context, r *http.Request, query QueryParams) (*types.PublicReportWater, *nhttp.ErrorWithStatus) {
	return res.byID(ctx, r, true)
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
	setter_report := modelpublicreport.Report{
		//AddressID:              omitnull.From(...),
		AddressGid: w.Address.GID,
		AddressRaw: w.Address.Raw,
		ClientUUID: &w.ClientID,
		Created:    time.Now(),
		//H3cell:              omitnull.From(latlng.Cell.String()),
		LatlngAccuracyType:  modelpublicreport.Accuracytype_Browser,
		LatlngAccuracyValue: accuracy,
		//Location: omitnull.From(fmt.Sprintf("ST_GeometryFromText(Point(%s %s))", longitude, latitude)),
		Location: nil,
		MapZoom:  float32(0.0),
		//OrganizationID:      ,
		//PublicID:
		ReporterEmail:       "",
		ReporterName:        "",
		ReporterPhone:       "",
		ReporterPhoneCanSms: true,
		ReportType:          modelpublicreport.Reporttype_Water,
		Status:              modelpublicreport.Reportstatustype_Reported,
	}
	setter_water := modelpublicreport.Water{
		AccessComments:         w.AccessComments,
		AccessDog:              w.AccessDog,
		AccessFence:            w.AccessFence,
		AccessGate:             w.AccessGate,
		AccessLocked:           w.AccessLocked,
		AccessOther:            w.AccessOther,
		Comments:               w.Comments,
		Duration:               w.Duration.GetOr(modelpublicreport.Nuisancedurationtype_None),
		HasAdult:               w.HasAdult,
		HasBackyardPermission:  w.HasBackyardPermission,
		HasLarvae:              w.HasLarvae,
		HasPupae:               w.HasPupae,
		IsReporterConfidential: w.IsReporterConfidential,
		IsReporterOwner:        w.IsReporter_owner,
		OwnerEmail:             w.OwnerEmail,
		OwnerName:              w.OwnerName,
		OwnerPhone:             w.OwnerPhone,
		//ReportID               omit.Val[int32]
	}
	report, err := platform.PublicReportWaterCreate(ctx, setter_report, setter_water, w.Location, w.Address, uploads)
	if err != nil {
		return nil, nhttp.NewError("Failed to save new report: %w", err)
	}
	uri, err := res.router.IDStrToURI("publicreport.ByIDGetPublic", report.PublicID)
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
func (res *waterR) byID(ctx context.Context, r *http.Request, is_public bool) (*types.PublicReportWater, *nhttp.ErrorWithStatus) {
	vars := mux.Vars(r)
	public_id := vars["id"]
	if public_id == "" {
		return nil, nhttp.NewBadRequest("You must provid an ID")
	}
	report, err := platform.PublicReportByIDWater(ctx, public_id, is_public)
	if err != nil {
		return nil, nhttp.NewError("get report: %w", err)
	}
	populateDistrictURI(&report.PublicReport, res.router)
	populateReportURI(&report.PublicReport, res.router, is_public)
	return report, nil
}
