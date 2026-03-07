package api

import (
	"context"
	"net/http"
	"time"

	"github.com/Gleipnir-Technology/bob"
	"github.com/Gleipnir-Technology/bob/dialect/psql"
	"github.com/Gleipnir-Technology/bob/dialect/psql/sm"
	"github.com/Gleipnir-Technology/nidus-sync/config"
	"github.com/Gleipnir-Technology/nidus-sync/db"
	"github.com/Gleipnir-Technology/nidus-sync/db/models"
	nhttp "github.com/Gleipnir-Technology/nidus-sync/http"
	"github.com/google/uuid"
	//"github.com/Gleipnir-Technology/nidus-sync/platform"
	//"github.com/aarondl/opt/null"
	"github.com/stephenafamo/scan"
)

type reporter struct {
	HasEmail bool   `json:"has_email"`
	HasPhone bool   `json:"has_phone"`
	Name     string `json:"name"`
}
type publicReport struct {
	AdditionalInfo      string   `json:"additional_info"`
	Address             Address  `json:"address"`
	Duration            string   `json:"duration"`
	Images              []string `json:"images"`
	IsLocationBackyard  bool     `json:"is_location_backyard"`
	IsLocationFrontyard bool     `json:"is_location_frontyard"`
	IsLocationGarden    bool     `json:"is_location_garden"`
	IsLocationOther     bool     `json:"is_location_other"`
	IsLocationPool      bool     `json:"is_location_pool"`
	Location            Location `json:"location"`
	Reporter            reporter `json:"reporter"`
	SourceContainer     bool     `json:"source_container"`
	SourceDescription   string   `json:"source_description"`
	SourceGutter        bool     `json:"source_gutter"`
	SourceStagnant      bool     `json:"source_stagnant"`
	TODDay              bool     `json:"time_of_day_day"`
	TODEarly            bool     `json:"time_of_day_early"`
	TODEvening          bool     `json:"time_of_day_evening"`
	TODNight            bool     `json:"time_of_day_night"`
}
type communication struct {
	Created      time.Time    `json:"created"`
	ID           string       `json:"id"`
	PublicReport publicReport `json:"public_report"`
	Type         string       `json:"type"`
}
type contentListCommunication struct {
	Communications []communication `json:"communications"`
}

