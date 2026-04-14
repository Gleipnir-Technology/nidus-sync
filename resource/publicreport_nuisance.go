package resource

import (
	"context"
	"net/http"
	"slices"
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
	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"
)

func Nuisance(r *router) *nuisanceR {
	return &nuisanceR{
		router: r,
	}
}

type nuisanceR struct {
	router *router
}
type nuisance struct {
	District string `json:"district"`
	PublicID string `json:"public_id"`
	URI      string `json:"uri"`
}
type nuisanceForm struct {
	Address           types.Address  `schema:"address"`
	AdditionalInfo    string         `schema:"additional-info"`
	ClientID          uuid.UUID      `schema:"client_id" json:"client_id"`
	Duration          string         `schema:"duration"`
	Location          types.Location `schema:"location"`
	MapZoom           string         `schema:"map-zoom"`
	SourceStagnant    bool           `schema:"source-stagnant"`
	SourceContainer   bool           `schema:"source-container"`
	SourceDescription string         `schema:"source-description"`
	SourceGutters     bool           `schema:"source-gutters"`
	SourceLocations   []string       `schema:"source-location"`
	TODEarly          bool           `schema:"tod-early"`
	TODDay            bool           `schema:"tod-day"`
	TODEvening        bool           `schema:"tod-evening"`
	TODNight          bool           `schema:"tod-night"`
}

func (res *nuisanceR) ByID(ctx context.Context, r *http.Request, query QueryParams) (*types.PublicReportNuisance, *nhttp.ErrorWithStatus) {
	vars := mux.Vars(r)
	public_id := vars["id"]
	if public_id == "" {
		return nil, nhttp.NewBadRequest("You must provid an ID")
	}
	report, err := platform.PublicreportByIDNuisance(ctx, public_id)
	if err != nil {
		return nil, nhttp.NewError("get report: %w", err)
	}
	populateDistrictURI(&report.PublicReport, res.router)
	populateReportURI(&report.PublicReport, res.router)
	return report, nil
}
func (res *nuisanceR) Create(ctx context.Context, r *http.Request, n nuisanceForm) (*nuisance, *nhttp.ErrorWithStatus) {
	user_agent := r.Header.Get("User-Agent")
	err := platform.EnsureClient(ctx, n.ClientID, user_agent)
	if err != nil {
		return nil, nhttp.NewError("Failed to ensure client: %w", err)
	}
	duration := enums.PublicreportNuisancedurationtypeNone
	is_location_frontyard := slices.Contains(n.SourceLocations, "frontyard")
	is_location_backyard := slices.Contains(n.SourceLocations, "backyard")
	is_location_garden := slices.Contains(n.SourceLocations, "garden")
	is_location_pool := slices.Contains(n.SourceLocations, "pool-area")
	is_location_other := slices.Contains(n.SourceLocations, "other")

	err = duration.Scan(n.Duration)
	if err != nil {
		log.Warn().Err(err).Str("duration_str", n.Duration).Msg("Failed to interpret 'duration'")
	}

	uploads, err := html.ExtractImageUploads(r)
	log.Info().Int("len", len(uploads)).Msg("extracted nuisance uploads")
	if err != nil {
		return nil, nhttp.NewError("Failed to extract image uploads: %w", err)
	}
	accuracy := float32(0.0)
	if n.Location.Accuracy != nil {
		accuracy = *n.Location.Accuracy
	}
	setter_report := models.PublicreportReportSetter{
		//AddressID:              omitnull.From(latlng.Cell.String()),
		AddressGid: omit.From(n.Address.GID),
		AddressRaw: omit.From(n.Address.Raw),
		ClientUUID: omitnull.From(n.ClientID),
		Created:    omit.From(time.Now()),
		//H3cell:              omitnull.From(latlng.Cell.String()),
		LatlngAccuracyType:  omit.From(enums.PublicreportAccuracytypeBrowser),
		LatlngAccuracyValue: omit.From(accuracy),
		//Location: omitnull.From(fmt.Sprintf("ST_GeometryFromText(Point(%s %s))", longitude, latitude)),
		Location: omitnull.FromPtr[string](nil),
		MapZoom:  omit.From(float32(0.0)),
		//OrganizationID:    omitnull.FromPtr(organization_id),
		//PublicID:          omit.From(public_id),
		ReporterEmail:       omit.From(""),
		ReporterName:        omit.From(""),
		ReporterPhone:       omit.From(""),
		ReporterPhoneCanSMS: omit.From(true),
		ReportType:          omit.From(enums.PublicreportReporttypeNuisance),
		Status:              omit.From(enums.PublicreportReportstatustypeReported),
	}
	setter_nuisance := models.PublicreportNuisanceSetter{
		AdditionalInfo:      omit.From(n.AdditionalInfo),
		Duration:            omit.From(duration),
		IsLocationBackyard:  omit.From(is_location_backyard),
		IsLocationFrontyard: omit.From(is_location_frontyard),
		IsLocationGarden:    omit.From(is_location_garden),
		IsLocationOther:     omit.From(is_location_other),
		IsLocationPool:      omit.From(is_location_pool),
		//ReportID            omit.Val[int32]
		SourceContainer:   omit.From(n.SourceContainer),
		SourceDescription: omit.From(n.SourceDescription),
		SourceGutter:      omit.From(n.SourceGutters),
		SourceStagnant:    omit.From(n.SourceStagnant),
		TodDay:            omit.From(n.TODDay),
		TodEarly:          omit.From(n.TODEarly),
		TodEvening:        omit.From(n.TODEvening),
		TodNight:          omit.From(n.TODNight),
	}
	report, err := platform.PublicReportNuisanceCreate(ctx, setter_report, setter_nuisance, n.Location, n.Address, uploads)
	if err != nil {
		return nil, nhttp.NewError("create nuisance report: %w", err)
	}
	uri, err := res.router.IDStrToURI("publicreport.ByIDGet", report.PublicID)
	if err != nil {
		return nil, nhttp.NewError("generate uri: %w", err)
	}
	district_uri, err := res.router.IDToURI("district.ByIDGet", int(report.OrganizationID))
	if err != nil {
		return nil, nhttp.NewError("generate district uri: %w", err)
	}
	return &nuisance{
		District: district_uri,
		PublicID: report.PublicID,
		URI:      uri,
	}, nil
}
