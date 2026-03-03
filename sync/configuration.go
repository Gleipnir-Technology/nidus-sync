package sync

import (
	"context"
	"net/http"

	"github.com/Gleipnir-Technology/nidus-sync/arcgis"
	"github.com/Gleipnir-Technology/nidus-sync/db"
	"github.com/Gleipnir-Technology/nidus-sync/db/models"
	"github.com/Gleipnir-Technology/nidus-sync/html"
	nhttp "github.com/Gleipnir-Technology/nidus-sync/http"
	//"github.com/rs/zerolog/log"
)

type contentConfigurationRoot struct{}

func getConfigurationRoot(ctx context.Context, r *http.Request, org *models.Organization, user *models.User) (*html.Response[contentConfigurationRoot], *nhttp.ErrorWithStatus) {
	return html.NewResponse("sync/configuration/root.html", contentConfigurationRoot{}), nil
}

type contentSettingOrganization struct {
	Organization *models.Organization
}

type contentSettingIntegration struct {
	ArcGISOAuth *models.ArcgisOauthToken
}

func getConfigurationOrganization(ctx context.Context, r *http.Request, org *models.Organization, u *models.User) (*html.Response[contentSettingOrganization], *nhttp.ErrorWithStatus) {
	org, err := u.Organization().One(ctx, db.PGInstance.BobDB)
	if err != nil {
		return nil, nhttp.NewError("get organization: %w", err)
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
	return html.NewResponse("sync/configuration/organization.html", data), nil
}
func getConfigurationIntegration(ctx context.Context, r *http.Request, org *models.Organization, u *models.User) (*html.Response[contentSettingIntegration], *nhttp.ErrorWithStatus) {
	oauth, err := arcgis.GetOAuthForUser(ctx, u)
	if err != nil {
		return nil, nhttp.NewError("Failed to get oauth: %w", err)
	}
	data := contentSettingIntegration{
		ArcGISOAuth: oauth,
	}
	return html.NewResponse("sync/configuration/integration.html", data), nil
}
func getConfigurationIntegrationArcgis(ctx context.Context, r *http.Request, org *models.Organization, u *models.User) (*html.Response[contentSettingIntegration], *nhttp.ErrorWithStatus) {
	oauth, err := arcgis.GetOAuthForUser(ctx, u)
	if err != nil {
		return nil, nhttp.NewError("Failed to get oauth: %w", err)
	}
	data := contentSettingIntegration{
		ArcGISOAuth: oauth,
	}
	return html.NewResponse("sync/configuration/integration-arcgis.html", data), nil
}

type contentSettingPlaceholder struct{}

func getConfigurationPesticide(ctx context.Context, r *http.Request, org *models.Organization, user *models.User) (*html.Response[contentSettingPlaceholder], *nhttp.ErrorWithStatus) {
	content := contentSettingPlaceholder{}
	return html.NewResponse("sync/configuration/pesticide.html", content), nil
}
func getConfigurationPesticideAdd(ctx context.Context, r *http.Request, org *models.Organization, user *models.User) (*html.Response[contentSettingPlaceholder], *nhttp.ErrorWithStatus) {
	content := contentSettingPlaceholder{}
	return html.NewResponse("sync/configuration/pesticide-add.html", content), nil
}
func getConfigurationUserAdd(ctx context.Context, r *http.Request, org *models.Organization, user *models.User) (*html.Response[contentSettingPlaceholder], *nhttp.ErrorWithStatus) {
	content := contentSettingPlaceholder{}
	return html.NewResponse("sync/configuration/user-add.html", content), nil
}
func getConfigurationUserList(ctx context.Context, r *http.Request, org *models.Organization, user *models.User) (*html.Response[contentSettingPlaceholder], *nhttp.ErrorWithStatus) {
	content := contentSettingPlaceholder{}
	return html.NewResponse("sync/configuration/user-list.html", content), nil
}
