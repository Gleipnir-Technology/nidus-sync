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
	"github.com/aarondl/opt/omit"
	//"github.com/aarondl/opt/omitnull"
	//"github.com/rs/zerolog/log"
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
	ID string `json:"id"`
}
type waterForm struct {
	AccessComments         string `schema:"access-comments"`
	AccessDog              bool   `schema:"access-dog"`
	AccessFence            bool   `schema:"access-fence"`
	AccessGate             bool   `schema:"access-gate"`
	AccessLocked           bool   `schema:"access-locked"`
	AccessOther            bool   `schema:"access-other"`
	AddressRaw             string `schema:"address"`
	AddressCountry         string `schema:"address-country"`
	AddressLocality        string `schema:"address-locality"`
	AddressNumber          string `schema:"address-number"`
	AddressPostalCode      string `schema:"address-postalcode"`
	AddressRegion          string `schema:"address-region"`
	AddressStreet          string `schema:"address-street"`
	Comments               string `schema:"comments"`
	HasAdult               bool   `schema:"has-adult"`
	HasBackyardPermission  bool   `schema:"backyard-permission"`
	HasLarvae              bool   `schema:"has-larvae"`
	HasPupae               bool   `schema:"has-pupae"`
	IsReporterConfidential bool   `schema:"reporter-confidential"`
	IsReporter_owner       bool   `schema:"property-ownership"`
	OwnerEmail             string `schema:"owner-email"`
	OwnerName              string `schema:"owner-name"`
	OwnerPhone             string `schema:"owner-phone"`
}

func (res *waterR) Create(ctx context.Context, r *http.Request, w waterForm) (*water, *nhttp.ErrorWithStatus) {
	latlng, err := parseLatLng(r)
	if err != nil {
		return nil, nhttp.NewError("Failed to parse lat lng for water report: %w", err)
	}

	uploads, err := html.ExtractImageUploads(r)
	if err != nil {
		return nil, nhttp.NewError("Failed to extract image uploads: %w", err)
	}

	address := platform.Address{
		Country:    w.AddressCountry,
		Locality:   w.AddressLocality,
		Number:     w.AddressNumber,
		PostalCode: w.AddressPostalCode,
		Raw:        w.AddressRaw,
		Region:     w.AddressRegion,
		Street:     w.AddressStreet,
		Unit:       "",
	}
	setter_report := models.PublicreportReportSetter{
		AddressRaw:        omit.From(address.Raw),
		AddressCountry:    omit.From(address.Country),
		AddressNumber:     omit.From(address.Number),
		AddressLocality:   omit.From(address.Locality),
		AddressPostalCode: omit.From(address.PostalCode),
		AddressRegion:     omit.From(address.Region),
		AddressStreet:     omit.From(address.Street),
		Created:           omit.From(time.Now()),
		//H3cell:       omitnull.From(geospatial.Cell.String()),
		LatlngAccuracyType:  omit.From(latlng.AccuracyType),
		LatlngAccuracyValue: omit.From(float32(latlng.AccuracyValue)),
		//Location: add later
		MapZoom: omit.From(latlng.MapZoom),
		//OrganizationID: omitnull.FromPtr(organization_id),
		//PublicID:       omit.From(public_id),
		ReporterEmail: omit.From(""),
		ReporterName:  omit.From(""),
		ReporterPhone: omit.From(""),
		ReportType:    omit.From(enums.PublicreportReporttypeWater),
		Status:        omit.From(enums.PublicreportReportstatustypeReported),
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
	report, err := platform.ReportWaterCreate(ctx, setter_report, setter_water, latlng, address, uploads)
	if err != nil {
		return nil, nhttp.NewError("Failed to save new report: %w", err)
	}
	return &water{
		ID: report.PublicID,
	}, nil
}
