package sync

import (
	"github.com/Gleipnir-Technology/nidus-sync/config"
)

type contentURL struct {
	OAuthRefreshArcGIS string
	Root               string
	Route              string
	SamplePoolCSV      string
	Setting            contentURLSetting
	Tegola             string
	UploadCSVPool      string
}

type contentURLSetting struct {
	Fieldseeker  string
	Integration  string
	Organization string
	Pesticide    string
	PesticideAdd string
	Root         string
	User         string
	UserAdd      string
}

func newContentURL() contentURL {
	return contentURL{
		OAuthRefreshArcGIS: config.MakeURLNidus("/arcgis/oauth/begin"),
		Root:               config.MakeURLNidus("/"),
		Route:              config.MakeURLNidus("/route"),
		SamplePoolCSV:      config.MakeURLNidus("/static/file/sample-pool.csv"),
		Setting:            newContentURLSetting(),
		Tegola:             config.MakeURLTegola("/"),
		UploadCSVPool:      config.MakeURLNidus("/upload/pool"),
	}
}
func newContentURLSetting() contentURLSetting {
	return contentURLSetting{
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
