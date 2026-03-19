package platform

import (
	"context"
	"errors"
	"fmt"
	"time"

	//"github.com/Gleipnir-Technology/bob"
	"github.com/Gleipnir-Technology/bob/dialect/psql"
	"github.com/Gleipnir-Technology/bob/dialect/psql/um"
	"github.com/Gleipnir-Technology/nidus-sync/db"
	"github.com/Gleipnir-Technology/nidus-sync/db/enums"
	"github.com/Gleipnir-Technology/nidus-sync/db/models"
	//"github.com/Gleipnir-Technology/nidus-sync/platform/geocode"
	//"github.com/Gleipnir-Technology/nidus-sync/platform/geom"
	"github.com/aarondl/opt/omit"
	"github.com/aarondl/opt/omitnull"
	//"github.com/rs/zerolog/log"
)

// Create a lead from the given signal and site
func SignalCreateFromPublicreport(ctx context.Context, user User, report_id string) (*int32, error) {
	txn, err := db.PGInstance.BobDB.BeginTx(ctx, nil)
	defer txn.Rollback(ctx)
	if err != nil {
		return nil, fmt.Errorf("start transaction: %w", err)
	}

	report, err := models.PublicreportReports.Query(
		models.SelectWhere.PublicreportReports.PublicID.EQ(report_id),
		models.SelectWhere.PublicreportReports.OrganizationID.EQ(user.Organization.ID()),
	).One(ctx, txn)
	if err != nil {
		return nil, fmt.Errorf("query report existence: %w", err)
	}

	// At this point we have a report. We need to decide where to put it based on either the address or
	// the location.
	var site_id int32
	if report.AddressID.IsValue() {
		site, err := siteFromAddress(ctx, txn, user, report.AddressID.MustGet())
		if err != nil {
			return nil, fmt.Errorf("site from address: %w", err)
		}
		site_id = site.ID
	} else if report.LocationLatitude.IsValue() && report.LocationLongitude.IsValue() {
		site, err := siteFromLocation(ctx, txn, user, Location{
			Latitude:  report.LocationLatitude.MustGet(),
			Longitude: report.LocationLongitude.MustGet(),
		})
		if err != nil {
			return nil, fmt.Errorf("site from address: %w", err)
		}
		site_id = site.ID

	} else if report.AddressRaw != "" {
		// At this point we don't have an address, and we don't have GPS
		// We'll try geocoding and creating an address from that.
		site, err := siteFromAddressRaw(ctx, txn, user, report.AddressRaw)
		if err != nil {
			return nil, fmt.Errorf("site from address: %w", err)
		}
		site_id = site.ID
	} else {
		// We have no structured address, no GPS, no unstructued address.
		// There's really nothing we can make this lead from and have it be meaningful
		return nil, errors.New("Refusing to create a lead with no location data.")
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
	signal, err := models.Signals.Insert(&models.SignalSetter{
		Addressed: omitnull.FromPtr[time.Time](nil),
		Addressor: omitnull.FromPtr[int32](nil),
		Created:   omit.From(time.Now()),
		Creator:   omit.From(int32(user.ID)),
		// ID
		OrganizationID: omit.From(int32(user.Organization.ID())),
		Species:        omitnull.FromPtr[enums.Mosquitospecies](nil),
		SiteID:         omitnull.From(site_id),
		Title:          omit.From[string](""),
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
	txn.Commit(ctx)

	return &signal.ID, nil
}
