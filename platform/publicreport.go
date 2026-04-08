package platform

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Gleipnir-Technology/bob"
	"github.com/Gleipnir-Technology/bob/dialect/psql"
	"github.com/Gleipnir-Technology/bob/dialect/psql/um"
	"github.com/Gleipnir-Technology/nidus-sync/db"
	"github.com/Gleipnir-Technology/nidus-sync/db/enums"
	"github.com/Gleipnir-Technology/nidus-sync/db/models"
	"github.com/aarondl/opt/omit"
	"github.com/aarondl/opt/omitnull"
	//"github.com/Gleipnir-Technology/nidus-sync/platform/background"
	"github.com/Gleipnir-Technology/nidus-sync/platform/email"
	"github.com/Gleipnir-Technology/nidus-sync/platform/event"
	"github.com/Gleipnir-Technology/nidus-sync/platform/geocode"
	"github.com/Gleipnir-Technology/nidus-sync/platform/report"
	"github.com/Gleipnir-Technology/nidus-sync/platform/text"
	"github.com/Gleipnir-Technology/nidus-sync/platform/types"
	"github.com/rs/zerolog/log"
)

func PublicreportByID(ctx context.Context, report_id string) (*models.PublicreportReport, error) {
	report, err := models.PublicreportReports.Query(
		models.SelectWhere.PublicreportReports.PublicID.EQ(report_id),
	).One(ctx, db.PGInstance.BobDB)
	if err != nil {
		return nil, err
	}
	return report, nil
}
func PublicreportInvalid(ctx context.Context, user User, report_id string) error {
	report, err := reportFromID(ctx, user, report_id)
	if err != nil {
		return fmt.Errorf("query report existence: %w", err)
	}

	err = report.Update(ctx, db.PGInstance.BobDB, &models.PublicreportReportSetter{
		Reviewed:   omitnull.From(time.Now()),
		ReviewerID: omitnull.From(int32(user.ID)),
		Status:     omit.From(enums.PublicreportReportstatustypeInvalidated),
	})

	log.Info().Int32("id", report.ID).Msg("Report marked as invalid")
	event.Updated(event.TypeRMOReport, user.Organization.ID, report_id)
	return nil
}

func PublicReportMessageCreate(ctx context.Context, user User, report_id, message string) (message_id *int32, err error) {
	txn, err := db.PGInstance.BobDB.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("create txn: %w", err)
	}
	defer txn.Rollback(ctx)

	report, err := reportFromID(ctx, user, report_id)
	if err != nil {
		return nil, fmt.Errorf("query report existence: %w", err)
	}
	if report.ReporterPhone != "" {
		log.Debug().Str("report_id", report_id).Msg("contacting via phone")
		p, err := text.ParsePhoneNumber(report.ReporterPhone)
		if err != nil {
			return nil, fmt.Errorf("parse phone: %w", err)
		}
		msg_id, err := text.ReportMessage(ctx, txn, int32(user.ID), int32(report.ID), *p, message)
		if err != nil {
			return nil, fmt.Errorf("send text: %w", err)
		}
		txn.Commit(ctx)
		//log.Debug().Int32("msg_id", *msg_id).Msg("Created text.ReportMessage")
		return msg_id, nil
	} else if report.ReporterEmail != "" {
		msg_id, err := email.ReportMessage(ctx, int32(user.ID), report_id, report.ReporterEmail, message)
		if err != nil {
			return nil, fmt.Errorf("send email: %w", err)
		}
		txn.Commit(ctx)
		return msg_id, nil
	} else {
		log.Debug().Str("report_id", report_id).Msg("contacting via email")
		return nil, errors.New("no contact methods available")
	}
}
func PublicReportReporterUpdated(ctx context.Context, org_id int32, report_id string) {
	event.Updated(event.TypeRMOReport, org_id, report_id)
}
func ReportNuisanceCreate(ctx context.Context, setter_report models.PublicreportReportSetter, setter_nuisance models.PublicreportNuisanceSetter, latlng LatLng, address Address, images []ImageUpload) (*models.PublicreportReport, error) {
	return reportCreate(ctx, setter_report, latlng, address, images, func(ctx context.Context, txn bob.Executor, report_id int32) error {
		setter_nuisance.ReportID = omit.From(report_id)
		_, err := models.PublicreportNuisances.Insert(&setter_nuisance).One(ctx, txn)
		if err != nil {
			return fmt.Errorf("Failed to create nuisance database record: %w", err)
		}
		return nil
	})
}

