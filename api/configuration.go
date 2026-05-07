package api

import (
	"context"
	"net/http"

	"github.com/Gleipnir-Technology/bob/dialect/psql"
	"github.com/Gleipnir-Technology/bob/dialect/psql/um"
	"github.com/Gleipnir-Technology/nidus-sync/db"
	"github.com/Gleipnir-Technology/nidus-sync/db/gen/nidus-sync/arcgis/model"
	queryarcgis "github.com/Gleipnir-Technology/nidus-sync/db/query/arcgis"
	"github.com/Gleipnir-Technology/nidus-sync/html"
	nhttp "github.com/Gleipnir-Technology/nidus-sync/http"
	"github.com/Gleipnir-Technology/nidus-sync/platform"
	"github.com/rs/zerolog/log"
)

type contentConfigurationRoot struct{}

func getConfigurationRoot(ctx context.Context, r *http.Request, user platform.User) (*html.Response[contentConfigurationRoot], *nhttp.ErrorWithStatus) {
	return html.NewResponse("sync/configuration/root.html", contentConfigurationRoot{}), nil
}

type contentSettingOrganization struct {
	Organization platform.Organization
}

type contentSettingIntegration struct {
	ArcGISAccount *model.Account
	ArcGISOAuth   *model.OAuthToken
	ServiceMaps   []model.ServiceMap
}

func getConfigurationOrganization(ctx context.Context, r *http.Request, u platform.User) (*html.Response[contentSettingOrganization], *nhttp.ErrorWithStatus) {
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
		Organization: u.Organization,
	}
	return html.NewResponse("sync/configuration/organization.html", data), nil
}
func getConfigurationIntegration(ctx context.Context, r *http.Request, u platform.User) (*html.Response[contentSettingIntegration], *nhttp.ErrorWithStatus) {
	oauth, err := platform.GetOAuthForUser(ctx, u)
	if err != nil {
		return nil, nhttp.NewError("Failed to get oauth: %w", err)
	}
	data := contentSettingIntegration{
		ArcGISOAuth: oauth,
	}
	return html.NewResponse("sync/configuration/integration.html", data), nil
}
func getConfigurationIntegrationArcgis(ctx context.Context, r *http.Request, u platform.User) (*html.Response[contentSettingIntegration], *nhttp.ErrorWithStatus) {
	oauth, err := platform.GetOAuthForUser(ctx, u)
	if err != nil {
		return nil, nhttp.NewError("Failed to get oauth: %w", err)
	}
	var account model.Account
	var service_maps []model.ServiceMap
	account_id := u.Organization.ArcgisAccountID()
	if account_id != "" {
		account, err = queryarcgis.AccountFromID(ctx, account_id)
		if err != nil {
			return nil, nhttp.NewError("Failed to get arcgis: %w", err)
		}
		service_maps, err = queryarcgis.ServiceMapsFromAccountID(ctx, account.ID)
		if err != nil {
			return nil, nhttp.NewError("Failed to get map services: %w", err)
		}
	}
	data := contentSettingIntegration{
		ArcGISAccount: &account,
		ArcGISOAuth:   oauth,
		ServiceMaps:   service_maps,
	}
	return html.NewResponse("sync/configuration/integration-arcgis.html", data), nil
}

type contentSettingPlaceholder struct{}

func getConfigurationPesticide(ctx context.Context, r *http.Request, user platform.User) (*html.Response[contentSettingPlaceholder], *nhttp.ErrorWithStatus) {
	content := contentSettingPlaceholder{}
	return html.NewResponse("sync/configuration/pesticide.html", content), nil
}
func getConfigurationPesticideAdd(ctx context.Context, r *http.Request, user platform.User) (*html.Response[contentSettingPlaceholder], *nhttp.ErrorWithStatus) {
	content := contentSettingPlaceholder{}
	return html.NewResponse("sync/configuration/pesticide-add.html", content), nil
}
func getConfigurationUserAdd(ctx context.Context, r *http.Request, user platform.User) (*html.Response[contentSettingPlaceholder], *nhttp.ErrorWithStatus) {
	content := contentSettingPlaceholder{}
	return html.NewResponse("sync/configuration/user-add.html", content), nil
}
func getConfigurationUserList(ctx context.Context, r *http.Request, user platform.User) (*html.Response[contentSettingPlaceholder], *nhttp.ErrorWithStatus) {
	content := contentSettingPlaceholder{}
	return html.NewResponse("sync/configuration/user-list.html", content), nil
}

type formArcgisConfiguration struct {
	MapService *string `schema:"map-service"`
}

func postConfigurationIntegrationArcgis(ctx context.Context, r *http.Request, u platform.User, f formArcgisConfiguration) (string, *nhttp.ErrorWithStatus) {
	if f.MapService != nil {
		_, err := psql.Update(
			um.Table("organization"),
			um.SetCol("arcgis_map_service_id").ToArg(f.MapService),
			um.Where(psql.Quote("id").EQ(psql.Arg(u.Organization.ID))),
		).Exec(ctx, db.PGInstance.BobDB)
		if err != nil {
			return "", nhttp.NewError("Failed to update map service config: %w", err)
		}
		log.Info().Str("map-service", *f.MapService).Int32("org-id", u.Organization.ID).Msg("changed map service")
	} else {
		log.Info().Msg("no map service")
	}
	return "/configuration/integration/arcgis", nil
}