func listCommunication(ctx context.Context, r *http.Request, org *models.Organization, user *models.User, query queryParams) (*contentListCommunication, *nhttp.ErrorWithStatus) {
	type _Report struct {
		AdditionalInfo      string    `db:"additional_info"`
		AddressCountry      string    `db:"address_country" `
		AddressPlace        string    `db:"address_place" `
		AddressPostcode     string    `db:"address_postcode" `
		AddressRegion       string    `db:"address_region" `
		AddressStreet       string    `db:"address_street" `
		Created             time.Time `db:"created" `
		Duration            string    `db:"duration" `
		IsLocationBackyard  bool      `db:"is_location_backyard"`
		IsLocationFrontyard bool      `db:"is_location_frontyard"`
		IsLocationGarden    bool      `db:"is_location_garden"`
		IsLocationOther     bool      `db:"is_location_other"`
		IsLocationPool      bool      `db:"is_location_pool"`
		Latitude            float64   `db:"latitude"`
		Longitude           float64   `db:"longitude"`
		PublicID            string    `db:"public_id" `
		ReporterEmail       *string   `db:"reporter_email" `
		ReporterName        *string   `db:"reporter_name" `
		ReporterPhone       *string   `db:"reporter_phone" `
		SourceContainer     bool      `db:"source_container"`
		SourceDescription   string    `db:"source_description"`
		SourceGutter        bool      `db:"source_gutter"`
		SourceStagnant      bool      `db:"source_stagnant"`
		TODDay              bool      `db:"tod_day"`
		TODEarly            bool      `db:"tod_early"`
		TODEvening          bool      `db:"tod_evening"`
		TODNight            bool      `db:"tod_night"`
	}
	reports, err := bob.All(ctx, db.PGInstance.BobDB, psql.Select(
		sm.Columns(
			"additional_info",
			"address_country",
			"address_place",
			"address_postcode",
			"address_region",
			"address_street",
			"created",
			"duration",
			"is_location_backyard",
			"is_location_frontyard",
			"is_location_garden",
			"is_location_other",
			"is_location_pool",
			"ST_Y(location::geometry::geometry(point, 4326)) AS latitude",
			"ST_X(location::geometry::geometry(point, 4326)) AS longitude",
			"public_id",
			"reporter_email",
			"reporter_phone",
			"reporter_name",
			"source_container",
			"source_description",
			"source_gutter",
			"source_stagnant",
			"tod_day",
			"tod_early",
			"tod_evening",
			"tod_night",
		),
		sm.From("publicreport.nuisance"),
		sm.Where(psql.Quote("publicreport", "nuisance", "organization_id").EQ(psql.Arg(org.ID))),
	), scan.StructMapper[_Report]())
	if err != nil {
		return nil, nhttp.NewError("get reports: %w", err)
	}
	type _Row struct {
		PublicID    string    `db:"nuisance_public_id"`
		StorageUUID uuid.UUID `db:"storage_uuid"`
	}
	rows, err := bob.All(ctx, db.PGInstance.BobDB, psql.Select(
		sm.Columns(
			"n.public_id AS nuisance_public_id",
			"i.storage_uuid AS storage_uuid",
		),
		sm.From("publicreport.nuisance").As("n"),
		sm.InnerJoin("publicreport.nuisance_image").As("ni").OnEQ(
			psql.Quote("n", "id"),
			psql.Quote("ni", "nuisance_id"),
		),
		sm.InnerJoin("publicreport.image").As("i").OnEQ(
			psql.Quote("ni", "image_id"),
			psql.Quote("i", "id"),
		),
		sm.Where(psql.Quote("n", "organization_id").EQ(psql.Arg(org.ID))),
	), scan.StructMapper[_Row]())
	if err != nil {
		return nil, nhttp.NewError("get images: %w")
	}
	id_to_images := make(map[string][]uuid.UUID, len(reports))
	for _, row := range rows {
		r, ok := id_to_images[row.PublicID]
		if !ok {
			r = make([]uuid.UUID, 0)
		}
		r = append(r, row.StorageUUID)
		id_to_images[row.PublicID] = r
	}
	comms := make([]communication, len(reports))
	for i, report := range reports {
		name := ""
		if report.ReporterName != nil {
			name = *report.ReporterName
		}
		comms[i] = communication{
			Created: report.Created,
			ID:      report.PublicID,
			PublicReport: publicReport{
				Address: Address{
					Country:  report.AddressCountry,
					Locality: report.AddressPlace,
					//Number: report.Address
					PostalCode: report.AddressPostcode,
					Region:     report.AddressRegion,
					Street:     report.AddressStreet,
				},
				AdditionalInfo:      report.AdditionalInfo,
				Duration:            report.Duration,
				Images:              toImageURLs(id_to_images, report.PublicID),
				IsLocationBackyard:  report.IsLocationBackyard,
				IsLocationFrontyard: report.IsLocationFrontyard,
				IsLocationGarden:    report.IsLocationGarden,
				IsLocationOther:     report.IsLocationOther,
				IsLocationPool:      report.IsLocationPool,
				Location: Location{
					Latitude:  report.Latitude,
					Longitude: report.Longitude,
				},
				Reporter: reporter{
					Name:     name,
					HasEmail: report.ReporterEmail != nil,
					HasPhone: report.ReporterPhone != nil,
				},
				SourceContainer:   report.SourceContainer,
				SourceDescription: report.SourceDescription,
				SourceGutter:      report.SourceGutter,
				SourceStagnant:    report.SourceStagnant,
				TODDay:            report.TODDay,
				TODEarly:          report.TODEarly,
				TODEvening:        report.TODEvening,
				TODNight:          report.TODNight,
			},
			Type: "nuisance",
		}
	}
	return &contentListCommunication{
		Communications: comms,
	}, nil
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
