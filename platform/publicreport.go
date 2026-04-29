package platform

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"math/big"
	"strings"
	"time"

	"github.com/Gleipnir-Technology/bob"
	"github.com/Gleipnir-Technology/bob/dialect/psql"
	"github.com/Gleipnir-Technology/bob/dialect/psql/sm"
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
	"github.com/Gleipnir-Technology/nidus-sync/platform/publicreport"
	"github.com/Gleipnir-Technology/nidus-sync/platform/text"
	"github.com/Gleipnir-Technology/nidus-sync/platform/types"
	"github.com/rs/zerolog/log"
)

// GenerateReportID creates a 12-character random string using only unambiguous
// capital letters and numbers
func GenerateReportID() (string, error) {
	// Define character set (no O/0, I/l/1, 2/Z to avoid confusion)
	const charset = "ABCDEFGHJKLMNPQRSTUVWXY3456789"
	const length = 12

	var builder strings.Builder
	builder.Grow(length)

	// Use crypto/rand for secure randomness
	for i := 0; i < length; i++ {
		// Generate a random index within our charset
		n, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			return "", fmt.Errorf("failed to generate random number: %w", err)
		}

		// Add the randomly selected character to our ID
		builder.WriteByte(charset[n.Int64()])
	}

	return builder.String(), nil
}

func PublicReportByID(ctx context.Context, report_id string, is_public bool) (*types.PublicReport, error) {
	return publicreport.ByID(ctx, report_id, is_public)
}
func PublicReportByIDCompliance(ctx context.Context, report_id string, is_public bool) (*types.PublicReportCompliance, error) {
	result, err := publicreport.ByIDCompliance(ctx, report_id, is_public)
	if err != nil {
		return nil, fmt.Errorf("byidcompliance: %w", err)
	}
	// Check for evidence if this is a mailer-based compliance request
	crr, err := ComplianceReportRequestFromPublicID(ctx, result.PublicID)
	if err != nil {
		return nil, fmt.Errorf("compliance report request by public id: %w", err)
	}
	if crr != nil {
		result.Concerns = []*types.ConcernComplianceReportRequest{
			&types.ConcernComplianceReportRequest{
				ComplianceReportRequestPublicID: crr.PublicID,
			},
		}
	}
	return result, nil
}
func PublicReportByIDNuisance(ctx context.Context, report_id string, is_public bool) (*types.PublicReportNuisance, error) {
	return publicreport.ByIDNuisance(ctx, report_id, is_public)
}
func PublicReportByIDWater(ctx context.Context, report_id string, is_public bool) (*types.PublicReportWater, error) {
	return publicreport.ByIDWater(ctx, report_id, is_public)
}
func PublicReportComplianceSubmit(ctx context.Context, report_id string, is_public bool) (*types.PublicReportCompliance, error) {
	report, err := publicreport.ByIDCompliance(ctx, report_id, is_public)
	if err != nil {
		return nil, fmt.Errorf("byidcompliance: %w", err)
	}
	_, err = psql.Update(
		um.Table(models.PublicreportCompliances.NameAs()),
		um.SetCol(models.PublicreportCompliances.Columns.Submitted.String()).ToArg(time.Now()),
		um.Where(models.PublicreportCompliances.Columns.ReportID.EQ(psql.Arg(report.ReportID))),
	).Exec(ctx, db.PGInstance.BobDB)
	if err != nil {
		return nil, fmt.Errorf("update report submitted: %w", err)
	}
	return publicreport.ByIDCompliance(ctx, report_id, is_public)
}
func PublicReportInvalid(ctx context.Context, user User, public_id string) error {
	report, err := publicReportFromID(ctx, public_id)
	if err != nil {
		return fmt.Errorf("query report existence: %w", err)
	}
	if report.OrganizationID != user.Organization.ID {
		return fmt.Errorf("user is from a different organization")
	}

	err = report.Update(ctx, db.PGInstance.BobDB, &models.PublicreportReportSetter{
		Reviewed:   omitnull.From(time.Now()),
		ReviewerID: omitnull.From(int32(user.ID)),
		Status:     omit.From(enums.PublicreportReportstatustypeInvalidated),
	})

	log.Info().Int32("id", report.ID).Msg("Report marked as invalid")
	event.Updated(event.TypeRMOPublicReport, user.Organization.ID, public_id)
	return nil
}

