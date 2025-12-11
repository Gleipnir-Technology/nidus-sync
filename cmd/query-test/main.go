package main

import (
	"context"
	"os"
	"time"

	fslayer "github.com/Gleipnir-Technology/arcgis-go/fieldseeker/layer"
	"github.com/Gleipnir-Technology/nidus-sync/db"
	"github.com/alexedwards/scs/v2"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var sessionManager *scs.SessionManager

var BaseURL, ClientID, ClientSecret, Environment, FieldseekerSchemaDirectory, MapboxToken string

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	ClientID = os.Getenv("ARCGIS_CLIENT_ID")
	if ClientID == "" {
		log.Error().Msg("You must specify a non-empty ARCGIS_CLIENT_ID")
		os.Exit(1)
	}
	ClientSecret = os.Getenv("ARCGIS_CLIENT_SECRET")
	if ClientSecret == "" {
		log.Error().Msg("You must specify a non-empty ARCGIS_CLIENT_SECRET")
		os.Exit(1)
	}
	BaseURL = os.Getenv("BASE_URL")
	if BaseURL == "" {
		log.Error().Msg("You must specify a non-empty BASE_URL")
		os.Exit(1)
	}
	bind := os.Getenv("BIND")
	if bind == "" {
		bind = ":9001"
	}
	Environment = os.Getenv("ENVIRONMENT")
	if Environment == "" {
		log.Error().Msg("You must specify a non-empty ENVIRONMENT")
		os.Exit(1)
	}
	if !(Environment == "PRODUCTION" || Environment == "DEVELOPMENT") {
		log.Error().Str("ENVIRONMENT", Environment).Msg("ENVIRONMENT should be either DEVELOPMENT or PRODUCTION")
		os.Exit(2)
	}
	MapboxToken = os.Getenv("MAPBOX_TOKEN")
	if MapboxToken == "" {
		log.Error().Msg("You must specify a non-empty MAPBOX_TOKEN")
		os.Exit(1)
	}
	pg_dsn := os.Getenv("POSTGRES_DSN")
	if pg_dsn == "" {
		log.Error().Msg("You must specify a non-empty POSTGRES_DSN")
		os.Exit(1)
	}
	FieldseekerSchemaDirectory = os.Getenv("FIELDSEEKER_SCHEMA_DIRECTORY")
	if FieldseekerSchemaDirectory == "" {
		log.Error().Msg("You must specify a non-empty FIELDSEEKER_SCHEMA_DIRECTORY")
		os.Exit(1)
	}

	log.Info().Msg("Starting...")
	err := db.InitializeDatabase(context.TODO(), pg_dsn)
	if err != nil {
		log.Error().Str("err", err.Error()).Msg("Failed to connect to database")
		os.Exit(2)
	}
	ctx := context.Background()
	row := fslayer.RodentLocation{
		ObjectID:     1,
		LocationName: "some location",
		Zone:         "",
		Zone2:        "",
		//Habitat:                      fslayer.RodentLocationRodentLocationHabitatCommercial,
		//Priority:                     fslayer.RodentLocationLocationPriority1None,
		//Usetype:                      fslayer.RodentLocationLocationUseType1Residential,
		//Active:                       fslayer.RodentLocationNotInUITF1True,
		Description: "",
		Accessdesc:  "",
		Comments:    "",
		//Symbology:                    fslayer.RodentLocationRodentLocationSymbologyActionrequired,
		ExternalID:                   "",
		Nextactiondatescheduled:      time.Now(),
		Locationnumber:               1,
		LastInspectionDate:           time.Now(),
		LastInspectionSpecies:        "",
		LastInspectionAction:         "",
		LastInspectionConditions:     "",
		LastInspectionRodentEvidence: "",
		GlobalID:                     uuid.New(),
		CreatedUser:                  "",
		CreatedDate:                  time.Now(),
		LastEditedUser:               "",
		LastEditedDate:               time.Now(),
		CreationDate:                 time.Now(),
		Creator:                      "",
		EditDate:                     time.Now(),
		Editor:                       "",
		Jurisdiction:                 "",
	}
	err = db.TestPreparedQuery(ctx, &row)
	if err != nil {
		log.Error().Str("err", err.Error()).Msg("Failed to run prepared query")
		os.Exit(3)
	}
	log.Info().Msg("Complete.")
}
