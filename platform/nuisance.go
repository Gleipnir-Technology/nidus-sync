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

func NuisanceCreate(ctx context.Context, setter models.PublicreportNuisanceSetter, latlng LatLng, address Address, images []ImageUpload) (public_id string, err error) {
	txn, err := db.PGInstance.BobDB.BeginTx(ctx, nil)
	if err != nil {
		return "", fmt.Errorf("create txn: %w", err)
	}
	defer txn.Rollback(ctx)

	public_id, err = report.GenerateReportID()
	if err != nil {
		return "", fmt.Errorf("create public ID: %w", err)
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
		return "", fmt.Errorf("Failed to save image uploads: %w", err)
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
		setter.OrganizationID = omit.FromPtr(organization_id)
	}
	nuisance, err := models.PublicreportNuisances.Insert(&setter).One(ctx, txn)
	if err != nil {
		return "", fmt.Errorf("Failed to create database record: %w", err)
	}
	if latlng.Latitude != nil && latlng.Longitude != nil {
		h3cell, _ := latlng.H3Cell()
		geom_query, _ := latlng.GeometryQuery()
		_, err = psql.Update(
			um.Table("publicreport.nuisance"),
			um.SetCol("h3cell").ToArg(h3cell),
			um.SetCol("location").To(geom_query),
			um.Where(psql.Quote("id").EQ(psql.Arg(nuisance.ID))),
		).Exec(ctx, txn)
		if err != nil {
			return "", fmt.Errorf("Failed to insert publicreport.nuisance geospatial", err)
		}
	}
	log.Info().Str("public_id", public_id).Int32("id", nuisance.ID).Msg("Created nuisance report")
	if len(saved_images) > 0 {
		setters := make([]*models.PublicreportNuisanceImageSetter, 0)
		for _, image := range saved_images {
			setters = append(setters, &models.PublicreportNuisanceImageSetter{
				ImageID:    omit.From(int32(image.ID)),
				NuisanceID: omit.From(int32(nuisance.ID)),
			})
		}
		_, err = models.PublicreportNuisanceImages.Insert(bob.ToMods(setters...)).Exec(ctx, txn)
		if err != nil {
			return "", fmt.Errorf("Failed to save reference to images: %w", err)
		}
		log.Info().Int("len", len(images)).Msg("saved uploaded images")
	}
	txn.Commit(ctx)

	if organization_id != nil {
		event.Created(
			event.TypeRMONuisance,
			*organization_id,
			nuisance.PublicID,
		)
	}
	return nuisance.PublicID, nil
}
