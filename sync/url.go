package sync

import (
	"github.com/Gleipnir-Technology/nidus-sync/config"
)

type ContentURL struct {
	PoolCSVUpload string
}

func newContentURL() ContentURL {
	return ContentURL{
		PoolCSVUpload: config.MakeURLNidus("/pool/upload"),
	}
}
