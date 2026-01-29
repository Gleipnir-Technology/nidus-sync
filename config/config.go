package config

import (
	"fmt"
	"net/url"
	"os"

	"github.com/nyaruka/phonenumbers"
)

var (
	Bind                       string
	ClientID                   string
	ClientSecret               string
	DomainRMO                  string
	DomainNidus                string
	DomainTegola               string
	Environment                string
	FilesDirectoryLogo         string
	FilesDirectoryPublic       string
	FilesDirectoryUser         string
	FieldseekerSchemaDirectory string
	ForwardEmailAPIToken       string
	ForwardEmailReportAddress  string
	ForwardEmailReportPassword string
	ForwardEmailReportUsername string
	MapboxToken                string
	PGDSN                      string
	PhoneNumberReport          phonenumbers.PhoneNumber
	PhoneNumberReportStr       string
	PhoneNumberSupport         phonenumbers.PhoneNumber
	PhoneNumberSupportStr      string
	SentryDSN                  string
	TextProvider               string
	TwilioAuthToken            string
	TwilioAccountSID           string
	TwilioMessagingServiceSID  string
	TwilioRCSSenderRMO         string
	VoipMSNumber               string
	VoipMSPassword             string
	VoipMSUsername             string
)

func IsProductionEnvironment() bool {
	return Environment == "PRODUCTION"
}

func makeURL(domain, path string, args ...string) string {
	to_add := make([]any, 0)
	for _, a := range args {
		to_add = append(to_add, url.QueryEscape(a))
	}
	pattern := "https://" + domain + path
	return fmt.Sprintf(pattern, to_add...)
}

func MakeURLNidus(path string, args ...string) string {
	return makeURL(DomainNidus, path, args...)
}
func MakeURLReport(path string, args ...string) string {
	return makeURL(DomainRMO, path, args...)
}
func MakeURLTegola(path string, args ...string) string {
	return makeURL(DomainTegola, path, args...)
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
	DomainNidus = os.Getenv("DOMAIN_NIDUS")
	if DomainNidus == "" {
		return fmt.Errorf("You must specify a non-empty DOMAIN_NIDUS")
	}
	DomainRMO = os.Getenv("DOMAIN_RMO")
	if DomainRMO == "" {
		return fmt.Errorf("You must specify a non-empty DOMAIN_RMO")
	}
	DomainTegola = os.Getenv("DOMAIN_TEGOLA")
	if DomainTegola == "" {
		return fmt.Errorf("You must specify a non-empty DOMAIN_TEGOLA")
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
	FilesDirectoryLogo = os.Getenv("FILES_DIRECTORY_LOGO")
	if FilesDirectoryLogo == "" {
		return fmt.Errorf("You must specify a non-empty FILES_DIRECTORY_LOGO")
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
	PhoneNumberReportStr = os.Getenv("PHONE_NUMBER_RMO")
	if PhoneNumberReportStr == "" {
		return fmt.Errorf("You must specify a non-empty PHONE_NUMBER_RMO")
	}
	p, err := phonenumbers.Parse(PhoneNumberReportStr, "US")
	if err != nil {
		return fmt.Errorf("Failed to parse '%s' as a valid phone number: %w", PhoneNumberReportStr, err)
	}
	PhoneNumberReport = *p

	PhoneNumberSupportStr = os.Getenv("PHONE_NUMBER_SUPPORT")
	if PhoneNumberSupportStr == "" {
		return fmt.Errorf("You must specify a non-empty PHONE_NUMBER_SUPPORT")
	}
	p, err = phonenumbers.Parse(PhoneNumberSupportStr, "US")
	if err != nil {
		return fmt.Errorf("Failed to parse '%s' as a valid phone number: %w", PhoneNumberSupportStr, err)
	}
	PhoneNumberSupport = *p

	SentryDSN = os.Getenv("SENTRY_DSN")
	if SentryDSN == "" {
		return fmt.Errorf("You must specify a non-empty SENTRY_DSN")
	}
	TextProvider = os.Getenv("TEXT_PROVIDER")
	switch TextProvider {
	case "":
		return fmt.Errorf("You must specify a non-empty TEXT_PROVIDER")
	case "twilio":
	case "voipms":
		break
	default:
		return fmt.Errorf("Unrecognized text provider '%s'", TextProvider)
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
	TwilioRCSSenderRMO = os.Getenv("TWILIO_RCS_SENDER_RMO")
	if TwilioRCSSenderRMO == "" {
		return fmt.Errorf("You must specify a non-empty TWILIO_RCS_SENDER_RMO")
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
	if VoipMSPassword == "" {
		return fmt.Errorf("You must specify a non-empty VOIPMS_USERNAME")
	}
	return nil
}

func ArcGISOauthRedirectURL() string {
	return MakeURLNidus("/arcgis/oauth/callback")
}
