package sync

import (
	"github.com/Gleipnir-Technology/nidus-sync/config"
)

type ContentURL struct {
	PoolCSVUpload       string
	SamplePoolCSV       string
	Setting             string
	SettingIntegration  string
	SettingPesticide    string
	SettingPesticideAdd string
	SettingUser         string
	SettingUserAdd      string
	Tegola              string
}

func newContentURL() ContentURL {
	return ContentURL{
		PoolCSVUpload:       config.MakeURLNidus("/pool/upload"),
		SamplePoolCSV:       config.MakeURLNidus("/static/file/sample-pool.csv"),
		Setting:             config.MakeURLNidus("/setting"),
		SettingIntegration:  config.MakeURLNidus("/setting/integration"),
		SettingPesticide:    config.MakeURLNidus("/setting/pesticide"),
		SettingPesticideAdd: config.MakeURLNidus("/setting/pesticide/add"),
		SettingUser:         config.MakeURLNidus("/setting/user"),
		SettingUserAdd:      config.MakeURLNidus("/setting/user/add"),
		Tegola:              config.MakeURLTegola("/"),
	}
}
