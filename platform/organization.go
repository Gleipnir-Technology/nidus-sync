package platform

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/Gleipnir-Technology/bob/dialect/psql/sm"
	"github.com/Gleipnir-Technology/nidus-sync/db"
	"github.com/Gleipnir-Technology/nidus-sync/db/models"
	"github.com/Gleipnir-Technology/nidus-sync/platform/types"
	//"github.com/google/uuid"
)

type Organization struct {
	ID          int32              `json:"id"`
	ServiceArea *types.ServiceArea `json:"service_area"`

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
func (o Organization) HasServiceArea() bool {
	return o.model.ServiceAreaGeometry.IsValue()
}
func (o Organization) IsCatchall() bool {
	return o.model.IsCatchall
}
func (o Organization) MarshalJSON() ([]byte, error) {
	to_marshal := map[string]any{}
	to_marshal["id"] = o.ID
	to_marshal["name"] = o.Name()
	to_marshal["service_area"] = o.ServiceArea
	return json.Marshal(to_marshal)
}
func (o Organization) Name() string {
	return o.model.Name
}
func (o Organization) IsSyncOngoing() bool {
	return IsSyncOngoing(o.ID)
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

func (o Organization) ServiceRequestRecent(ctx context.Context) ([]*models.FieldseekerServicerequest, error) {
	results, err := o.model.Servicerequests(sm.OrderBy("creationdate").Desc(), sm.Limit(10)).All(ctx, db.PGInstance.BobDB)
	if err != nil {
		return []*models.FieldseekerServicerequest{}, fmt.Errorf("query service request: %w", err)
	}
	return results, nil
}
func (o Organization) Slug() string {
	return o.model.Slug.GetOr("")
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
func OrganizationList(ctx context.Context) ([]*Organization, error) {
	rows, err := models.Organizations.Query().All(ctx, db.PGInstance.BobDB)
	if err != nil {
		return nil, fmt.Errorf("query orgs: %w", err)
	}
	results := make([]*Organization, len(rows))
	for i, row := range rows {
		o := newOrganization(row)
		results[i] = &o
	}
	return results, err
}
func newOrganization(org *models.Organization) Organization {
	var sa *types.ServiceArea
	if org.ServiceAreaXmax.IsValue() &&
		org.ServiceAreaXmin.IsValue() &&
		org.ServiceAreaYmax.IsValue() &&
		org.ServiceAreaYmin.IsValue() {
		sa = &types.ServiceArea{
			Min: types.Location{
				Longitude: org.ServiceAreaXmin.MustGet(),
				Latitude:  org.ServiceAreaYmin.MustGet(),
			},
			Max: types.Location{
				Longitude: org.ServiceAreaXmax.MustGet(),
				Latitude:  org.ServiceAreaYmax.MustGet(),
			},
		}
	}
	return Organization{
		ID:          org.ID,
		ServiceArea: sa,
		model:       org,
	}
}
