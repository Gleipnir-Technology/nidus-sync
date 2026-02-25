package main

import (
	"context"
	"flag"
	//"fmt"
	//"net/url"
	"os"

	"github.com/Gleipnir-Technology/arcgis-go"
	"github.com/Gleipnir-Technology/arcgis-go/fieldseeker"
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

	ctx := context.TODO()
	custom_logger := log.With().Logger().Level(zerolog.DebugLevel)
	ctx = arcgis.WithLogger(ctx, custom_logger)
	err = db.InitializeDatabase(ctx, config.PGDSN)
	if err != nil {
		log.Error().Err(err).Msg("Failed to connect to database")
		os.Exit(3)
	}

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

	fieldseeker, err := background.NewFieldSeeker(ctx, oauth)
	if err != nil {
		log.Error().Err(err).Msg("Failed to create fieldseeker ")
		os.Exit(6)
	}

	//ar.GeocodeFindAddressCandidates("1 Infinite Loop, Cupertino, CA")
	info, err := fieldseeker.Arcgis.Info(ctx)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get info")
		os.Exit(7)
	}
	log.Info().Float64("current_version", info.CurrentVersion).Str("full_version", info.FullVersion).Str("owning_system_url", info.OwningSystemUrl).Str("owning_tenant", info.OwningTenant).Msg("Got info")

	err = fieldseeker.Arcgis.SwitchHostByPortal(ctx)
	search_resp, err := fieldseeker.Arcgis.Search(ctx, "type:Map Service")
	if err != nil || search_resp == nil {
		log.Error().Err(err).Msg("Failed to make search")
		os.Exit(9)
	}
	for _, item := range search_resp.Results {
		log.Info().Str("name", item.Name).Str("id", item.ID).Str("url", item.URL).Msg("Found a search result")
	}
	//for _, portal := range portals. {
	//}
	/*
		u, err := url.Parse("https://tiles.arcgis.com/tiles/pV7SH1EgRc6tpxlJ/arcgis/rest/services/TrimmedFlyover2025/MapServer?f=json")
		if err != nil || u == nil {
			log.Error().Err(err).Msg("Failed to make url")
			os.Exit(8)
		}
		body, err := fieldseeker.Arcgis.RawGet(ctx, *u)
		if err != nil {
			log.Error().Err(err).Msg("Failed to raw get url")
			os.Exit(9)
		}
		log.Info().Str("body", string(body)).Msg("Got it")
	*/
}

func printServices(ctx context.Context, fs fieldseeker.FieldSeeker) {
	//shows services, but there's no TrimmedFlyover2025 in them
	services, err := fs.Arcgis.Services(ctx)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get services")
		os.Exit(10)
	}
	for _, service := range services.Services {
		fs, err := fs.Arcgis.GetFeatureServer(ctx, service.Name)
		if err != nil {
			log.Error().Err(err).Msg("Failed to get FS")
			os.Exit(11)
		}
		log.Info().Str("name", service.Name).Str("type", service.Type).Str("url", service.URL).Msg("Found service")
		for _, l := range fs.Layers {
			log.Info().Str("name", l.Name).Str("type", l.Type).Msg("Layer")
		}
	}
}