func PublicReportMessageCreate(ctx context.Context, user User, public_id, message string) (message_id *int32, err error) {
	txn, err := db.PGInstance.BobDB.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("create txn: %w", err)
	}
	defer txn.Rollback(ctx)

	report, err := publicReportFromID(ctx, public_id)
	if err != nil {
		return nil, fmt.Errorf("query report existence: %w", err)
	}
	if report.OrganizationID != user.Organization.ID {
		return nil, fmt.Errorf("user is from a different organization")
	}
	if report.ReporterPhone != "" {
		log.Debug().Str("public_id", public_id).Msg("contacting via phone")
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
		msg_id, err := email.ReportMessage(ctx, int32(user.ID), public_id, report.ReporterEmail, message)
		if err != nil {
			return nil, fmt.Errorf("send email: %w", err)
		}
		txn.Commit(ctx)
		return msg_id, nil
	} else {
		log.Debug().Str("public_id", public_id).Msg("contacting via email")
		return nil, errors.New("no contact methods available")
	}
}
func PublicReportUpdateCompliance(ctx context.Context, public_id string, report_setter *models.PublicreportReportSetter, compliance_setter *models.PublicreportComplianceSetter, address *types.Address, location *types.Location) (*types.PublicReportCompliance, error) {
	txn, err := db.PGInstance.BobDB.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("create txn: %w", err)
	}
	defer txn.Rollback(ctx)
	report, err := publicReportFromID(ctx, public_id)
	if err != nil {
		return nil, fmt.Errorf("query report existence: %w", err)
	}
	compliance, err := models.FindPublicreportCompliance(ctx, txn, report.ID)
	if err != nil {
		return nil, fmt.Errorf("find compliance %d: %w", report.ID, err)
	}
	// Avoid attempting to perform an empty update
	if report_setter.LatlngAccuracyValue.IsValue() ||
		report_setter.ReporterEmail.IsValue() ||
		report_setter.ReporterName.IsValue() ||
		report_setter.ReporterPhone.IsValue() {
		err = report.Update(ctx, txn, report_setter)
		if err != nil {
			return nil, fmt.Errorf("update report: %w", err)
		}
	}
	// Avoid attempting to perform an empty update
	if compliance_setter.AccessInstructions.IsValue() ||
		compliance_setter.AvailabilityNotes.IsValue() ||
		compliance_setter.Comments.IsValue() ||
		compliance_setter.GateCode.IsValue() ||
		compliance_setter.HasDog.IsValue() ||
		compliance_setter.PermissionType.IsValue() ||
		compliance_setter.ReportPhoneCanText.IsValue() ||
		compliance_setter.WantsScheduled.IsValue() {
		err = compliance.Update(ctx, txn, compliance_setter)
		if err != nil {
			return nil, fmt.Errorf("update compliance: %w", err)
		}
	}
	if address != nil {
		err = publicReportUpdateAddress(ctx, txn, report, *address)
		if err != nil {
			return nil, fmt.Errorf("update address: %w", err)
		}
	}
	if location != nil {
		err = publicReportUpdateLocation(ctx, txn, report.ID, *location)
		if err != nil {
			return nil, fmt.Errorf("update location: %w", err)
		}
	}
	txn.Commit(ctx)
	return publicreport.ByIDCompliance(ctx, public_id, false)
}
func PublicReportReporterUpdated(ctx context.Context, org_id int32, report_id string) {
	event.Updated(event.TypeRMOPublicReport, org_id, report_id)
}
func PublicReportsForOrganization(ctx context.Context, org_id int32, is_public bool) ([]*types.PublicReport, error) {
	return publicreport.ReportsForOrganization(ctx, org_id, is_public)
}
func PublicReportComplianceCreate(ctx context.Context, setter_report models.PublicreportReportSetter, setter_compliance models.PublicreportComplianceSetter, org_id int32) (*models.PublicreportReport, error) {
	return publicReportCreate(ctx, setter_report, nil, nil, nil, org_id, func(ctx context.Context, txn bob.Executor, report_id int32) error {
		setter_compliance.ReportID = omit.From(report_id)
		_, err := models.PublicreportCompliances.Insert(&setter_compliance).One(ctx, txn)
		if err != nil {
			return fmt.Errorf("Failed to create nuisance database record: %w", err)
		}
		return nil
	})
}
func PublicReportImageCreate(ctx context.Context, public_id string, images []ImageUpload) error {
	txn, err := db.PGInstance.BobDB.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("create txn: %w", err)
	}
	defer txn.Rollback(ctx)

	report, err := publicReportFromID(ctx, public_id)
	if err != nil {
		return fmt.Errorf("report from ID: %w", err)
	}
	saved_images, err := saveImageUploads(ctx, txn, images)
	if err != nil {
		return fmt.Errorf("Failed to save image uploads: %w", err)
	}
	if len(saved_images) > 0 {
		setters := make([]*models.PublicreportReportImageSetter, 0)
		for _, image := range saved_images {
			setters = append(setters, &models.PublicreportReportImageSetter{
				ImageID:  omit.From(int32(image.ID)),
				ReportID: omit.From(int32(report.ID)),
			})
		}
		_, err = models.PublicreportReportImages.Insert(bob.ToMods(setters...)).Exec(ctx, txn)
		if err != nil {
			return fmt.Errorf("Failed to save reference to images: %w", err)
		}
		log.Info().Int("len", len(images)).Msg("saved uploaded images")
	}
	txn.Commit(ctx)
	return nil
}
func PublicReportNuisanceCreate(ctx context.Context, setter_report models.PublicreportReportSetter, setter_nuisance models.PublicreportNuisanceSetter, location types.Location, address Address, images []ImageUpload) (*models.PublicreportReport, error) {
	return publicReportCreate(ctx, setter_report, &location, &address, images, 0, func(ctx context.Context, txn bob.Executor, report_id int32) error {
		setter_nuisance.ReportID = omit.From(report_id)
		_, err := models.PublicreportNuisances.Insert(&setter_nuisance).One(ctx, txn)
		if err != nil {
			return fmt.Errorf("Failed to create nuisance database record: %w", err)
		}
		return nil
	})
}

