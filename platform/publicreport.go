package platform

import (
	"context"
	"fmt"
	"time"

	//"github.com/Gleipnir-Technology/bob"
	"github.com/Gleipnir-Technology/bob/dialect/psql"
	"github.com/Gleipnir-Technology/bob/dialect/psql/um"
	"github.com/Gleipnir-Technology/nidus-sync/db"
	"github.com/Gleipnir-Technology/nidus-sync/db/enums"
	"github.com/Gleipnir-Technology/nidus-sync/db/models"
	"github.com/Gleipnir-Technology/nidus-sync/platform/event"
	"github.com/rs/zerolog/log"
)

func PublicreportInvalid(ctx context.Context, user User, report_id string) error {
	location, err := models.PublicreportReportLocations.Query(
		models.SelectWhere.PublicreportReportLocations.PublicID.EQ(report_id),
		models.SelectWhere.PublicreportReportLocations.OrganizationID.EQ(user.Organization.ID()),
	).One(ctx, db.PGInstance.BobDB)
	if err != nil {
		return fmt.Errorf("query report existence: %w", err)
	}

	tablename := location.TableName.MustGet()
	_, err = psql.Update(
		um.Table("publicreport."+tablename),
		um.SetCol("reviewed").ToArg(time.Now()),
		um.SetCol("reviewer_id").ToArg(user.ID),
		um.SetCol("status").ToArg(enums.PublicreportReportstatustypeInvalidated),
		um.Where(psql.Quote("public_id").EQ(psql.Arg(report_id))),
	).Exec(ctx, db.PGInstance.BobDB)
	if err != nil {
		return fmt.Errorf("update report %s.%s: %w", tablename, report_id, err)
	}

	log.Info().Str("report-id", report_id).Str("tablename", tablename).Msg("Marked as invalid")
	resource := resourceTypeFromTablename(tablename)
	event.Updated(resource, user.Organization.ID(), report_id)
	return nil
}

func PublicReportReporterUpdated(ctx context.Context, org_id int32, report_id string, tablename string) {
	resource := resourceTypeFromTablename(tablename)
	event.Updated(resource, org_id, report_id)
}
func resourceTypeFromTablename(tablename string) event.ResourceType {
	switch tablename {
	case "nuisance":
		return event.TypeRMONuisance
	case "water":
		return event.TypeRMOWater
	default:
		return event.TypeUnknown
	}
}
