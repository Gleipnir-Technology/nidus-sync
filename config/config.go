package config

import (
	"fmt"
	"net/url"
	"os"
	"strconv"
)

var (
	Bind string
	ClientID string
	ClientSecret string
	Environment string
	FilesDirectoryPublic string
	FilesDirectoryUser string
	FieldseekerSchemaDirectory string
	ForwardEmailAPIToken string
	ForwardEmailReportPassword string
	ForwardEmailReportUsername string
	MapboxToken string
	PGDSN string
	URLReport string
	URLSync string
	URLTegola string
	VoipMSPassword string
	VoipMSNumber string
	VoipMSUsername string
)

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
	Bind = os.Getenv("BIND")
	if Bind == "" {
		Bind = ":9001"
	}
	ClientID = os.Getenv("ARCGIS_CLIENT_ID")
	if ClientID == "" {
		return fmt.Errorf("You must specify a non-empty ARCGIS_CLIENT_ID")
	}
	ClientSecret = os.Getenv("ARCGIS_CLIENT_SECRET")
	if ClientSecret == "" {
		return fmt.Errorf("You must specify a non-empty ARCGIS_CLIENT_SECRET")
	}
	Environment = os.Getenv("ENVIRONMENT")
	if Environment == "" {
		return fmt.Errorf("You must specify a non-empty ENVIRONMENT")
	}
	if !(Environment == "PRODUCTION" || Environment == "DEVELOPMENT") {
		return fmt.Errorf("ENVIRONMENT should be either DEVELOPMENT or PRODUCTION")
	}
	FieldseekerSchemaDirectory = os.Getenv("FIELDSEEKER_SCHEMA_DIRECTORY")
	if FieldseekerSchemaDirectory == "" {
		return fmt.Errorf("You must specify a non-empty FIELDSEEKER_SCHEMA_DIRECTORY")
	}
	FilesDirectoryPublic = os.Getenv("FILES_DIRECTORY_PUBLIC")
	if FilesDirectoryPublic == "" {
		return fmt.Errorf("You must specify a non-empty FILES_DIRECTORY_PUBLIC")
	}
	FilesDirectoryUser = os.Getenv("FILES_DIRECTORY_USER")
	if FilesDirectoryUser == "" {
		return fmt.Errorf("You must specify a non-empty FILES_DIRECTORY_USER")
	}
	ForwardEmailAPIToken = os.Getenv("FORWARDEMAIL_API_TOKEN")
	if ForwardEmailAPIToken == "" {
		return fmt.Errorf("You must specify a non-empty FORWARDEMAIL_API_TOKEN")
	}
	ForwardEmailReportUsername = os.Getenv("FORWARDEMAIL_REPORT_USERNAME")
	if ForwardEmailReportUsername == "" {
		return fmt.Errorf("You must specify a non-empty FORWARDEMAIL_REPORT_USERNAME")
	}
	ForwardEmailReportPassword = os.Getenv("FORWARDEMAIL_REPORT_PASSWORD")
	if ForwardEmailReportPassword == "" {
		return fmt.Errorf("You must specify a non-empty FORWARDEMAIL_REPORT_PASSWORD")
	}
	MapboxToken = os.Getenv("MAPBOX_TOKEN")
	if MapboxToken == "" {
		return fmt.Errorf("You must specify a non-empty MAPBOX_TOKEN")
	}
	PGDSN = os.Getenv("POSTGRES_DSN")
	if PGDSN == "" {
		return fmt.Errorf("You must specify a non-empty POSTGRES_DSN")
	}
	URLReport = os.Getenv("URL_REPORT")
	if URLReport == "" {
		return fmt.Errorf("You must specify a non-empty URL_REPORT")
	}
	URLSync = os.Getenv("URL_SYNC")
	if URLSync == "" {
		return fmt.Errorf("You must specify a non-empty URL_SYNC")
	}
	URLTegola = os.Getenv("URL_TEGOLA")
	if URLTegola == "" {
		return fmt.Errorf("You must specify a non-empty URL_TEGOLA")
	}
	VoipMSNumber = os.Getenv("VOIPMS_NUMBER")
	if VoipMSNumber == "" {
		return fmt.Errorf("You must specify a non-empty VOIPMS_NUMBER")
	}
	VoipMSPassword = os.Getenv("VOIPMS_PASSWORD")
	if VoipMSPassword == "" {
		return fmt.Errorf("You must specify a non-empty VOIPMS_PASSWORD")
	}
	VoipMSUsername = os.Getenv("VOIPMS_USERNAME")
	if VoipMSUsername == "" {
		return fmt.Errorf("You must specify a non-empty VOIPMS_USERNAME")
	}
	return nil
}

func RedirectURL() string {
	return MakeURLSync("/arcgis/oauth/callback")
}
