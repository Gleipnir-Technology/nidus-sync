package sync

import (
	"github.com/Gleipnir-Technology/nidus-sync/config"
)

type contentConfig struct {
	IsProductionEnvironment bool
}

func newContentConfig() contentConfig {
	return contentConfig{
		IsProductionEnvironment: config.IsProductionEnvironment(),
	}
}
