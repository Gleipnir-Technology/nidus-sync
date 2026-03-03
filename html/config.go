package html

import (
	"github.com/Gleipnir-Technology/nidus-sync/config"
)

type ContentConfig struct {
	IsProductionEnvironment bool
}

func NewContentConfig() ContentConfig {
	return ContentConfig{
		IsProductionEnvironment: config.IsProductionEnvironment(),
	}
}
