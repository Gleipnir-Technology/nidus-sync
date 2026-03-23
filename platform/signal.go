package platform

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/Gleipnir-Technology/bob"
	"github.com/Gleipnir-Technology/bob/dialect/psql"
	"github.com/Gleipnir-Technology/bob/dialect/psql/sm"
	"github.com/Gleipnir-Technology/bob/dialect/psql/um"
	"github.com/Gleipnir-Technology/nidus-sync/db"
	"github.com/Gleipnir-Technology/nidus-sync/db/enums"
	"github.com/Gleipnir-Technology/nidus-sync/db/models"
	"github.com/Gleipnir-Technology/nidus-sync/platform/event"
	"github.com/Gleipnir-Technology/nidus-sync/platform/publicreport"
	"github.com/Gleipnir-Technology/nidus-sync/platform/types"
	//"github.com/Gleipnir-Technology/nidus-sync/platform/geocode"
	//"github.com/Gleipnir-Technology/nidus-sync/platform/geom"
	"github.com/aarondl/opt/omit"
	"github.com/aarondl/opt/omitnull"
	"github.com/rs/zerolog/log"
	"github.com/stephenafamo/scan"
)

type Signal struct {
	Address   *types.Address `db:"address" json:"address"`
	Addressed *time.Time     `db:"addressed" json:"addressed"`
	Addressor *int32         `db:"addressor" json:"addressor"`
	Created   time.Time      `db:"created" json:"created"`
	Creator   int32          `db:"creator" json:"creator"`
	ID        int32          `db:"id" json:"id"`
	Location  types.Location `db:"location" json:"location"`
	Pool      *Pool          `db:"pool" json:"pool"`
	Report    *types.Report  `db:"report" json:"report"`
	Species   *string        `db:"species" json:"species"`
	Type      string         `db:"type" json:"type"`
}

// Create a lead from the given signal and site
func SignalCreateFromPublicreport(ctx context.Context, user User, report_id string) (*int32, error) {
	txn, err := db.PGInstance.BobDB.BeginTx(ctx, nil)
	defer txn.Rollback(ctx)
	if err != nil {
		return nil, fmt.Errorf("start transaction: %w", err)
	}

	report, err := models.PublicreportReports.Query(
		models.SelectWhere.PublicreportReports.PublicID.EQ(report_id),
		models.SelectWhere.PublicreportReports.OrganizationID.EQ(user.Organization.ID),
	).One(ctx, txn)
	if err != nil {
		return nil, fmt.Errorf("query report existence: %w", err)
	}

	// At this point we have a report. We need to decide where to put it based on either the address or
	// the location.
	var site_id int32
	var location string
	if report.AddressID.IsValue() {
		address_id := report.AddressID.MustGet()
		address, err := models.FindAddress(ctx, txn, address_id)
		if err != nil {
			return nil, fmt.Errorf("find address: %w", err)
		}
		site, err := siteFromAddress(ctx, txn, user, address_id)
		if err != nil {
			return nil, fmt.Errorf("site from address: %w", err)
		}
		site_id = site.ID
		lat := address.LocationY.GetOr(0.0)
		lng := address.LocationX.GetOr(0.0)
		location = fmt.Sprintf("POINT(%f %f)", lng, lat)
	} else if report.LocationLatitude.IsValue() && report.LocationLongitude.IsValue() {
		lat := report.LocationLatitude.MustGet()
		lng := report.LocationLongitude.MustGet()
		site, err := siteFromLocation(ctx, txn, user, Location{
			Latitude:  lat,
			Longitude: lng,
		})
		if err != nil {
			return nil, fmt.Errorf("site from address: %w", err)
		}
		site_id = site.ID
		location = fmt.Sprintf("POINT(%f %f)", lng, lat)
	} else if report.AddressRaw != "" {
		// At this point we don't have an address, and we don't have GPS
		// We'll try geocoding and creating an address from that.
		site, err := siteFromAddressRaw(ctx, txn, user, report.AddressRaw)
		if err != nil {
			return nil, fmt.Errorf("site from address: %w", err)
		}
		address, err := models.FindAddress(ctx, txn, site.AddressID)
		if err != nil {
			return nil, fmt.Errorf("find address from raw: %w", err)
		}
		site_id = site.ID
		lat := address.LocationY.GetOr(0.0)
		lng := address.LocationX.GetOr(0.0)
		location = fmt.Sprintf("POINT(%f %f)", lng, lat)
	} else {
		// We have no structured address, no GPS, no unstructued address.
		// There's really nothing we can make this lead from and have it be meaningful
		return nil, errors.New("Refusing to create a signal with no location data.")
	}

	var signal_type enums.Signaltype
	switch report.ReportType {
	case enums.PublicreportReporttypeNuisance:
		signal_type = enums.SignaltypePublicreportNuisance
	case enums.PublicreportReporttypeWater:
		signal_type = enums.SignaltypePublicreportWater
	default:
		return nil, fmt.Errorf("Unrecognized report type %s", string(report.ReportType))
	}
	log.Debug().Str("location", location).Msg("inserting signal")
	signal, err := models.Signals.Insert(&models.SignalSetter{
		Addressed:            omitnull.FromPtr[time.Time](nil),
		Addressor:            omitnull.FromPtr[int32](nil),
		Created:              omit.From(time.Now()),
		Creator:              omit.From(int32(user.ID)),
		FeaturePoolFeatureID: omitnull.FromPtr[int32](nil),
		// ID
		OrganizationID: omit.From(int32(user.Organization.ID)),
		Location:       omit.From(location),
		ReportID:       omitnull.From(report.ID),
		Species:        omitnull.FromPtr[enums.Mosquitospecies](nil),
		SiteID:         omitnull.From(site_id),
		Type:           omit.From[enums.Signaltype](signal_type),
	}).One(ctx, txn)
	if err != nil {
		return nil, fmt.Errorf("create signal: %w", err)
	}
	_, err = psql.Update(
		um.Table(psql.Quote("publicreport", "report")),
		um.SetCol("reviewed").ToArg(time.Now()),
		um.SetCol("reviewer_id").ToArg(user.ID),
		um.SetCol("status").ToArg(enums.PublicreportReportstatustypeReviewed),
		um.Where(psql.Quote("public_id").EQ(psql.Arg(report_id))),
	).Exec(ctx, txn)
	if err != nil {
		return nil, fmt.Errorf("failed to update report %d: %w", report_id, err)
	}
	event.Created(event.TypeSignal, user.Organization.ID, strconv.Itoa(int(signal.ID)))
	txn.Commit(ctx)

	return &signal.ID, nil
}

