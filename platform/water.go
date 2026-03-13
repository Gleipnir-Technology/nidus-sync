package platform

import (
	"context"
	"fmt"

	"github.com/Gleipnir-Technology/bob"
	"github.com/Gleipnir-Technology/bob/dialect/psql"
	"github.com/Gleipnir-Technology/bob/dialect/psql/um"
	"github.com/Gleipnir-Technology/nidus-sync/db"
	"github.com/Gleipnir-Technology/nidus-sync/db/models"
	"github.com/Gleipnir-Technology/nidus-sync/platform/event"
	"github.com/Gleipnir-Technology/nidus-sync/platform/geocode"
	"github.com/Gleipnir-Technology/nidus-sync/platform/report"
	"github.com/Gleipnir-Technology/nidus-sync/platform/types"
	"github.com/aarondl/opt/omit"
	"github.com/aarondl/opt/omitnull"
	"github.com/rs/zerolog/log"
)

func WaterCreate(ctx context.Context, setter models.PublicreportWaterSetter, latlng LatLng, address Address, images []ImageUpload) (public_id string, err error) {
	txn, err := db.PGInstance.BobDB.BeginTx(ctx, nil)
	if err != nil {
		return "", fmt.Errorf("Failed to create transaction: %w", err)
	}
	defer txn.Rollback(ctx)
	public_id, err = report.GenerateReportID()
	if err != nil {
		return "", fmt.Errorf("Failed to create water report public ID", err)
	}
	setter.PublicID = omit.From(public_id)

	// If we've got an locality value it was set by geocoding so we should save it
	var a *models.Address
	if address.Locality != "" && latlng.Latitude != nil && latlng.Longitude != nil {
		a, err = geocode.EnsureAddress(ctx, txn, address, types.Location{
			Latitude:  *latlng.Latitude,
			Longitude: *latlng.Longitude,
		})
		if err != nil {
			return "", fmt.Errorf("Failed to ensure address: %w", err)
		}
	}

	saved_images, err := saveImageUploads(ctx, txn, images)
	if err != nil {
		return "", fmt.Errorf("Failed to save image uploads", err)
	}

	var organization_id *int32
	organization_id, err = MatchDistrict(ctx, latlng.Longitude, latlng.Latitude, images)
	if err != nil {
		log.Warn().Err(err).Msg("Failed to match district")
	}
	if a != nil {
		setter.AddressID = omitnull.From(a.ID)
	}
	if organization_id != nil {
		setter.OrganizationID = omitnull.FromPtr(organization_id)
	}

	water, err := models.PublicreportWaters.Insert(&setter).One(ctx, txn)
	if err != nil {
		return "", fmt.Errorf("Failed to create database record", err)
	}

	if latlng.Latitude != nil && latlng.Longitude != nil {
		h3cell, _ := latlng.H3Cell()
		geom_query, _ := latlng.GeometryQuery()
		_, err = psql.Update(
			um.Table("publicreport.water"),
			um.SetCol("h3cell").ToArg(h3cell),
			um.SetCol("location").To(geom_query),
			um.Where(psql.Quote("id").EQ(psql.Arg(water.ID))),
		).Exec(ctx, txn)
		if err != nil {
			return "", fmt.Errorf("Failed to update publicreport.water geospatial", err)
		}
	}
	log.Info().Int32("id", water.ID).Str("public_id", water.PublicID).Msg("Created water report")
	setters := make([]*models.PublicreportWaterImageSetter, 0)
	for _, image := range saved_images {
		setters = append(setters, &models.PublicreportWaterImageSetter{
			ImageID: omit.From(int32(image.ID)),
			WaterID: omit.From(int32(water.ID)),
		})
	}
	if len(setters) > 0 {
		_, err = models.PublicreportWaterImages.Insert(bob.ToMods(setters...)).Exec(ctx, txn)
		if err != nil {
			return "", fmt.Errorf("Failed to save upload relationships", err)
		}
	}
	txn.Commit(ctx)

	if organization_id != nil {
		event.Created(
			event.TypeRMOWater,
			*organization_id,
			water.PublicID,
		)
	}
	return water.PublicID, nil
}
