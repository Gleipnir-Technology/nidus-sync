package static

import (
	"embed"
)

//go:embed css gen file ico js vendor
var EmbeddedStaticFS embed.FS
