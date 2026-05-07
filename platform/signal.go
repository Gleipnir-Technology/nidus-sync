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
	"github.com/Gleipnir-Technology/nidus-sync/db"
	modelpublic "github.com/Gleipnir-Technology/nidus-sync/db/gen/nidus-sync/public/model"
	modelpublicreport "github.com/Gleipnir-Technology/nidus-sync/db/gen/nidus-sync/publicreport/model"
	tablepublicreport "github.com/Gleipnir-Technology/nidus-sync/db/gen/nidus-sync/publicreport/table"
	querypublic "github.com/Gleipnir-Technology/nidus-sync/db/query/public"
	querypublicreport "github.com/Gleipnir-Technology/nidus-sync/db/query/publicreport"
	"github.com/Gleipnir-Technology/nidus-sync/platform/event"
	"github.com/Gleipnir-Technology/nidus-sync/platform/publicreport"
	"github.com/Gleipnir-Technology/nidus-sync/platform/types"
	"github.com/rs/zerolog/log"
	"github.com/stephenafamo/scan"
	"github.com/twpayne/go-geom"
)

type Signal struct {
	Address   *types.Address      `db:"address" json:"address"`
	Addressed *time.Time          `db:"addressed" json:"addressed"`
	Addressor *int32              `db:"addressor" json:"addressor"`
	Created   time.Time           `db:"created" json:"created"`
	Creator   int32               `db:"creator" json:"creator"`
	ID        int32               `db:"id" json:"id"`
	Location  types.Location      `db:"location" json:"location"`
	Pool      *Pool               `db:"pool" json:"pool"`
	Report    *types.PublicReport `db:"report" json:"report"`
	Species   *string             `db:"species" json:"species"`
	Type      string              `db:"type" json:"type"`
}

type _rowWithID struct {
	ID int32 `db:"id"`
}

func SignalCreateFromPool(ctx context.Context, txn db.Ex, user User, site_id int32, feature_id int32, location types.Location) (modelpublic.Signal, error) {
	g := location.ToGeom()
	signal := modelpublic.Signal{
		Addressed:            nil,
		Addressor:            nil,
		Created:              time.Now(),
		Creator:              int32(user.ID),
		FeaturePoolFeatureID: &feature_id,
		//ID
		Location:       g,
		OrganizationID: user.Organization.ID,
		ReportID:       nil,
		SiteID:         &site_id,
		Species:        nil,
		Type:           modelpublic.Signaltype_FlyoverPool,
	}
	var err error
	signal, err = querypublic.SignalInsert(ctx, txn, signal)
	if err != nil {
		return modelpublic.Signal{}, fmt.Errorf("insert signal: %w", err)
	}
	return signal, nil
}

