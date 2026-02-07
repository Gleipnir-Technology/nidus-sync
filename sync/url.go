package sync

import (
	"github.com/Gleipnir-Technology/nidus-sync/config"
)

type ContentURL struct {
	PoolCSVUpload string
	SamplePoolCSV string
}

func newContentURL() ContentURL {
	return ContentURL{
		PoolCSVUpload: config.MakeURLNidus("/pool/upload"),
		SamplePoolCSV: config.MakeURLNidus("/static/file/sample-pool.csv"),
	}
}