func PublicReportWaterCreate(ctx context.Context, setter_report models.PublicreportReportSetter, setter_water models.PublicreportWaterSetter, location types.Location, address Address, images []ImageUpload) (*models.PublicreportReport, error) {
	return publicReportCreate(ctx, setter_report, &location, &address, images, 0, func(ctx context.Context, txn bob.Executor, report_id int32) error {
		setter_water.ReportID = omit.From(report_id)
		_, err := models.PublicreportWaters.Insert(&setter_water).One(ctx, txn)
		if err != nil {
			return fmt.Errorf("Failed to create water database record: %w", err)
		}
		return nil
	})
}
func PublicReportTypeByID(ctx context.Context, public_id string) (string, error) {
	report, err := models.PublicreportReports.Query(
		models.SelectWhere.PublicreportReports.PublicID.EQ(public_id),
	).One(ctx, db.PGInstance.BobDB)
	if err != nil {
		return "", fmt.Errorf("query report '%s': %w", public_id, err)
	}
	return report.ReportType.String(), nil
}

type funcSetReportDetail = func(context.Context, bob.Executor, int32) error

func publicReportCreate(ctx context.Context, setter_report models.PublicreportReportSetter, location *types.Location, address *Address, images []ImageUpload, organization_id int32, detail_setter funcSetReportDetail) (result *models.PublicreportReport, err error) {
	txn, err := db.PGInstance.BobDB.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("create txn: %w", err)
	}
	defer txn.Rollback(ctx)

	if setter_report.PublicID.IsUnset() {
		public_id, err := GenerateReportID()
		if err != nil {
			return nil, fmt.Errorf("create public ID: %w", err)
		}
		setter_report.PublicID = omit.From(public_id)
	}

	var addr *models.Address
	if address != nil {
		if address.GID != "" {
			addr, err = geocode.EnsureAddress(ctx, txn, *address)
			if err != nil {
				return nil, fmt.Errorf("Failed to ensure address: %w", err)
			}
		} else if address.Raw != "" {
			geo_res, err := geocode.GeocodeRaw(ctx, nil, address.Raw)
			if err != nil {
				return nil, fmt.Errorf("Failed to geocode raw: %w", err)
			}
			addr, err = models.FindAddress(ctx, txn, *geo_res.Address.ID)
			if err != nil {
				return nil, fmt.Errorf("Failed to lookup address: %w", err)
			}
		} else {
			return nil, fmt.Errorf("empty address")
		}
	}

	saved_images, err := saveImageUploads(ctx, txn, images)
	if err != nil {
		return nil, fmt.Errorf("Failed to save image uploads: %w", err)
	}
	if organization_id == 0 {
		organization_id, err = matchDistrict(ctx, location, images, addr)
		if err != nil {
			log.Warn().Err(err).Msg("Failed to match district")
		}
	}
	setter_report.OrganizationID = omit.From(organization_id)

	if addr != nil {
		setter_report.AddressID = omitnull.From(addr.ID)
	}
	result, err = models.PublicreportReports.Insert(&setter_report).One(ctx, txn)
	if err != nil {
		return nil, fmt.Errorf("Failed to create report database record: %w", err)
	}
	if location != nil {
		l := *location
		if l.Latitude != 0 && l.Longitude != 0 {
			publicReportUpdateLocation(ctx, txn, result.ID, l)
		}
	}
	log.Info().Str("public_id", setter_report.PublicID.GetOr("")).Int32("id", result.ID).Msg("Created base report")

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

	event.Created(
		event.TypeRMOPublicReport,
		organization_id,
		result.PublicID,
	)
	return result, nil
}
func publicReportFromID(ctx context.Context, public_id string) (*models.PublicreportReport, error) {
	report, err := models.PublicreportReports.Query(
		models.SelectWhere.PublicreportReports.PublicID.EQ(public_id),
	).One(ctx, db.PGInstance.BobDB)
	if err != nil {
		return nil, err
	}
	return report, nil
}
func publicReportUpdateAddress(ctx context.Context, txn bob.Executor, report *models.PublicreportReport, address types.Address) error {
	err := report.Update(ctx, txn, &models.PublicreportReportSetter{
		AddressGid: omit.From(address.GID),
		AddressRaw: omit.From(address.Raw),
	})
	if err != nil {
		return fmt.Errorf("update report: %w", err)
	}
	_, err = psql.Update(
		um.Table("publicreport.report"),
		um.SetCol("address_id").To(
			psql.Select(
				sm.Columns("id"),
				sm.From("address"),
				sm.Where(psql.Quote("gid").EQ(psql.Arg(address.GID))),
				sm.Limit(1),
			),
		),
		um.Where(psql.Quote("publicreport", "report", "id").EQ(psql.Arg(report.ID))),
	).Exec(ctx, txn)
	if err != nil {
		return fmt.Errorf("update report address_id: %w", err)
	}
	return nil
}
func publicReportUpdateLocation(ctx context.Context, txn bob.Executor, id int32, location types.Location) error {
	h3cell, _ := location.H3Cell()
	geom_query, _ := location.GeometryQuery()
	_, err := psql.Update(
		um.Table("publicreport.report"),
		um.SetCol("h3cell").ToArg(h3cell),
		um.SetCol("location").To(geom_query),
		um.Where(psql.Quote("id").EQ(psql.Arg(id))),
	).Exec(ctx, txn)
	if err != nil {
		return fmt.Errorf("Failed to insert publicreport.report geospatial", err)
	}
	return nil
}
