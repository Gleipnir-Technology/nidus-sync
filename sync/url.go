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
		Setting:            newContentURLSetting(),
		Sidebar:            newContentURLSidebar(),
		Tegola:             config.MakeURLTegola("/"),
		UploadCSVPool:      config.MakeURLNidus("/upload/pool"),
	}
}

type contentURLConfiguration struct {
	Upload string
}

func newContentURLConfiguration() contentURLConfiguration {
	return contentURLConfiguration{
		Upload: config.MakeURLNidus("/configuration/upload"),
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

type contentURLSetting struct {
	ArcGIS       string
	Fieldseeker  string
	Integration  string
	Organization string
	Pesticide    string
	PesticideAdd string
	Root         string
	User         string
	UserAdd      string
}

func newContentURLSetting() contentURLSetting {
	return contentURLSetting{
		ArcGIS:       config.MakeURLNidus("/setting/integration/arcgis"),
		Fieldseeker:  config.MakeURLNidus("/setting/integration/fieldseeker"),
		Integration:  config.MakeURLNidus("/setting/integration"),
		Organization: config.MakeURLNidus("/setting/organization"),
		Pesticide:    config.MakeURLNidus("/setting/pesticide"),
		PesticideAdd: config.MakeURLNidus("/setting/pesticide/add"),
		Root:         config.MakeURLNidus("/setting"),
		User:         config.MakeURLNidus("/setting/user"),
		UserAdd:      config.MakeURLNidus("/setting/user/add"),
	}
}
