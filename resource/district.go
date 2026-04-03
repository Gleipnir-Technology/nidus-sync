package resource

import (
	"context"
	nhttp "github.com/Gleipnir-Technology/nidus-sync/http"
	"github.com/Gleipnir-Technology/nidus-sync/platform"
	"net/http"
	//"github.com/rs/zerolog/log"
)

type districtR struct {
	router *router
}

type district struct {
	Name    string `json:"name"`
	Slug    string `json:"slug"`
	URLLogo string `json:"url_logo"`
}

func District(r *router) *districtR {
	return &districtR{
		router: r,
	}
}

func (res *districtR) List(ctx context.Context, r *http.Request, query QueryParams) ([]*district, *nhttp.ErrorWithStatus) {
	organizations, err := platform.OrganizationList(ctx)
	if err != nil {
		return nil, nhttp.NewError("list orgs: %w", err)
	}
	districts := make([]*district, 0)
	for _, org := range organizations {
		slug := org.Slug()
		if slug == "" {
			continue
		}
		logo, err := res.router.SlugToURI("district.logo.BySlug", slug)
		if err != nil {
			return nil, nhttp.NewError("logo url: %w", err)
		}
		districts = append(districts, &district{
			Name:    org.Name(),
			Slug:    slug,
			URLLogo: logo,
		})
	}
	return districts, nil
}