func SignalList(ctx context.Context, user User, limit int) ([]*Signal, error) {
	org_id := user.Organization.ID
	rows, err := bob.All(ctx, db.PGInstance.BobDB, psql.Select(
		sm.Columns(
			"signal.addressed AS addressed",
			"signal.addressor AS addressor",
			"signal.created AS created",
			"signal.creator AS creator",
			"signal.id AS id",
			"COALESCE(signal.feature_pool_feature_id, 0) AS \"pool.id\"",
			"COALESCE(signal.report_id, 0) AS \"report.id\"",
			"signal.species AS species",
			"signal.type_ AS type",
			"COALESCE(address.country, '') AS \"address.country\"",
			"COALESCE(address.locality, '') AS \"address.locality\"",
			"COALESCE(address.number_, '') AS \"address.number\"",
			"COALESCE(address.postal_code, '') AS \"address.postal_code\"",
			"COALESCE(address.region, '') AS \"address.region\"",
			"COALESCE(address.street, '') AS \"address.street\"",
			"COALESCE(address.unit, '') AS \"address.unit\"",
			// This will work great, up until we add polygons to signal
			"ST_Y(signal.location) AS \"location.latitude\"",
			"ST_X(signal.location) AS \"location.longitude\"",
		),
		sm.From("signal"),
		sm.LeftJoin("site").OnEQ(
			psql.Quote("signal", "site_id"),
			psql.Quote("site", "id"),
		),
		sm.LeftJoin("address").OnEQ(
			psql.Quote("site", "address_id"),
			psql.Quote("address", "id"),
		),
		sm.Where(psql.Quote("signal", "organization_id").EQ(psql.Arg(org_id))),
		sm.Where(psql.Quote("signal", "addressed").IsNull()),
		sm.Limit(limit),
	), scan.StructMapper[*Signal]())
	log.Debug().Int("len", len(rows)).Msg("got signals")
	if err != nil {
		return nil, fmt.Errorf("failed to get signals: %w", err)
	}
	report_ids := make([]int32, 0)
	pool_ids := make([]int32, 0)
	for _, row := range rows {
		if row.Report.ID != 0 {
			report_ids = append(report_ids, row.Report.ID)
		} else if row.Pool.ID != 0 {
			pool_ids = append(pool_ids, row.Pool.ID)
		}
	}
	pools, err := poolList(ctx, org_id, pool_ids)
	if err != nil {
		return nil, fmt.Errorf("getting pools by ID: %w", err)
	}
	reports, err := publicreport.Reports(ctx, org_id, report_ids)
	if err != nil {
		return nil, fmt.Errorf("getting reports by ID: %w", err)
	}
	pool_map := make(map[int32]*Pool, len(pools))
	for _, pool := range pools {
		pool_map[pool.ID] = pool
		log.Debug().Int32("pool", pool.ID).Msg("Added to map")
	}
	report_map := make(map[int32]*types.Report, len(report_ids))
	for _, report := range reports {
		report_map[report.ID] = report
	}
	for _, row := range rows {
		if row.Pool.ID != 0 {
			p, ok := pool_map[row.Pool.ID]
			if !ok {
				return nil, fmt.Errorf("failed to get pool %d for %d", row.Pool.ID, row.ID)
			}
			if p == nil {
				return nil, fmt.Errorf("got nil pool from %d for %d", row.Pool.ID, row.ID)
			}
			row.Pool = p
			row.Report = nil
		} else if row.Report.ID != 0 {
			row.Pool = nil
			row.Report = report_map[row.Report.ID]
		}
		if row.Address.Street == "" {
			row.Address = nil
		}
	}
	return rows, nil
}
