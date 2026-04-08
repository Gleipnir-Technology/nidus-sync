package resource

import (
	"context"
	"fmt"
	"strconv"

	nhttp "github.com/Gleipnir-Technology/nidus-sync/http"
	"github.com/Gleipnir-Technology/nidus-sync/platform"
	"github.com/gorilla/mux"
	"net/http"
	//"github.com/rs/zerolog/log"
)

type districtR struct {
	router *router
}

type district struct {
	Name        string `json:"name"`
	PhoneOffice string `json:"phone_office"`
	Slug        string `json:"slug"`
	URI         string `json:"uri"`
	URLLogo     string `json:"url_logo"`
	URLWebsite  string `json:"url_website"`
}

func District(r *router) *districtR {
	return &districtR{
		router: r,
	}
}

func (res *districtR) GetByID(ctx context.Context, r *http.Request, query QueryParams) (*district, *nhttp.ErrorWithStatus) {
	vars := mux.Vars(r)
	id_str := vars["id"]
	id, err := strconv.Atoi(id_str)
	if err != nil {
		return nil, nhttp.NewBadRequest("id conversion: %w", err)
	}
	org, err := platform.OrganizationByID(ctx, id)
	if err != nil {
		return nil, nhttp.NewError("get org: %w", err)
	}
	district, err := newDistrict(res.router, org)
	if err != nil {
		return nil, nhttp.NewError("new district: %w", err)
	}
	return district, nil
}
func (res *districtR) List(ctx context.Context, r *http.Request, query QueryParams) ([]*district, *nhttp.ErrorWithStatus) {
	organizations, err := platform.OrganizationList(ctx)
	if err != nil {
		return nil, nhttp.NewError("list orgs: %w", err)
	}
	districts := make([]*district, 0)
	for _, org := range organizations {
		district, err := newDistrict(res.router, org)
		if err != nil {
			return nil, nhttp.NewError("make district: %w", err)
		}
		if district == nil {
			continue
		}
		districts = append(districts, district)
	}
	return districts, nil
}

func newDistrict(r *router, org *platform.Organization) (*district, error) {
	slug := org.Slug()
	if slug == "" {
		return nil, nil
	}
	logo, err := r.SlugToURI("district.logo.BySlug", slug)
	if err != nil {
		return nil, fmt.Errorf("logo url: %w", err)
	}
	uri, err := r.IDToURI("district.ByIDGet", int(org.ID))
	if err != nil {
		return nil, nhttp.NewError("district uri: %w", err)
	}
	return &district{
		Name:        org.Name(),
		PhoneOffice: org.PhoneOffice(),
		Slug:        slug,
		URI:         uri,
		URLLogo:     logo,
		URLWebsite:  org.Website(),
	}, nil
}
