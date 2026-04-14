package platform

import (
	"context"
	"fmt"

	"github.com/Gleipnir-Technology/bob"
	"github.com/Gleipnir-Technology/bob/dialect/psql"
	//"github.com/Gleipnir-Technology/bob/dialect/psql/dialect"
	"github.com/Gleipnir-Technology/bob/dialect/psql/sm"
	"github.com/Gleipnir-Technology/nidus-sync/db"
	//"github.com/Gleipnir-Technology/nidus-sync/db/models"
	//"github.com/Gleipnir-Technology/nidus-sync/platform/fieldseeker"
	"github.com/Gleipnir-Technology/nidus-sync/platform/types"
	//"github.com/google/uuid"
	"github.com/stephenafamo/scan"
)

func ServiceRequestList(ctx context.Context, user User, limit int) ([]*types.ServiceRequest, error) {
	query := psql.Select(
		sm.Columns(
			"COALESCE(sr.reqaddr1, '') AS \"address.raw\"",
			"COALESCE(sr.assignedtech, '') AS \"assigned_technician\"",
			"COALESCE(sr.reqcity, '') AS \"city\"",
			"sr.creationdate AS \"created\"",
			//"COALESCE(sr.h3cell, 0) AS \"h3cell\"",
			"COALESCE(sr.dog, 0) AS \"has_dog\"",
			"COALESCE(sr.spanish, 0) AS \"has_spanish_speaker\"",
			"sr.globalid AS \"id\"",
			"sr.priority AS \"priority\"",
			"sr.recdatetime AS \"recorded_date\"",
			"sr.source AS \"source\"",
			"sr.reqtarget AS \"target\"",
			"sr.reqzip AS \"zip\"",
			"COALESCE(ST_X(pl.geospatial), 0) AS \"address.location.longitude\"",
			"COALESCE(ST_Y(pl.geospatial), 0) AS \"address.location.latitude\"",
		),
		sm.From("fieldseeker.servicerequest").As("sr"),
		sm.LeftJoin("fieldseeker.pointlocation").As("pl").OnEQ(
			psql.Quote("sr", "pointlocid"),
			psql.Quote("pl", "globalid"),
		),
	)
	results, err := bob.All(ctx, db.PGInstance.BobDB, query, scan.StructMapper[*types.ServiceRequest]())
	if err != nil {
		return nil, fmt.Errorf("query service requests: %w", err)
	}
	/*
		service_requests, err := models.FieldseekerServicerequests.Query(
			models.SelectWhere.FieldseekerServicerequests.OrganizationID.EQ(user.Organization.ID),
			//sm.OrderBy(models.FieldseekerServicerequests.Columns.Created).Desc(),
		).All(ctx, db.PGInstance.BobDB)
		if err != nil {
			return nil, fmt.Errorf("query sync: %w", err)
		}
		point_location_ids := make([]uuid.UUID, len(service_requests))
		for i, s := range service_requests {
			p, ok := s.Pointlocid.Get()
			if ok {
				point_location_ids[i] = p
			}
		}
		point_locations, err := fieldseeker.PointLocationList(ctx, point_location_ids)
		if err != nil {
			return nil, fmt.Errorf("list point locations: %w", err)
		}
		point_location_by_id := make(map[uuid.UUID]*models.FieldseekerPointlocation, len(point_locations))
		for _, pl := range point_locations {
			point_location_by_id[pl.Globalid] = pl
		}
		results := make([]*types.ServiceRequest, len(service_requests))
		for i, s := range service_requests {
			r := types.ServiceRequestFromModel(s)
			loc_id, ok := s.Pointlocid.Get()
			if ok {
				pl, ok := point_location_by_id[loc_id]
				if ok {
					r.Location = types.LocationFromFS(pl)
				}
			}
			results[i] = &r
			point_location_ids[i]
		}
	*/
	return results, nil
}
