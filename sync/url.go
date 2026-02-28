package sync

import (
	"github.com/Gleipnir-Technology/nidus-sync/config"
)

type contentURL struct {
	Configuration      contentURLConfiguration
	OAuthRefreshArcGIS string
	Root               string
	Route              string
	SamplePoolCSV      string
	Sidebar            contentURLSidebar
	Tegola             string
	UploadCSVPool      string
}

func newContentURL() contentURL {
	return contentURL{
		Configuration:      newContentURLConfiguration(),
		OAuthRefreshArcGIS: config.MakeURLNidus("/arcgis/oauth/begin"),
		Root:               config.MakeURLNidus("/"),
		Route:              config.MakeURLNidus("/route"),
		SamplePoolCSV:      config.MakeURLNidus("/static/file/sample-pool.csv"),
		Sidebar:            newContentURLSidebar(),
		Tegola:             config.MakeURLTegola("/"),
		UploadCSVPool:      config.MakeURLNidus("/configuration/upload/pool"),
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
