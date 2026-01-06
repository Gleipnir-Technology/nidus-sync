package platform

import (
	"context"
	"fmt"
	"time"

	"github.com/Gleipnir-Technology/nidus-sync/db"
	"github.com/Gleipnir-Technology/nidus-sync/db/models"
	"github.com/google/uuid"
)

func fieldseeker(ctx context.Context, u *models.User, since *time.Time) (fsync FieldseekerRecordsSync, err error) {
	if u == nil {
		return fsync, fmt.Errorf("Wha! Nil user!")
	}
	org := u.R.Organization
	if org == nil {
		return fsync, fmt.Errorf("Whoa nil org from user %d and org %d.", u.ID, u.OrganizationID)
	}
	db_connection := db.PGInstance.BobDB
	pl, err := org.Pointlocations().All(ctx, db_connection)
	if err != nil {
		return fsync, fmt.Errorf("Failed to get point locations: %w", err)
	}
	inspections, err := u.R.Organization.Mosquitoinspections().All(ctx, db.PGInstance.BobDB)
	if err != nil {
		return fsync, fmt.Errorf("Failed to get mosquito inspections: %w", err)
	}
	inspections_by_location := make(map[uuid.UUID]models.FieldseekerMosquitoinspectionSlice, 0)
	for _, i := range inspections {
		if i.Pointlocid.IsNull() {
			continue
		}
		locid := i.Pointlocid.MustGet()
		insp, ok := inspections_by_location[locid]
		if !ok {
			insp = make(models.FieldseekerMosquitoinspectionSlice, 0)
		}
		insp = append(insp, i)
		inspections_by_location[locid] = insp
	}
	treatments, err := u.R.Organization.Treatments().All(ctx, db.PGInstance.BobDB)
	if err != nil {
		return fsync, fmt.Errorf("Failed to get treatment data: %w", err)
	}
	treatments_by_location := make(map[uuid.UUID]models.FieldseekerTreatmentSlice, 0)
	for _, t := range treatments {
		if t.Pointlocid.IsNull() {
			continue
		}
		locid := t.Pointlocid.MustGet()
		ts, ok := treatments_by_location[locid]
		if !ok {
			ts = make(models.FieldseekerTreatmentSlice, 0)
		}
		ts = append(ts, t)
		treatments_by_location[locid] = ts
	}
	sources := make([]*MosquitoSource, 0)
	for _, p := range pl {
		inspections, ok := inspections_by_location[p.Globalid]
		if !ok {
			inspections = make(models.FieldseekerMosquitoinspectionSlice, 0)
		}
		treatments, ok := treatments_by_location[p.Globalid]
		if !ok {
			treatments = make(models.FieldseekerTreatmentSlice, 0)
		}
		ms := MosquitoSource{
			PointLocation: p,
			Inspections:   &inspections,
			Treatments:    &treatments,
		}
		sources = append(sources, &ms)
	}
	fsync.MosquitoSources = &sources
	return fsync, err
}

func ContentClientIos(ctx context.Context, u *models.User, since *time.Time) (csync ClientSync, err error) {
	fsync, err := fieldseeker(ctx, u, since)
	return ClientSync{
		Fieldseeker: fsync,
	}, err
}
