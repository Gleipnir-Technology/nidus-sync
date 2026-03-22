package html

import (
	"strconv"

	"github.com/Gleipnir-Technology/nidus-sync/config"
)

type ContentURL struct {
	API                contentURLAPI
	Configuration      contentURLConfiguration
	OAuthRefreshArcGIS string
	RMO                contentURLRMO
	Root               string
	Route              string
	Sidebar            contentURLSidebar
	Tegola             string
	Upload             contentURLUpload
}

func NewContentURL() ContentURL {
	return ContentURL{
		API:                newContentURLAPI(),
		Configuration:      newContentURLConfiguration(),
		OAuthRefreshArcGIS: config.MakeURLNidus("/arcgis/oauth/begin"),
		RMO:                newContentURLRMO(),
		Root:               config.MakeURLNidus("/"),
		Route:              config.MakeURLNidus("/route"),
		Sidebar:            newContentURLSidebar(),
		Tegola:             config.MakeURLTegola("/"),
		Upload:             newContentURLUpload(),
	}
}

type contentURLAPI struct {
	Communication string
	Publicreport  contentURLAPIPublicreport
}

func newContentURLAPI() contentURLAPI {
	return contentURLAPI{
		Communication: config.MakeURLNidus("/api/communication"),
	}
}

type contentURLAPIPublicreport struct {
	Message string
}

func newContentURLAPIPublicreport() contentURLAPIPublicreport {
	return contentURLAPIPublicreport{
		Message: config.MakeURLNidus("/api/publicreport/message"),
	}
}

type contentURLConfiguration struct {
	ArcGIS       string
	Fieldseeker  string
	Integration  string
	Organization string
	Pesticide    string
	PesticideAdd string
	Root         string
	User         string
	Upload       string
	UserAdd      string
}

func newContentURLConfiguration() contentURLConfiguration {
	return contentURLConfiguration{
		ArcGIS:       config.MakeURLNidus("/configuration/integration/arcgis"),
		Fieldseeker:  config.MakeURLNidus("/configuration/integration/fieldseeker"),
		Integration:  config.MakeURLNidus("/configuration/integration"),
		Organization: config.MakeURLNidus("/configuration/organization"),
		Pesticide:    config.MakeURLNidus("/configuration/pesticide"),
		PesticideAdd: config.MakeURLNidus("/configuration/pesticide/add"),
		Root:         config.MakeURLNidus("/configuration"),
		User:         config.MakeURLNidus("/configuration/user"),
		Upload:       config.MakeURLNidus("/configuration/upload"),
		UserAdd:      config.MakeURLNidus("/configuration/user/add"),
	}
}

type contentURLRMO struct {
	Mailer contentURLRMOMailer
}

func newContentURLRMO() contentURLRMO {
	return contentURLRMO{
		Mailer: newContentURLRMOMailer(),
	}
}

type contentURLRMOMailer struct {
	AppointmentConfirmed urlWithParams
	Confirm              urlWithParams
	Contribute           urlWithParams
	Evidence             urlWithParams
	Root                 urlWithParams
	Schedule             urlWithParams
	Update               urlWithParams
}

func newContentURLRMOMailer() contentURLRMOMailer {
	return contentURLRMOMailer{
		AppointmentConfirmed: makeURLWithParams(config.MakeURLReport, "/mailer/%s/appointment-confirmed"),
		Confirm:              makeURLWithParams(config.MakeURLReport, "/mailer/%s/confirm"),
		Contribute:           makeURLWithParams(config.MakeURLReport, "/mailer/%s/contribute"),
		Evidence:             makeURLWithParams(config.MakeURLReport, "/mailer/%s/evidence"),
		Root:                 makeURLWithParams(config.MakeURLReport, "/mailer/%s"),
		Schedule:             makeURLWithParams(config.MakeURLReport, "/mailer/%s/schedule"),
		Update:               makeURLWithParams(config.MakeURLReport, "/mailer/%s/update"),
	}
}

type contentURLSidebar struct {
	Communication string
	Configuration string
	Intelligence  string
	Operations    string
	Planning      string
	Review        string
}

func newContentURLSidebar() contentURLSidebar {
	return contentURLSidebar{
		Communication: config.MakeURLNidus("/communication"),
		Configuration: config.MakeURLNidus("/configuration"),
		Intelligence:  config.MakeURLNidus("/intelligence"),
		Operations:    config.MakeURLNidus("/operations"),
		Planning:      config.MakeURLNidus("/planning"),
		Review:        config.MakeURLNidus("/review"),
	}
}

type urlForID = func(int) string
type urlWithParams = func(...string) string

type urlMaker func(path string, args ...string) string

func makeURLForID(maker urlMaker, pattern string) urlForID {
	return func(id int) string {
		params := []string{
			strconv.Itoa(id),
		}
		return maker(pattern, params...)
	}
}
func makeURLWithParams(maker urlMaker, pattern string, args ...string) urlWithParams {
	return func(args ...string) string {
		return maker(pattern, args...)
	}
}

type contentURLUpload struct {
	Commit        urlForID
	Discard       urlForID
	Pool          string
	PoolCustom    string
	PoolFlyover   string
	SamplePoolCSV string
}

func newContentURLUpload() contentURLUpload {
	return contentURLUpload{
		Commit:        makeURLForID(config.MakeURLNidus, "/configuration/upload/%s/commit"),
		Discard:       makeURLForID(config.MakeURLNidus, "/configuration/upload/%s/discard"),
		Pool:          config.MakeURLNidus("/configuration/upload/pool"),
		PoolFlyover:   config.MakeURLNidus("/configuration/upload/pool/flyover"),
		PoolCustom:    config.MakeURLNidus("/configuration/upload/pool/custom"),
		SamplePoolCSV: config.MakeURLNidus("/static/file/sample-pool.csv"),
	}
}
