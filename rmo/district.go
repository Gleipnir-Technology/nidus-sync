package rmo

import (
	"github.com/Gleipnir-Technology/nidus-sync/config"
	"github.com/Gleipnir-Technology/nidus-sync/db/models"
)

type ContentDistrict struct {
	Name       string
	URLLogo    string
	URLWebsite string
}

func newContentDistrict(d *models.Organization) *ContentDistrict {
	if d == nil {
		return nil
	}
	return &ContentDistrict{
		Name:       d.Name,
		URLLogo:    config.MakeURLNidus("/api/district/%s/logo", d.Slug.GetOr("unset")),
		URLWebsite: d.Website.GetOr(""),
	}
}
