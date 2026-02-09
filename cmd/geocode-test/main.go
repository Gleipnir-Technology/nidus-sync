package main

import (
	"context"
	"flag"
	//"fmt"
	"os"

	"github.com/Gleipnir-Technology/arcgis-go"
	//fslayer "github.com/Gleipnir-Technology/arcgis-go/fieldseeker/layer"
	"github.com/Gleipnir-Technology/nidus-sync/background"
	"github.com/Gleipnir-Technology/nidus-sync/config"
	"github.com/Gleipnir-Technology/nidus-sync/db"
	"github.com/Gleipnir-Technology/nidus-sync/db/models"
	"github.com/alexedwards/scs/v2"
	//"github.com/google/uuid"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var sessionManager *scs.SessionManager

var BaseURL, ClientID, ClientSecret, Environment, FieldseekerSchemaDirectory, MapboxToken string

func main() {
	org_id := flag.Int("org", 0, "The ID of the organization to use")
	flag.Parse()
	if org_id == nil || *org_id == 0 {
		log.Error().Msg("You must specify -org_id")
		os.Exit(1)
	}
	err := config.Parse()
	if err != nil {
		log.Error().Err(err).Msg("Failed to parse config")
		os.Exit(2)
	}
	log.Info().Msg("Starting...")
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	err = db.InitializeDatabase(context.TODO(), config.PGDSN)
	if err != nil {
		log.Error().Err(err).Msg("Failed to connect to database")
		os.Exit(3)
	}

	ctx := context.TODO()
	org, err := models.FindOrganization(ctx, db.PGInstance.BobDB, int32(*org_id))
	if err != nil {
		log.Error().Err(err).Msg("Failed to get org")
		os.Exit(4)
	}

	oauth, err := background.GetOAuthForOrg(ctx, org)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get oauth for org")
		os.Exit(5)
	}
	ar := arcgis.NewArcGIS(
		arcgis.AuthenticatorOAuth{
			AccessToken:         oauth.AccessToken,
			AccessTokenExpires:  oauth.AccessTokenExpires,
			RefreshToken:        oauth.RefreshToken,
			RefreshTokenExpires: oauth.RefreshTokenExpires,
		},
	)
	ar.GeocodeFindAddressCandidates("1 Infinite Loop, Cupertino, CA")
}
