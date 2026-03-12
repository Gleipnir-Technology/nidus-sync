package platform

import (
	"context"
	"fmt"

	"github.com/Gleipnir-Technology/bob/dialect/psql/sm"
	"github.com/Gleipnir-Technology/nidus-sync/db"
	"github.com/Gleipnir-Technology/nidus-sync/db/models"
	"github.com/Gleipnir-Technology/nidus-sync/platform/background"
	//"github.com/google/uuid"
)

type Organization struct {
	ServiceAreaXmax float64
	ServiceAreaXmin float64
	ServiceAreaYmax float64
	ServiceAreaYmin float64

	model *models.Organization
}

func (o Organization) ArcgisAccountID() string {
	if o.model.ArcgisAccountID.IsNull() {
		return ""
	}
	return o.model.ArcgisAccountID.MustGet()
}
func (o Organization) CountServiceRequest(ctx context.Context) (uint, error) {
	result, err := o.model.Servicerequests().Count(ctx, db.PGInstance.BobDB)
	if err != nil {
		return 0, fmt.Errorf("get service request count: %w", err)
	}
	return uint(result), nil
}
func (o Organization) CountSource(ctx context.Context) (uint, error) {
	result, err := o.model.Pointlocations().Count(ctx, db.PGInstance.BobDB)
	if err != nil {
		return 0, fmt.Errorf("get source count: %w", err)
	}
	return uint(result), nil
}
func (o Organization) CountTrap(ctx context.Context) (uint, error) {
	result, err := o.model.Traplocations().Count(ctx, db.PGInstance.BobDB)
	if err != nil {
		return 0, fmt.Errorf("get trap count: %w", err)
	}
	return uint(result), nil
}
func (o Organization) Name() string {
	return o.model.Name
}
func (o Organization) ID() int32 {
	return o.model.ID
}
func (o Organization) IsSyncOngoing() bool {
	return background.IsSyncOngoing(o.ID())
}
func (o Organization) FieldseekerSyncLatest(ctx context.Context) (*models.FieldseekerSync, error) {
	sync, err := o.model.FieldseekerSyncs(sm.OrderBy("created").Desc()).One(ctx, db.PGInstance.BobDB)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			return nil, nil
		}
		return nil, fmt.Errorf("get syncs: %w", err)
	}
	return sync, nil
}

type ServiceArea struct {
	Min Point
	Max Point
}

func (o Organization) ServiceArea() ServiceArea {
	if o.model.ServiceAreaXmax.IsNull() ||
		o.model.ServiceAreaXmin.IsNull() ||
		o.model.ServiceAreaYmax.IsNull() ||
		o.model.ServiceAreaYmin.IsNull() {
		return ServiceArea{}
	}
	return ServiceArea{
		Min: Point{
			X: o.model.ServiceAreaXmin.MustGet(),
			Y: o.model.ServiceAreaYmin.MustGet(),
		},
		Max: Point{
			X: o.model.ServiceAreaXmax.MustGet(),
			Y: o.model.ServiceAreaYmax.MustGet(),
		},
	}
}
func (o Organization) ServiceRequestRecent(ctx context.Context) ([]*models.FieldseekerServicerequest, error) {
	results, err := o.model.Servicerequests(sm.OrderBy("creationdate").Desc(), sm.Limit(10)).All(ctx, db.PGInstance.BobDB)
	if err != nil {
		return []*models.FieldseekerServicerequest{}, fmt.Errorf("query service request: %w", err)
	}
	return results, nil
}
func OrganizationByID(ctx context.Context, id int) (*Organization, error) {
	org, err := models.FindOrganization(ctx, db.PGInstance.BobDB, int32(id))
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			return nil, nil
		}
		return nil, fmt.Errorf("query org: %w", err)
	}
	o := newOrganization(org)
	return &o, nil
}
func newOrganization(org *models.Organization) Organization {
	return Organization{
		model: org,
	}
}
