package sync

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/Gleipnir-Technology/nidus-sync/db"
	"github.com/Gleipnir-Technology/nidus-sync/db/models"
	"github.com/Gleipnir-Technology/nidus-sync/db/sql"
	"github.com/Gleipnir-Technology/nidus-sync/notification"
	"github.com/google/uuid"
	"github.com/stephenafamo/bob"
	"github.com/stephenafamo/bob/dialect/psql"
	"github.com/stephenafamo/bob/dialect/psql/sm"
	"github.com/uber/h3-go/v4"
)

func breedingSourcesByCell(ctx context.Context, org *models.Organization, c h3.Cell) ([]BreedingSourceSummary, error) {
	var results []BreedingSourceSummary

	boundary, err := c.Boundary()
	if err != nil {
		return results, fmt.Errorf("Failed to get cell boundary: %w", err)
	}
	geom_query := gisStatement(boundary)
	rows, err := org.Pointlocations(
		sm.Where(
			psql.F("ST_Within", "geospatial", geom_query),
		),
		sm.OrderBy("lasttreatdate"),
	).All(ctx, db.PGInstance.BobDB)
	if err != nil {
		return results, fmt.Errorf("Failed to query rows: %w", err)
	}
	for _, r := range rows {
		var last_inspected *time.Time
		if !r.Lastinspectdate.IsNull() {
			l := r.Lastinspectdate.MustGet()
			last_inspected = &l
		}
		var last_treat_date *time.Time
		if !r.Lasttreatdate.IsNull() {
			l := r.Lasttreatdate.MustGet()
			last_treat_date = &l
		}
		results = append(results, BreedingSourceSummary{
			ID:            r.Globalid,
			LastInspected: last_inspected,
			LastTreated:   last_treat_date,
			Type:          r.Habitat.GetOr("none"),
		})
	}
	return results, nil
}
func gisStatement(cb h3.CellBoundary) string {
	var content strings.Builder
	for i, p := range cb {
		if i != 0 {
			content.WriteString(", ")
		}
		content.WriteString(fmt.Sprintf("%f %f", p.Lng, p.Lat))
	}
	// Repeat the first coordinate to close the polygon
	content.WriteString(fmt.Sprintf(", %f %f", cb[0].Lng, cb[0].Lat))
	return fmt.Sprintf("ST_GeomFromText('POLYGON((%s))', 3857)", content.String())
}

func sourceByGlobalId(ctx context.Context, org *models.Organization, id uuid.UUID) (*BreedingSourceDetail, error) {
	row, err := org.Pointlocations(
		sm.Where(models.FieldseekerPointlocations.Columns.Globalid.EQ(psql.Arg(id))),
	).One(ctx, db.PGInstance.BobDB)
	if err != nil {
		return nil, fmt.Errorf("Failed to get point location: %w", err)
	}
	return toTemplateBreedingSource(row), nil
}

func extractInitials(name string) string {
	parts := strings.Fields(name)
	var initials strings.Builder

	for _, part := range parts {
		if len(part) > 0 {
			initials.WriteString(strings.ToUpper(string(part[0])))
		}
	}

	return initials.String()
}

func contentForUser(ctx context.Context, user *models.User) (User, error) {
	notifications, err := notification.ForUser(ctx, user)
	if err != nil {
		return User{}, err
	}
	return User{
		DisplayName:   user.DisplayName,
		Initials:      extractInitials(user.DisplayName),
		Notifications: notifications,
		Username:      user.Username,
	}, nil

}

