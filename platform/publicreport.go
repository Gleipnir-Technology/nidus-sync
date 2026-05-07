package platform

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"math/big"
	"strings"
	"time"

	"github.com/Gleipnir-Technology/nidus-sync/db"
	"github.com/Gleipnir-Technology/nidus-sync/db/gen/nidus-sync/public/model"
	tablepublic "github.com/Gleipnir-Technology/nidus-sync/db/gen/nidus-sync/public/table"
	modelpublicreport "github.com/Gleipnir-Technology/nidus-sync/db/gen/nidus-sync/publicreport/model"
	tablepublicreport "github.com/Gleipnir-Technology/nidus-sync/db/gen/nidus-sync/publicreport/table"
	querypublic "github.com/Gleipnir-Technology/nidus-sync/db/query/public"
	querypublicreport "github.com/Gleipnir-Technology/nidus-sync/db/query/publicreport"
	"github.com/Gleipnir-Technology/nidus-sync/platform/email"
	"github.com/Gleipnir-Technology/nidus-sync/platform/event"
	"github.com/Gleipnir-Technology/nidus-sync/platform/geocode"
	"github.com/Gleipnir-Technology/nidus-sync/platform/publicreport"
	"github.com/Gleipnir-Technology/nidus-sync/platform/text"
	"github.com/Gleipnir-Technology/nidus-sync/platform/types"
	"github.com/go-jet/jet/v2/postgres"
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

