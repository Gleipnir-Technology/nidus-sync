package publicreport

import (
	"context"
	"fmt"

	"github.com/Gleipnir-Technology/bob"
	"github.com/Gleipnir-Technology/bob/dialect/psql"
	"github.com/Gleipnir-Technology/bob/dialect/psql/sm"
	//"github.com/Gleipnir-Technology/nidus-sync/config"
	"github.com/Gleipnir-Technology/nidus-sync/db"
	//"github.com/Gleipnir-Technology/nidus-sync/db/models"
	"github.com/Gleipnir-Technology/nidus-sync/platform/types"
	//"github.com/google/uuid"
	//"github.com/rs/zerolog/log"
	"github.com/stephenafamo/scan"
)

/*
SELECT
    i.*,
    MAX(e.value) FILTER (WHERE e.name = 'Make') as exif_make,
    MAX(e.value) FILTER (WHERE e.name = 'Model') as exif_model,
    MAX(e.value) FILTER (WHERE e.name = 'DateTime') as exif_datetime,
    MAX(e.value) FILTER (WHERE e.name = 'GPSLatitude') as exif_gps_lat
FROM publicreport.image i
LEFT JOIN publicreport.image_exif e ON i.id = e.image_id
WHERE i.id IN (1, 2, 3, 4)
GROUP BY i.id;
*/
// Get all the images that belong to the list of nuisance report IDs
func loadImagesForReportNuisance(ctx context.Context, org_id int32, report_ids []int32) (results map[int32][]types.Image, err error) {
	rows, err := bob.All(ctx, db.PGInstance.BobDB, psql.Select(
		sm.Columns(
			"i.storage_uuid AS uuid",
			"COALESCE(ST_X(i.location), 0) AS \"location.longitude\"",
			"COALESCE(ST_Y(i.location), 0) AS \"location.latitude\"",
			"ST_Distance(i.location::geography, n.location::geography) AS \"distance_from_reporter_meters\"",
			"COALESCE(MAX(e.value) FILTER (WHERE e.name = 'Make'), '') AS exif_make",
			"COALESCE(MAX(e.value) FILTER (WHERE e.name = 'Model'), '') AS exif_model",
			"COALESCE(MAX(e.value) FILTER (WHERE e.name = 'DateTime'), '') AS exif_datetime",
			"ni.nuisance_id AS report_id",
		),
		sm.From("publicreport.image").As("i"),
		sm.LeftJoin("publicreport.image_exif").As("e").OnEQ(
			psql.Quote("i", "id"),
			psql.Quote("e", "image_id"),
		),
		sm.InnerJoin("publicreport.nuisance_image").As("ni").OnEQ(
			psql.Quote("ni", "image_id"),
			psql.Quote("i", "id"),
		),
		sm.InnerJoin("publicreport.nuisance").As("n").OnEQ(
			psql.Quote("ni", "nuisance_id"),
			psql.Quote("n", "id"),
		),
		sm.Where(psql.Quote("ni", "nuisance_id").EQ(psql.Any(report_ids))),
		sm.GroupBy(
			//psql.Quote("i", "id"),
			//psql.Quote("ni", "nuisance_id"),
			psql.Raw("i.id, ni.nuisance_id, n.location"),
		),
	), scan.StructMapper[types.Image]())
	if err != nil {
		return nil, fmt.Errorf("get images: %w", err)
	}
	results = make(map[int32][]types.Image, len(report_ids))
	for _, row := range rows {
		r, ok := results[row.ReportID]
		if !ok {
			r = make([]types.Image, 0)
		}
		r = append(r, row)
		results[row.ReportID] = r
	}
	return results, nil
}

// Get all the images that belong to the list of water report IDs
func loadImagesForReportWater(ctx context.Context, org_id int32, report_ids []int32) (results map[int32][]types.Image, err error) {
	rows, err := bob.All(ctx, db.PGInstance.BobDB, psql.Select(
		sm.Columns(
			"i.storage_uuid AS uuid",
			"COALESCE(ST_X(i.location), 0) AS \"location.longitude\"",
			"COALESCE(ST_Y(i.location), 0) AS \"location.latitude\"",
			"ST_Distance(i.location::geography, w.location::geography) AS \"distance_from_reporter_meters\"",
			"COALESCE(MAX(e.value) FILTER (WHERE e.name = 'Make'), '') AS exif_make",
			"COALESCE(MAX(e.value) FILTER (WHERE e.name = 'Model'), '') AS exif_model",
			"COALESCE(MAX(e.value) FILTER (WHERE e.name = 'DateTime'), '') AS exif_datetime",
			"wi.water_id AS report_id",
		),
		sm.From("publicreport.image").As("i"),
		sm.LeftJoin("publicreport.image_exif").As("e").OnEQ(
			psql.Quote("i", "id"),
			psql.Quote("e", "image_id"),
		),
		sm.InnerJoin("publicreport.water_image").As("wi").OnEQ(
			psql.Quote("wi", "image_id"),
			psql.Quote("i", "id"),
		),
		sm.InnerJoin("publicreport.water").As("w").OnEQ(
			psql.Quote("wi", "water_id"),
			psql.Quote("w", "id"),
		),
		sm.Where(psql.Quote("wi", "water_id").EQ(psql.Any(report_ids))),
		sm.GroupBy(
			//psql.Quote("i", "id"),
			//psql.Quote("ni", "nuisance_id"),
			psql.Raw("i.id, wi.water_id, w.location"),
		),
	), scan.StructMapper[types.Image]())
	if err != nil {
		return nil, fmt.Errorf("get images: %w", err)
	}
	results = make(map[int32][]types.Image, len(report_ids))
	for _, row := range rows {
		r, ok := results[row.ReportID]
		if !ok {
			r = make([]types.Image, 0)
		}
		r = append(r, row)
		results[row.ReportID] = r
	}
	return results, nil
}
