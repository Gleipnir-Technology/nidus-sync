package config

import (
	"fmt"
	"net/url"
	"os"
	"strconv"
)

var Bind, ClientID, ClientSecret, Environment, FieldseekerSchemaDirectory, MapboxToken, PGDSN, URLReport, URLSync, UserFilesDirectory string

// Build the ArcGIS authorization URL with PKCE
func BuildArcGISAuthURL(clientID string) string {
	baseURL := "https://www.arcgis.com/sharing/rest/oauth2/authorize/"

	params := url.Values{}
	params.Add("client_id", clientID)
	params.Add("redirect_uri", RedirectURL())
	params.Add("response_type", "code")
	//params.Add("code_challenge", generateCodeChallenge(codeVerifier))
	//params.Add("code_challenge_method", "S256")

	// See https://developers.arcgis.com/rest/users-groups-and-items/token/
	// expiration is defined in minutes
	var expiration int
	if IsProductionEnvironment() {
		// 2 weeks is the maximum allowed
		expiration = 20160
	} else {
		expiration = 20
	}
	params.Add("expiration", strconv.Itoa(expiration))

	return baseURL + "?" + params.Encode()
}

func IsProductionEnvironment() bool {
	return Environment == "PRODUCTION"
}

func MakeURLSync(path string) string {
	return fmt.Sprintf("https://%s%s", URLSync, path)
}

func Parse() error {
	ClientID = os.Getenv("ARCGIS_CLIENT_ID")
	if ClientID == "" {
		return fmt.Errorf("You must specify a non-empty ARCGIS_CLIENT_ID")
	}
	ClientSecret = os.Getenv("ARCGIS_CLIENT_SECRET")
	if ClientSecret == "" {
		return fmt.Errorf("You must specify a non-empty ARCGIS_CLIENT_SECRET")
	}
	URLReport = os.Getenv("URL_REPORT")
	if URLReport == "" {
		return fmt.Errorf("You must specify a non-empty URL_REPORT")
	}
	URLSync = os.Getenv("URL_SYNC")
	if URLSync == "" {
		return fmt.Errorf("You must specify a non-empty URL_SYNC")
	}
	Bind = os.Getenv("BIND")
	if Bind == "" {
		Bind = ":9001"
	}
	Environment = os.Getenv("ENVIRONMENT")
	if Environment == "" {
		return fmt.Errorf("You must specify a non-empty ENVIRONMENT")
	}
	if !(Environment == "PRODUCTION" || Environment == "DEVELOPMENT") {
		return fmt.Errorf("ENVIRONMENT should be either DEVELOPMENT or PRODUCTION")
	}
	MapboxToken = os.Getenv("MAPBOX_TOKEN")
	if MapboxToken == "" {
		return fmt.Errorf("You must specify a non-empty MAPBOX_TOKEN")
	}
	PGDSN = os.Getenv("POSTGRES_DSN")
	if PGDSN == "" {
		return fmt.Errorf("You must specify a non-empty POSTGRES_DSN")
	}
	FieldseekerSchemaDirectory = os.Getenv("FIELDSEEKER_SCHEMA_DIRECTORY")
	if FieldseekerSchemaDirectory == "" {
		return fmt.Errorf("You must specify a non-empty FIELDSEEKER_SCHEMA_DIRECTORY")
	}
	UserFilesDirectory = os.Getenv("USER_FILES_DIRECTORY")
	if UserFilesDirectory == "" {
		return fmt.Errorf("You must specify a non-empty USER_FILES_DIRECTORY")
	}
	return nil
}

func RedirectURL() string {
	return MakeURLSync("/arcgis/oauth/callback")
}
