package html

import (
	"strconv"

	"github.com/Gleipnir-Technology/nidus-sync/config"
)

type ContentURL struct {
	Configuration      contentURLConfiguration
	OAuthRefreshArcGIS string
	Root               string
	Route              string
	Sidebar            contentURLSidebar
	Tegola             string
	Upload             contentURLUpload
}

func NewContentURL() ContentURL {
	return ContentURL{
		Configuration:      newContentURLConfiguration(),
		OAuthRefreshArcGIS: config.MakeURLNidus("/arcgis/oauth/begin"),
		Root:               config.MakeURLNidus("/"),
		Route:              config.MakeURLNidus("/route"),
		Sidebar:            newContentURLSidebar(),
		Tegola:             config.MakeURLTegola("/"),
		Upload:             newContentURLUpload(),
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

type urlForID = func(int) string

func makeURLForID(pattern string) urlForID {
	return func(id int) string {
		return config.MakeURLNidus(pattern, strconv.Itoa(id))
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
		Commit:        makeURLForID("/configuration/upload/%s/commit"),
		Discard:       makeURLForID("/configuration/upload/%s/discard"),
		Pool:          config.MakeURLNidus("/configuration/upload/pool"),
		PoolFlyover:   config.MakeURLNidus("/configuration/upload/pool/flyover"),
		PoolCustom:    config.MakeURLNidus("/configuration/upload/pool/custom"),
		SamplePoolCSV: config.MakeURLNidus("/static/file/sample-pool.csv"),
	}
}