func trapsBySource(ctx context.Context, org *models.Organization, sourceID uuid.UUID) ([]TrapNearby, error) {
	locations, err := sql.TrapLocationBySourceID(org.ID, sourceID).All(ctx, db.PGInstance.BobDB)
	if err != nil {
		return nil, fmt.Errorf("Failed to query rows: %w", err)
	}

	location_ids := make([]uuid.UUID, 0)
	var args []bob.Expression
	for _, location := range locations {
		location_ids = append(location_ids, location.TrapLocationGlobalid)
		args = append(args, psql.Arg(location.TrapLocationGlobalid))
	}
	/*
		trap_data, err := org.FSTrapdata(
			sm.Where(
				models.FSTrapdata.Columns.LocID.In(args...),
			),
			sm.OrderBy("enddatetime"),
		).All(ctx, db.PGInstance.BobDB)
	*/

	/*
		query := org.FSTrapdata(
			sm.From(
				psql.Select(
					sm.From(psql.F("ROW_NUMBER")(
						fm.Over(
							wm.PartitionBy(models.FSTrapdata.Columns.LocID),
							wm.OrderBy(models.FSTrapdata.Columns.Enddatetime).Desc(),
						),
					)).As("row_num"),
				sm.Where(models.FSTrapdata.Columns.LocID.In(args...))),
			),
			sm.Where(psql.Quote("row_num").LTE(psql.Arg(10))),
			sm.OrderBy(models.FSTrapdata.Columns.LocID),
			sm.OrderBy(models.FSTrapdata.Columns.Enddatetime).Desc(),
		)
	*/
	/*
		query := psql.Select(
			sm.From(
				psql.Select(
					sm.Columns(
						models.FSTrapdata.Columns.Globalid,
						psql.F("ROW_NUMBER")(
						fm.Over(
							wm.PartitionBy(models.FSTrapdata.Columns.LocID),
							wm.OrderBy(models.FSTrapdata.Columns.Enddatetime).Desc(),
						),
					).As("row_num"),
					sm.From(models.FSTrapdata.Name()),
				),
				sm.Where(models.FSTrapdata.Columns.LocID.In(args...))),
			),
			sm.Where(psql.Quote("row_num").LTE(psql.Arg(10))),
			sm.OrderBy(models.FSTrapdata.Columns.LocID),
			sm.OrderBy(models.FSTrapdata.Columns.Enddatetime).Desc(),
		)
		log.Info().Str("trapdata", queryToString(query)).Msg("Getting trap data")
		trap_data, err := query.Exec(ctx, db.PGInstance.BobDB)
	*/

	trap_data, err := sql.TrapDataByLocationIDRecent(org.ID, location_ids).All(ctx, db.PGInstance.BobDB)
	if err != nil {
		return nil, fmt.Errorf("Failed to query trap data: %w", err)
	}

	counts, err := sql.TrapCountByLocationID(org.ID, location_ids).All(ctx, db.PGInstance.BobDB)
	if err != nil {
		return nil, fmt.Errorf("Failed to query trap counts: %w", err)
	}

	traps, err := toTemplateTraps(locations, trap_data, counts)
	if err != nil {
		return nil, fmt.Errorf("Failed to convert trap data: %w", err)
	}
	return traps, nil
}

func treatmentsBySource(ctx context.Context, org *models.Organization, sourceID uuid.UUID) ([]Treatment, error) {
	var results []Treatment
	rows, err := org.Treatments(
		sm.Where(
			models.FieldseekerTreatments.Columns.Pointlocid.EQ(psql.Arg(sourceID)),
		),
		sm.OrderBy("enddatetime").Desc(),
	).All(ctx, db.PGInstance.BobDB)
	if err != nil {
		return results, fmt.Errorf("Failed to query rows: %w", err)
	}
	//log.Info().Int("row count", len(rows)).Msg("Getting treatments")
	return toTemplateTreatment(rows)
}

func treatmentsByCell(ctx context.Context, org *models.Organization, c h3.Cell) ([]Treatment, error) {
	var results []Treatment
	boundary, err := c.Boundary()
	if err != nil {
		return results, fmt.Errorf("Failed to get cell boundary: %w", err)
	}
	geom_query := gisStatement(boundary)
	rows, err := org.Treatments(
		sm.Where(
			psql.F("ST_Within", "geospatial", geom_query),
		),
		sm.OrderBy("pointlocid"),
		sm.OrderBy("enddatetime"),
	).All(ctx, db.PGInstance.BobDB)
	if err != nil {
		return results, fmt.Errorf("Failed to query rows: %w", err)
	}
	return toTemplateTreatment(rows)
}
func inspectionsByCell(ctx context.Context, org *models.Organization, c h3.Cell) ([]Inspection, error) {
	var results []Inspection

	boundary, err := c.Boundary()
	if err != nil {
		return results, fmt.Errorf("Failed to get cell boundary: %w", err)
	}
	geom_query := gisStatement(boundary)
	rows, err := org.Mosquitoinspections(
		sm.Where(
			psql.F("ST_Within", "geospatial", geom_query),
		),
		sm.OrderBy("pointlocid"),
		sm.OrderBy("enddatetime"),
	).All(ctx, db.PGInstance.BobDB)
	if err != nil {
		return results, fmt.Errorf("Failed to query rows: %w", err)
	}
	return toTemplateInspection(rows)
}
func inspectionsBySource(ctx context.Context, org *models.Organization, sourceID uuid.UUID) ([]Inspection, error) {
	var results []Inspection

	rows, err := org.Mosquitoinspections(
		sm.Where(
			models.FieldseekerMosquitoinspections.Columns.Pointlocid.EQ(psql.Arg(sourceID)),
		),
		sm.OrderBy("enddatetime").Desc(),
	).All(ctx, db.PGInstance.BobDB)
	if err != nil {
		return results, fmt.Errorf("Failed to query rows: %w", err)
	}
	return toTemplateInspection(rows)
}
