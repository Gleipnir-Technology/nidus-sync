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
	Address  Address  `json:"address"`
	Images   []string `json:"images"`
	Location Location `json:"location"`
	Reporter reporter `json:"reporter"`
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
		AddressCountry  string    `db:"address_country" `
		AddressPlace    string    `db:"address_place" `
		AddressPostcode string    `db:"address_postcode" `
		AddressRegion   string    `db:"address_region" `
		AddressStreet   string    `db:"address_street" `
		Created         time.Time `db:"created" `
		Latitude        float64   `db:"latitude"`
		Longitude       float64   `db:"longitude"`
		PublicID        string    `db:"public_id" `
		ReporterEmail   *string   `db:"reporter_email" `
		ReporterName    *string   `db:"reporter_name" `
		ReporterPhone   *string   `db:"reporter_phone" `
	}
	reports, err := bob.All(ctx, db.PGInstance.BobDB, psql.Select(
		sm.Columns(
			"address_country",
			"address_place",
			"address_postcode",
			"address_region",
			"address_street",
			"created",
			"ST_Y(location::geometry::geometry(point, 4326)) AS latitude",
			"ST_X(location::geometry::geometry(point, 4326)) AS longitude",
			"public_id",
			"reporter_email",
			"reporter_phone",
			"reporter_name",
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
			PublicReport: publicReport{
				Address: Address{
					Country:  report.AddressCountry,
					Locality: report.AddressPlace,
					//Number: report.Address
					PostalCode: report.AddressPostcode,
					Region:     report.AddressRegion,
					Street:     report.AddressStreet,
				},
				Images: toImageURLs(id_to_images, report.PublicID),
				Location: Location{
					Latitude:  report.Latitude,
					Longitude: report.Longitude,
				},
				Reporter: reporter{
					Name:     name,
					HasEmail: report.ReporterEmail != nil,
					HasPhone: report.ReporterPhone != nil,
				},
			},
			Created: report.Created,
			ID:      report.PublicID,
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
	for _, u := range uuids {
		urls = append(urls, config.MakeURLNidus("/api/image/%s", u.String()))
	}
	return urls
}
