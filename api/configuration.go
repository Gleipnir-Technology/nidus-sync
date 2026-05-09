package api

import (
	"context"
	"net/http"

	"github.com/Gleipnir-Technology/bob/dialect/psql"
	"github.com/Gleipnir-Technology/bob/dialect/psql/um"
	"github.com/Gleipnir-Technology/nidus-sync/db"
	nhttp "github.com/Gleipnir-Technology/nidus-sync/http"
	"github.com/Gleipnir-Technology/nidus-sync/platform"
	"github.com/rs/zerolog/log"
)

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