func PublicReportByIDCompliance(ctx context.Context, report_id string, is_public bool) (*types.PublicReportCompliance, error) {
	result, err := publicreport.ByIDCompliance(ctx, report_id, is_public)
	if err != nil {
		return nil, fmt.Errorf("byidcompliance: %w", err)
	}
	if result == nil {
		return nil, nil
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
func PublicReportInvalid(ctx context.Context, user User, public_id string) error {
	report, err := querypublicreport.ReportFromPublicID(ctx, db.PGInstance.PGXPool, public_id)
	if err != nil {
		return fmt.Errorf("query report existence: %w", err)
	}
	if report.OrganizationID != user.Organization.ID {
		return fmt.Errorf("user is from a different organization")
	}

	now := time.Now()
	report_updater := querypublicreport.ReportUpdater{}
	report_updater.Model.Reviewed = &now
	report_updater.Set(tablepublicreport.Report.Reviewed)
	reporter_id := int32(user.ID)
	report_updater.Model.ReviewerID = &reporter_id
	report_updater.Set(tablepublicreport.Report.ReviewerID)
	report_updater.Model.Status = modelpublicreport.Reportstatustype_Invalidated
	report_updater.Set(tablepublicreport.Report.Status)
	err = report_updater.Execute(ctx, db.PGInstance.PGXPool, report.ID)

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

	report, err := querypublicreport.ReportFromPublicID(ctx, db.PGInstance.PGXPool, public_id)
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
func PublicReportUpdateCompliance(ctx context.Context, public_id string, report_updates querypublicreport.ReportUpdater, compliance_updates querypublicreport.ComplianceUpdater, address *types.Address, location *types.Location) error {
	//txn, err := db.PGInstance.BobDB.BeginTx(ctx, nil)
	txn, err := db.BeginTxn(ctx)
	if err != nil {
		return fmt.Errorf("create txn: %w", err)
	}
	defer txn.Rollback(ctx)
	report, err := querypublicreport.ReportFromPublicID(ctx, db.PGInstance.PGXPool, public_id)
	if err != nil {
		return fmt.Errorf("query report existence: %w", err)
	}
	//compliance, err := models.FindPublicreportCompliance(ctx, txn, report.ID)
	compliance, err := querypublicreport.ComplianceFromID(ctx, txn, int64(report.ID))
	if err != nil {
		return fmt.Errorf("find compliance %d: %w", report.ID, err)
	}
	// Don't allow modifying of the submission date if it's set
	if compliance_updates.Has(tablepublicreport.Compliance.Submitted) {
		if compliance.Submitted != nil {
			compliance_updates.Unset(tablepublicreport.Compliance.Submitted)
		} else {
			comm := model.Communication{
				OrganizationID: report.OrganizationID,
				SourceReportID: &report.ID,
			}
			comm, err = querypublic.CommunicationInsert(ctx, txn, comm)
			if err != nil {
				return fmt.Errorf("insert communication: %w", err)
			}
			comm_log := model.CommunicationLogEntry{
				CommunicationID: comm.ID,
				Created:         time.Now(),
				Type:            model.Communicationlogentry_Created,
				User:            nil,
			}
			comm_log, err = querypublic.CommunicationLogEntryInsert(ctx, txn, comm_log)
			if err != nil {
				return fmt.Errorf("insert communication log entry: %w", err)
			}
			log.Debug().Int32("id", comm.ID).Msg("inserted new communication")
		}
	}

	// Avoid attempting to perform an empty update
	err = report_updates.Execute(ctx, txn, int64(report.ID))
	if err != nil {
		return fmt.Errorf("update report: %w", err)
	}
	err = compliance_updates.Execute(ctx, txn, int64(compliance.ReportID))
	if err != nil {
		return fmt.Errorf("update compliance: %w", err)
	}
	if address != nil {
		err = publicReportUpdateAddress(ctx, txn, report, *address)
		if err != nil {
			return fmt.Errorf("update address: %w", err)
		}
	}
	if location != nil {
		err = publicReportUpdateLocation(ctx, txn, report.ID, *location)
		if err != nil {
			return fmt.Errorf("update location: %w", err)
		}
	}
	txn.Commit(ctx)
	return nil
}
func PublicReportReporterUpdated(ctx context.Context, org_id int32, report_id string) {
	event.Updated(event.TypeRMOPublicReport, org_id, report_id)
}
func PublicReportsForOrganization(ctx context.Context, org_id int32, is_public bool) ([]*types.PublicReport, error) {
	return publicreport.ReportsForOrganization(ctx, org_id, is_public)
}
func PublicReportsFromIDs(ctx context.Context, report_ids []int64) ([]modelpublicreport.Report, error) {
	return querypublicreport.ReportsFromIDs(ctx, report_ids)
}
func PublicReportComplianceCreate(ctx context.Context, setter_report modelpublicreport.Report, setter_compliance modelpublicreport.Compliance, org_id int32) (modelpublicreport.Report, error) {
	return publicReportCreate(ctx, setter_report, nil, nil, nil, org_id, func(ctx context.Context, txn db.Ex, report_id int32) error {
		setter_compliance.ReportID = report_id
		_, err := querypublicreport.ComplianceInsert(ctx, txn, setter_compliance)
		if err != nil {
			return fmt.Errorf("Failed to create nuisance database record: %w", err)
		}
		return nil
	})
}
func PublicReportImageCreate(ctx context.Context, public_id string, images []ImageUpload) error {
	txn, err := db.BeginTxn(ctx)
	if err != nil {
		return fmt.Errorf("create txn: %w", err)
	}
	defer txn.Rollback(ctx)

	report, err := querypublicreport.ReportFromPublicID(ctx, db.PGInstance.PGXPool, public_id)
	if err != nil {
		return fmt.Errorf("report from ID: %w", err)
	}
	saved_images, err := saveImageUploads(ctx, txn, images)
	if err != nil {
		return fmt.Errorf("Failed to save image uploads: %w", err)
	}
	if len(saved_images) > 0 {
		report_images := make([]modelpublicreport.ReportImage, len(saved_images))
		for i, image := range saved_images {
			report_images[i] = modelpublicreport.ReportImage{
				ImageID:  image.ID,
				ReportID: report.ID,
			}
		}
		_, err := querypublicreport.ReportImagesInsert(ctx, txn, report_images)
		if err != nil {
			return fmt.Errorf("Failed to save reference to images: %w", err)
		}
		log.Info().Int("len", len(images)).Msg("saved uploaded images")
	}
	txn.Commit(ctx)
	return nil
}
func PublicReportNuisanceCreate(ctx context.Context, setter_report modelpublicreport.Report, setter_nuisance modelpublicreport.Nuisance, location types.Location, address Address, images []ImageUpload) (modelpublicreport.Report, error) {
	return publicReportCreate(ctx, setter_report, &location, &address, images, 0, func(ctx context.Context, txn db.Ex, report_id int32) error {
		setter_nuisance.ReportID = report_id
		_, err := querypublicreport.NuisanceInsert(ctx, txn, setter_nuisance)
		if err != nil {
			return fmt.Errorf("Failed to create nuisance database record: %w", err)
		}
		return nil
	})
}

func PublicReportWaterCreate(ctx context.Context, setter_report modelpublicreport.Report, setter_water modelpublicreport.Water, location types.Location, address Address, images []ImageUpload) (modelpublicreport.Report, error) {
	return publicReportCreate(ctx, setter_report, &location, &address, images, 0, func(ctx context.Context, txn db.Ex, report_id int32) error {
		setter_water.ReportID = report_id
		_, err := querypublicreport.WaterInsert(ctx, txn, setter_water)
		if err != nil {
			return fmt.Errorf("Failed to create water database record: %w", err)
		}
		return nil
	})
}
func PublicReportTypeByID(ctx context.Context, public_id string) (string, error) {
	report, err := querypublicreport.ReportFromPublicID(ctx, db.PGInstance.PGXPool, public_id)
	if err != nil {
		return "", fmt.Errorf("query report '%s': %w", public_id, err)
	}
	return report.ReportType.String(), nil
}

type funcSetReportDetail = func(context.Context, db.Ex, int32) error

func publicReportCreate(ctx context.Context, setter_report modelpublicreport.Report, location *types.Location, address *Address, images []ImageUpload, organization_id int32, detail_setter funcSetReportDetail) (result modelpublicreport.Report, err error) {
	txn, err := db.BeginTxn(ctx)
	if err != nil {
		return result, fmt.Errorf("create txn: %w", err)
	}
	defer txn.Rollback(ctx)

	if setter_report.PublicID == "" {
		public_id, err := GenerateReportID()
		if err != nil {
			return result, fmt.Errorf("create public ID: %w", err)
		}
		setter_report.PublicID = public_id
	}

	var addr *types.Address
	if address != nil {
		if address.GID != "" {
			addr_existing, err := geocode.EnsureAddress(ctx, txn, *address)
			if err != nil {
				return result, fmt.Errorf("Failed to ensure address: %w", err)
			}
			addr = &addr_existing
		} else if address.Raw != "" {
			geo_res, err := geocode.GeocodeRaw(ctx, nil, address.Raw)
			if err != nil {
				return result, fmt.Errorf("Failed to geocode raw: %w", err)
			}
			addr = &geo_res.Address
		} else {
			return result, fmt.Errorf("empty address")
		}
	}

	saved_images, err := saveImageUploads(ctx, txn, images)
	if err != nil {
		return result, fmt.Errorf("Failed to save image uploads: %w", err)
	}
	if organization_id == 0 {
		organization_id, err = matchDistrict(ctx, location, images, addr)
		if err != nil {
			log.Warn().Err(err).Msg("Failed to match district")
		}
	}
	setter_report.OrganizationID = organization_id

	if addr != nil {
		setter_report.AddressID = addr.ID
	}
	result, err = querypublicreport.ReportInsert(ctx, txn, setter_report)
	if err != nil {
		return result, fmt.Errorf("Failed to create report database record: %w", err)
	}
	if location != nil {
		l := *location
		if l.Latitude != 0 && l.Longitude != 0 {
			publicReportUpdateLocation(ctx, txn, result.ID, l)
		}
	}
	log.Info().Str("public_id", setter_report.PublicID).Int32("id", result.ID).Msg("Created base report")

	if len(saved_images) > 0 {
		setters := make([]modelpublicreport.ReportImage, len(saved_images))
		for i, image := range saved_images {
			setters[i] = modelpublicreport.ReportImage{
				ImageID:  int32(image.ID),
				ReportID: int32(result.ID),
			}
		}
		_, err = querypublicreport.ReportImagesInsert(ctx, txn, setters)
		if err != nil {
			return result, fmt.Errorf("Failed to save reference to images: %w", err)
		}
		log.Info().Int("len", len(images)).Msg("saved uploaded images")
	}

	err = detail_setter(ctx, txn, result.ID)
	if err != nil {
		return result, fmt.Errorf("detail setter: %w", err)
	}

	_, err = querypublicreport.ReportLogInsert(ctx, txn, modelpublicreport.ReportLog{
		Created:    time.Now(),
		EmailLogID: nil,
		// ID
		ReportID:  result.ID,
		TextLogID: nil,
		Type:      modelpublicreport.Reportlogtype_Created,
		UserID:    nil,
	})

	// Only create communication entries for compliance when they're submitted
	report_type := setter_report.ReportType
	if report_type != modelpublicreport.Reporttype_Compliance {
		comm := model.Communication{
			OrganizationID: result.OrganizationID,
			SourceReportID: &result.ID,
		}
		comm, err = querypublic.CommunicationInsert(ctx, txn, comm)
		if err != nil {
			return result, fmt.Errorf("insert communication: %w", err)
		}
		log.Debug().Int32("id", comm.ID).Msg("inserted new communication")
	}

	txn.Commit(ctx)

	event.Created(
		event.TypeRMOPublicReport,
		organization_id,
		result.PublicID,
	)
	return result, nil
}
func publicReportUpdateAddress(ctx context.Context, txn db.Tx, report *modelpublicreport.Report, address types.Address) error {
	statement := tablepublicreport.Report.UPDATE(
		tablepublicreport.Report.AddressGid,
		tablepublicreport.Report.AddressRaw,
	).SET(
		postgres.String(address.GID),
		postgres.String(address.Raw),
	).FROM(tablepublic.Address).
		WHERE(
			tablepublicreport.Report.ID.EQ(postgres.Int(int64(report.ID))),
		)
	err := db.ExecuteNoneTx(ctx, txn, statement)

	if err != nil {
		return fmt.Errorf("update report: %w", err)
	}
	statement = tablepublicreport.Report.UPDATE(
		tablepublicreport.Report.AddressID,
	).SET(
		tablepublic.Address.SELECT(
			tablepublic.Address.ID,
		).WHERE(
			tablepublic.Address.Gid.EQ(postgres.String(address.GID)),
		).LIMIT(1),
	).WHERE(
		tablepublicreport.Report.ID.EQ(postgres.Int(int64(report.ID))),
	)
	err = db.ExecuteNoneTx(ctx, txn, statement)
	if err != nil {
		return fmt.Errorf("update report address_id: %w", err)
	}
	return nil
}
func publicReportUpdateLocation(ctx context.Context, txn db.Tx, id int32, location types.Location) error {
	h3cell, _ := location.H3Cell()
	if h3cell == nil {
		return fmt.Errorf("nil h3 cell")
	}
	geom_query, _ := location.GeometryQuery()
	statement := tablepublicreport.Report.UPDATE(
		tablepublicreport.Report.H3cell,
		tablepublicreport.Report.Location,
	).SET(
		postgres.Int(int64(*h3cell)),
		postgres.Raw(geom_query),
	).WHERE(
		tablepublicreport.Report.ID.EQ(postgres.Int(int64(id))),
	)
	err := db.ExecuteNoneTx(ctx, txn, statement)
	if err != nil {
		return fmt.Errorf("Failed to insert publicreport.report geospatial", err)
	}
	return nil
}
