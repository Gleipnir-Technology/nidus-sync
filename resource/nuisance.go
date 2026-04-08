package resource

import (
	"context"
	"fmt"
	"net/http"
	"slices"
	"strconv"
	"time"

	"github.com/Gleipnir-Technology/nidus-sync/db/enums"
	"github.com/Gleipnir-Technology/nidus-sync/db/models"
	"github.com/Gleipnir-Technology/nidus-sync/html"
	nhttp "github.com/Gleipnir-Technology/nidus-sync/http"
	"github.com/Gleipnir-Technology/nidus-sync/platform"
	"github.com/aarondl/opt/omit"
	"github.com/aarondl/opt/omitnull"
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
	ID string `json:"id"`
}
type nuisanceForm struct {
	AdditionalInfo      string   `schema:"additional-info"`
	AddressGID          string   `schema:"address-gid"`
	Address             string   `schema:"address"`
	Duration            string   `schema:"duration"`
	Latitude            string   `schema:"latitude"`
	Longitude           string   `schema:"longitude"`
	LatlngAccuracyType  string   `schema:"latlng-accuracy-type"`
	LatlngAccuracyValue string   `schema:"latlng-accuracy-value"`
	MapZoom             string   `schema:"map-zoom"`
	SourceStagnant      bool     `schema:"source-stagnant"`
	SourceContainer     bool     `schema:"source-container"`
	SourceDescription   string   `schema:"source-description"`
	SourceGutters       bool     `schema:"source-gutters"`
	SourceLocations     []string `schema:"source-location"`
	TODEarly            bool     `schema:"tod-early"`
	TODDay              bool     `schema:"tod-day"`
	TODEvening          bool     `schema:"tod-evening"`
	TODNight            bool     `schema:"tod-night"`
}

func parseLatLng(r *http.Request) (platform.LatLng, error) {
	result := platform.LatLng{
		AccuracyType:  enums.PublicreportAccuracytypeNone,
		AccuracyValue: 0.0,
		Latitude:      nil,
		Longitude:     nil,
		MapZoom:       0.0,
	}
	latitude_str := r.FormValue("latitude")
	longitude_str := r.FormValue("longitude")
	latlng_accuracy_type_str := r.PostFormValue("latlng-accuracy-type")
	latlng_accuracy_value_str := r.PostFormValue("latlng-accuracy-value")
	map_zoom_str := r.PostFormValue("map-zoom")

	var err error
	if latlng_accuracy_type_str != "" {
		err := result.AccuracyType.Scan(latlng_accuracy_type_str)
		if err != nil {
			return result, fmt.Errorf("Failed to parse accuracy type '%s': %w", latlng_accuracy_type_str, err)
		}
	}
	if latlng_accuracy_value_str != "" {
		var t float64
		t, err = strconv.ParseFloat(latlng_accuracy_value_str, 32)
		if err != nil {
			return result, fmt.Errorf("Failed to parse latlng_accuracy_value '%s': %w", latlng_accuracy_value_str, err)
		}
		result.AccuracyValue = float64(t)
	}

	if latitude_str != "" {
		var t float64
		t, err = strconv.ParseFloat(latitude_str, 64)
		if err != nil {
			return result, fmt.Errorf("Failed to parse latitude '%s': %w", latitude_str, err)
		}
		result.Latitude = &t
	}
	if longitude_str != "" {
		var t float64
		t, err := strconv.ParseFloat(longitude_str, 64)
		if err != nil {
			return result, fmt.Errorf("Failed to parse longitude '%s': %w", longitude_str, err)
		}
		result.Longitude = &t
	}

	if map_zoom_str != "" {
		var t float64
		t, err = strconv.ParseFloat(map_zoom_str, 32)
		if err != nil {
			return result, fmt.Errorf("Failed to parse map_zoom_str '%s': %w", map_zoom_str, err)
		} else {
			result.MapZoom = float32(t)
		}
	}
	return result, nil
}

func (res *nuisanceR) Create(ctx context.Context, r *http.Request, n nuisanceForm) (*nuisance, *nhttp.ErrorWithStatus) {
	duration := enums.PublicreportNuisancedurationtypeNone
	is_location_frontyard := slices.Contains(n.SourceLocations, "frontyard")
	is_location_backyard := slices.Contains(n.SourceLocations, "backyard")
	is_location_garden := slices.Contains(n.SourceLocations, "garden")
	is_location_pool := slices.Contains(n.SourceLocations, "pool-area")
	is_location_other := slices.Contains(n.SourceLocations, "other")

	latlng, err := parseLatLng(r)

	err = duration.Scan(n.Duration)
	if err != nil {
		log.Warn().Err(err).Str("duration_str", n.Duration).Msg("Failed to interpret 'duration'")
	}

	uploads, err := html.ExtractImageUploads(r)
	log.Info().Int("len", len(uploads)).Msg("extracted uploads")
	if err != nil {
		return nil, nhttp.NewError("Failed to extract image uploads: %w", err)
	}
	address := platform.Address{
		GID: n.AddressGID,
		Raw: n.Address,
	}
	setter_report := models.PublicreportReportSetter{
		//AddressID:              omitnull.From(latlng.Cell.String()),
		AddressCountry:    omit.From(""),
		AddressGid:        omit.From(address.GID),
		AddressNumber:     omit.From(""),
		AddressLocality:   omit.From(""),
		AddressPostalCode: omit.From(""),
		AddressRaw:        omit.From(address.Raw),
		AddressRegion:     omit.From(""),
		AddressStreet:     omit.From(""),
		Created:           omit.From(time.Now()),
		//H3cell:              omitnull.From(latlng.Cell.String()),
		LatlngAccuracyType:  omit.From(latlng.AccuracyType),
		LatlngAccuracyValue: omit.From(float32(latlng.AccuracyValue)),
		//Location: omitnull.From(fmt.Sprintf("ST_GeometryFromText(Point(%s %s))", longitude, latitude)),
		Location: omitnull.FromPtr[string](nil),
		MapZoom:  omit.From(latlng.MapZoom),
		//OrganizationID:    omitnull.FromPtr(organization_id),
		//PublicID:          omit.From(public_id),
		ReporterEmail: omit.From(""),
		ReporterName:  omit.From(""),
		ReporterPhone: omit.From(""),
		ReportType:    omit.From(enums.PublicreportReporttypeNuisance),
		Status:        omit.From(enums.PublicreportReportstatustypeReported),
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
	report, err := platform.ReportNuisanceCreate(ctx, setter_report, setter_nuisance, latlng, address, uploads)
	if err != nil {
		return nil, nhttp.NewError("create nuisance report: %w", err)
	}
	return &nuisance{
		ID: report.PublicID,
	}, nil
}
