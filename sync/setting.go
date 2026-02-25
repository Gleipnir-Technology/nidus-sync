package sync

import (
	"context"
	"net/http"

	"github.com/Gleipnir-Technology/nidus-sync/arcgis"
	"github.com/Gleipnir-Technology/nidus-sync/db"
	"github.com/Gleipnir-Technology/nidus-sync/db/models"
	//"github.com/rs/zerolog/log"
)

type contentSettingOrganization struct {
	Organization *models.Organization
}

type contentSettingIntegration struct {
	ArcGISOAuth *models.OauthToken
}

type contentAuthenticatedPlaceholder struct {
}

func getSetting(ctx context.Context, r *http.Request, org *models.Organization, u *models.User) (*response[contentAuthenticatedPlaceholder], *errorWithStatus) {
	data := contentAuthenticatedPlaceholder{}
	return newResponse("sync/settings.html", data), nil
}
func getSettingOrganization(ctx context.Context, r *http.Request, org *models.Organization, u *models.User) (*response[contentSettingOrganization], *errorWithStatus) {
	org, err := u.Organization().One(ctx, db.PGInstance.BobDB)
	if err != nil {
		return nil, newError("get organization: %w", err)
	}
	/*
		var district contentDistrict
		district, err = bob.One[contentDistrict](ctx, db.PGInstance.BobDB, psql.Select(
			sm.From("import.district"),
			sm.Columns(
				"address",
				"agency",
				"area_4326_sqm",
				"city1",
				"city2",
				"contact",
				"fax1",
				"general_mg",
				"gid",
				"phone1",
				"phone2",
				"postal_c_1",
				"website",
				psql.F("ST_AsGeoJSON", "centroid_4326"),
				psql.F("ST_XMin", "extent_4326"),
				psql.F("ST_YMin", "extent_4326"),
				psql.F("ST_XMax", "extent_4326"),
				psql.F("ST_YMax", "extent_4326"),
			),
			sm.Where(psql.Quote("gid").EQ(psql.Arg(gid))),
		), scan.StructMapper[contentDistrict]())
		if err != nil {
			respondError(w, "Failed to get extents", err, http.StatusInternalServerError)
			return
		}
	*/
	data := contentSettingOrganization{
		Organization: org,
	}
	return newResponse("sync/setting-organization.html", data), nil
}
func getSettingIntegration(ctx context.Context, r *http.Request, org *models.Organization, u *models.User) (*response[contentSettingIntegration], *errorWithStatus) {
	oauth, err := arcgis.GetOAuthForUser(ctx, u)
	if err != nil {
		return nil, newError("Failed to get oauth: %w", err)
	}
	data := contentSettingIntegration{
		ArcGISOAuth: oauth,
	}
	return newResponse("sync/setting-integration.html", data), nil
}
func getSettingIntegrationArcgis(ctx context.Context, r *http.Request, org *models.Organization, u *models.User) (*response[contentSettingIntegration], *errorWithStatus) {
	oauth, err := arcgis.GetOAuthForUser(ctx, u)
	if err != nil {
		return nil, newError("Failed to get oauth: %w", err)
	}
	data := contentSettingIntegration{
		ArcGISOAuth: oauth,
	}
	return newResponse("sync/setting-integration-arcgis.html", data), nil
}

type contentSettingPlaceholder struct{}

func getSettingPesticide(ctx context.Context, r *http.Request, org *models.Organization, user *models.User) (*response[contentSettingPlaceholder], *errorWithStatus) {
	content := contentSettingPlaceholder{}
	return newResponse("sync/setting-pesticide.html", content), nil
}
func getSettingPesticideAdd(ctx context.Context, r *http.Request, org *models.Organization, user *models.User) (*response[contentSettingPlaceholder], *errorWithStatus) {
	content := contentSettingPlaceholder{}
	return newResponse("sync/setting-pesticide-add.html", content), nil
}
func getSettingUserAdd(ctx context.Context, r *http.Request, org *models.Organization, user *models.User) (*response[contentSettingPlaceholder], *errorWithStatus) {
	content := contentSettingPlaceholder{}
	return newResponse("sync/setting-user-add.html", content), nil
}
func getSettingUserList(ctx context.Context, r *http.Request, org *models.Organization, user *models.User) (*response[contentSettingPlaceholder], *errorWithStatus) {
	content := contentSettingPlaceholder{}
	return newResponse("sync/setting-user-list.html", content), nil
}
