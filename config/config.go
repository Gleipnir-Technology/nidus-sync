package config

import (
	"fmt"
	"net/url"
	"os"
	"strconv"

	"github.com/nyaruka/phonenumbers"
)

var (
	Bind                       string
	ClientID                   string
	ClientSecret               string
	Environment                string
	FilesDirectoryPublic       string
	FilesDirectoryUser         string
	FieldseekerSchemaDirectory string
	ForwardEmailAPIToken       string
	ForwardEmailReportAddress  string
	ForwardEmailReportPassword string
	ForwardEmailReportUsername string
	MapboxToken                string
	PGDSN                      string
	RMODomain                  string
	RMOPhoneNumber             phonenumbers.PhoneNumber
	URLReport                  string
	URLSync                    string
	URLTegola                  string
	TwilioAuthToken            string
	TwilioAccountSID           string
	TwilioMessagingServiceSID  string
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

func Parse() (err error) {
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
	ForwardEmailReportAddress = os.Getenv("FORWARDEMAIL_REPORT_ADDRESS")
	if ForwardEmailReportAddress == "" {
		return fmt.Errorf("You must specify a non-empty FORWARDEMAIL_REPORT_ADDRESS")
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
	RMODomain = os.Getenv("RMO_DOMAIN")
	if RMODomain == "" {
		return fmt.Errorf("You must specify a non-empty RMO_DOMAIN")
	}
	rmo_phone_number := os.Getenv("RMO_PHONE_NUMBER")
	if rmo_phone_number == "" {
		return fmt.Errorf("You must specify a non-empty RMO_PHONE_NUMBER")
	}
	p, err := phonenumbers.Parse(rmo_phone_number, "US")
	if err != nil {
		return fmt.Errorf("Failed to parse '%s' as a valid phone number: %w", rmo_phone_number, err)
	}
	RMOPhoneNumber = *p

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
	TwilioAccountSID = os.Getenv("TWILIO_ACCOUNT_SID")
	if TwilioAccountSID == "" {
		return fmt.Errorf("You must specify a non-empty TWILIO_ACCOUNT_SID")
	}
	TwilioAuthToken = os.Getenv("TWILIO_AUTH_TOKEN")
	if TwilioAuthToken == "" {
		return fmt.Errorf("You must specify a non-empty TWILIO_AUTH_TOKEN")
	}
	TwilioMessagingServiceSID = os.Getenv("TWILIO_MESSAGING_SERVICE_SID")
	if TwilioMessagingServiceSID == "" {
		return fmt.Errorf("You must specify a non-empty TWILIO_MESSAGING_SERVICE_SID")
	}
	return nil
}

func RedirectURL() string {
	return MakeURLSync("/arcgis/oauth/callback")
}