func ReportWaterCreate(ctx context.Context, setter_report models.PublicreportReportSetter, setter_water models.PublicreportWaterSetter, latlng LatLng, address Address, images []ImageUpload) (*models.PublicreportReport, error) {
	return reportCreate(ctx, setter_report, latlng, address, images, func(ctx context.Context, txn bob.Executor, report_id int32) error {
		setter_water.ReportID = omit.From(report_id)
		_, err := models.PublicreportWaters.Insert(&setter_water).One(ctx, txn)
		if err != nil {
			return fmt.Errorf("Failed to create water database record: %w", err)
		}
		return nil
	})
}

type funcSetReportDetail = func(context.Context, bob.Executor, int32) error

func reportCreate(ctx context.Context, setter_report models.PublicreportReportSetter, latlng LatLng, address Address, images []ImageUpload, detail_setter funcSetReportDetail) (result *models.PublicreportReport, err error) {
	txn, err := db.PGInstance.BobDB.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("create txn: %w", err)
	}
	defer txn.Rollback(ctx)

	public_id, err := report.GenerateReportID()
	if err != nil {
		return nil, fmt.Errorf("create public ID: %w", err)
	}
	setter_report.PublicID = omit.From(public_id)

	// If we've got an locality value it was set by geocoding so we should save it
	var a *models.Address
	if address.Locality != "" && latlng.Latitude != nil && latlng.Longitude != nil {
		a, err = geocode.EnsureAddress(ctx, txn, address, types.Location{
			Latitude:  *latlng.Latitude,
			Longitude: *latlng.Longitude,
		})
		if err != nil {
			return nil, fmt.Errorf("Failed to ensure address: %w", err)
		}
	}

	saved_images, err := saveImageUploads(ctx, txn, images)
	if err != nil {
		return nil, fmt.Errorf("Failed to save image uploads: %w", err)
	}
	var organization_id *int32
	organization_id, err = MatchDistrict(ctx, latlng.Longitude, latlng.Latitude, images)
	if err != nil {
		log.Warn().Err(err).Msg("Failed to match district")
	}

	if a != nil {
		setter_report.AddressID = omitnull.From(a.ID)
	}
	if organization_id != nil {
		setter_report.OrganizationID = omit.FromPtr(organization_id)
	}
	result, err = models.PublicreportReports.Insert(&setter_report).One(ctx, txn)
	if err != nil {
		return nil, fmt.Errorf("Failed to create report database record: %w", err)
	}
	if latlng.Latitude != nil && latlng.Longitude != nil {
		h3cell, _ := latlng.H3Cell()
		geom_query, _ := latlng.GeometryQuery()
		_, err = psql.Update(
			um.Table("publicreport.report"),
			um.SetCol("h3cell").ToArg(h3cell),
			um.SetCol("location").To(geom_query),
			um.Where(psql.Quote("id").EQ(psql.Arg(result.ID))),
		).Exec(ctx, txn)
		if err != nil {
			return nil, fmt.Errorf("Failed to insert publicreport.report geospatial", err)
		}
	}
	log.Info().Str("public_id", public_id).Int32("id", result.ID).Msg("Created base report")

	if len(saved_images) > 0 {
		setters := make([]*models.PublicreportReportImageSetter, 0)
		for _, image := range saved_images {
			setters = append(setters, &models.PublicreportReportImageSetter{
				ImageID:  omit.From(int32(image.ID)),
				ReportID: omit.From(int32(result.ID)),
			})
		}
		_, err = models.PublicreportReportImages.Insert(bob.ToMods(setters...)).Exec(ctx, txn)
		if err != nil {
			return nil, fmt.Errorf("Failed to save reference to images: %w", err)
		}
		log.Info().Int("len", len(images)).Msg("saved uploaded images")
	}

	err = detail_setter(ctx, txn, result.ID)
	if err != nil {
		return nil, fmt.Errorf("detail setter: %w", err)
	}

	models.PublicreportReportLogs.Insert(&models.PublicreportReportLogSetter{
		Created:    omit.From(time.Now()),
		EmailLogID: omitnull.FromPtr[int32](nil),
		// ID
		ReportID:  omit.From(result.ID),
		TextLogID: omitnull.FromPtr[int32](nil),
		Type:      omit.From(enums.PublicreportReportlogtypeCreated),
		UserID:    omitnull.FromPtr[int32](nil),
	}).One(ctx, txn)

	txn.Commit(ctx)

	if organization_id != nil {
		event.Created(
			event.TypeRMONuisance,
			*organization_id,
			result.PublicID,
		)
	}
	return result, nil
}
func reportFromID(ctx context.Context, user User, report_id string) (*models.PublicreportReport, error) {
	report, err := models.PublicreportReports.Query(
		models.SelectWhere.PublicreportReports.PublicID.EQ(report_id),
		models.SelectWhere.PublicreportReports.OrganizationID.EQ(user.Organization.ID),
	).One(ctx, db.PGInstance.BobDB)
	if err != nil {
		return nil, err
	}
	return report, nil
}
