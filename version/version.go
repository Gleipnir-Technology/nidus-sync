package version

import (
	"runtime/debug"
	"time"
)

type VersionInfo struct {
	BuildTime  time.Time `json:"build_time"`
	IsModified bool      `json:"is_modified"`
	Revision   string    `json:"revision"`
}

func Get() VersionInfo {
	info, ok := debug.ReadBuildInfo()
	if !ok {
		return VersionInfo{
			BuildTime:  time.Now(),
			IsModified: false,
			Revision:   "unknown",
		}
	}

	var version VersionInfo
	for _, setting := range info.Settings {
		switch setting.Key {
		case "vcs.modified":
			version.IsModified = setting.Value == "true"
		case "vcs.revision":
			if len(setting.Value) > 7 {
				version.Revision = setting.Value[:7]
			} else {
				version.Revision = setting.Value
			}
		case "vcs.time":
			if t, err := time.Parse(time.RFC3339, setting.Value); err == nil {
				version.BuildTime = t
			}
		}
	}

	return version
}