// Create a lead from the given signal and site
func SignalCreateFromPublicreport(ctx context.Context, user User, report_id string) (*int32, error) {
	txn, err := db.BeginTxn(ctx)
	defer txn.Rollback(ctx)
	if err != nil {
		return nil, fmt.Errorf("start transaction: %w", err)
	}

	report, err := querypublicreport.ReportFromPublicIDForOrg(ctx, txn, report_id, int64(user.Organization.ID))
	if err != nil {
		return nil, fmt.Errorf("query report existence: %w", err)
	}

	// At this point we have a report. We need to decide where to put it based on either the address or
	// the location.
	var site_id int32
	var location geom.T
	if report.AddressID != nil {
		address_id := *report.AddressID
		address, err := querypublic.AddressFromID(ctx, txn, int64(address_id))
		if err != nil {
			return nil, fmt.Errorf("find address: %w", err)
		}
		site, err := querypublic.SiteFromAddressIDForOrg(ctx, txn, int64(address_id), int64(user.Organization.ID))
		if err != nil {
			return nil, fmt.Errorf("site from address: %w", err)
		}
		site_id = site.ID
		location = address.Location
	} else if report.Location != nil {
		l, err := types.LocationFromGeom(*report.Location)
		if err != nil {
			return nil, fmt.Errorf("report location to geom: %w", err)
		}
		site, err := siteFromLocation(ctx, txn, user, l)
		if err != nil {
			return nil, fmt.Errorf("site from address: %w", err)
		}
		site_id = site.ID
		location = *report.Location
	} else if report.AddressRaw != "" {
		// At this point we don't have an address, and we don't have GPS
		// We'll try geocoding and creating an address from that.
		site, err := siteFromAddressRaw(ctx, txn, user, report.AddressRaw)
		if err != nil {
			return nil, fmt.Errorf("site from address: %w", err)
		}
		address, err := querypublic.AddressFromID(ctx, txn, int64(site.AddressID))
		if err != nil {
			return nil, fmt.Errorf("find address from raw: %w", err)
		}
		site_id = site.ID
		location = address.Location
	} else {
		// We have no structured address, no GPS, no unstructued address.
		// There's really nothing we can make this lead from and have it be meaningful
		return nil, errors.New("Refusing to create a signal with no location data.")
	}

	var signal_type modelpublic.Signaltype
	switch report.ReportType {
	case modelpublicreport.Reporttype_Nuisance:
		signal_type = modelpublic.Signaltype_PublicreportNuisance
	case modelpublicreport.Reporttype_Water:
		signal_type = modelpublic.Signaltype_PublicreportWater
	default:
		return nil, fmt.Errorf("Unrecognized report type %s", string(report.ReportType))
	}
	signal := modelpublic.Signal{
		Addressed:            nil,
		Addressor:            nil,
		Created:              time.Now(),
		Creator:              int32(user.ID),
		FeaturePoolFeatureID: nil,
		// ID
		OrganizationID: int32(user.Organization.ID),
		Location:       location,
		ReportID:       &report.ID,
		Species:        nil,
		SiteID:         &site_id,
		Type:           signal_type,
	}
	signal, err = querypublic.SignalInsert(ctx, txn, signal)
	if err != nil {
		return nil, fmt.Errorf("create signal: %w", err)
	}
	report_updater := querypublicreport.ReportUpdater{}
	now := time.Now()
	report_updater.Model.Reviewed = &now
	report_updater.Set(tablepublicreport.Report.Reviewed)
	user_id := int32(user.ID)
	report_updater.Model.ReviewerID = &user_id
	report_updater.Set(tablepublicreport.Report.ReviewerID)
	report_updater.Model.Status = modelpublicreport.Reportstatustype_Reviewed
	report_updater.Set(tablepublicreport.Report.Status)
	err = report_updater.Execute(ctx, txn, report_id)
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
			"COALESCE(address.country, 'usa') AS \"address.country\"",
			"COALESCE(address.locality, '') AS \"address.locality\"",
			"COALESCE(address.number_, '') AS \"address.number_\"",
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
	reports, err := publicreport.Reports(ctx, org_id, report_ids, false)
	if err != nil {
		return nil, fmt.Errorf("getting reports by ID: %w", err)
	}
	pool_map := make(map[int32]*Pool, len(pools))
	for _, pool := range pools {
		pool_map[pool.ID] = pool
		log.Debug().Int32("pool", pool.ID).Msg("Added to map")
	}
	report_map := make(map[int32]*types.PublicReport, len(report_ids))
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
			report, ok := report_map[row.Report.ID]
			if !ok {
				return nil, fmt.Errorf("failed to get report %d for %d", row.Report.ID, row.ID)
			}
			if report == nil {
				return nil, fmt.Errorf("got nil for report %d for %d", row.Report.ID, row.ID)
			}
			row.Pool = nil
			row.Report = report
		} else {
			log.Debug().Int32("id", row.ID).Msg("has no publicrreport nor pool")
			row.Pool = nil
			row.Report = nil
		}
		if row.Address.Street == "" {
			row.Address = nil
		}
	}
	return rows, nil
}
