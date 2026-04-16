package resource

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Gleipnir-Technology/nidus-sync/db"
	"github.com/Gleipnir-Technology/nidus-sync/db/models"
	nhttp "github.com/Gleipnir-Technology/nidus-sync/http"
	"github.com/Gleipnir-Technology/nidus-sync/platform"
	"github.com/Gleipnir-Technology/nidus-sync/platform/types"
	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"
)

func ComplianceRequest(r *router) *complianceRequestR {
	return &complianceRequestR{
		router: r,
	}
}

type complianceRequestR struct {
	router *router
}
type complianceRequestMailer struct {
	ID int32 `json:"id"`
}
type complianceRequestMailerForm struct {
	SiteID int32 `json:"site_id"`
}

func (res *complianceRequestR) CreateMailer(ctx context.Context, r *http.Request, user platform.User, n complianceRequestMailerForm) (*complianceRequestMailer, *nhttp.ErrorWithStatus) {
	id, err := platform.ComplianceRequestMailerCreate(ctx, user, n.SiteID)
	if err != nil {
		return nil, nhttp.NewError("create mailer: %w", err)
	}
	return &complianceRequestMailer{
		ID: id,
	}, nil
}
func (res *complianceRequestR) ImagePoolGet(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	public_id := vars["public_id"]
	if public_id == "" {
		http.Error(w, "need public ID", http.StatusNotFound)
		return
	}

	ctx := r.Context()
	err := imagePoolGet(ctx, w, public_id)
	if err != nil {
		log.Error().Err(err).Msg("failed to get image")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func imagePoolGet(ctx context.Context, w http.ResponseWriter, public_id string) error {
	txn := db.PGInstance.BobDB
	compliance_req, err := models.ComplianceReportRequests.Query(
		models.SelectWhere.ComplianceReportRequests.PublicID.EQ(public_id),
	).One(ctx, txn)
	if err != nil {
		return fmt.Errorf("find compliance report: %w", err)
	}

	if compliance_req.LeadID.IsNull() {
		return fmt.Errorf("no lead for compliance req %d", compliance_req.ID)
	}
	lead_id := compliance_req.LeadID.MustGet()
	lead, err := models.FindLead(ctx, txn, lead_id)
	if err != nil {
		return fmt.Errorf("find lead: %w", err)
	}

	if lead.SiteID.IsNull() {
		return fmt.Errorf("no site for lead %d", lead.ID)
	}
	site_id := lead.SiteID.MustGet()
	site, err := models.FindSite(ctx, txn, site_id)
	if err != nil {
		return fmt.Errorf("find site: %w", err)
	}
	organization, err := models.FindOrganization(ctx, txn, site.OrganizationID)
	if err != nil {
		return fmt.Errorf("find address: %w", err)
	}
	features, err := platform.FeaturesForSite(ctx, site.ID)
	if err != nil {
		return fmt.Errorf("get features: %w", err)
	}
	log.Debug().Int("len", len(features)).Int32("site", site.ID).Msg("got features for site")
	var pool_feature *types.Feature
	for _, f := range features {
		if f.Type == "pool" {
			pool_feature = &f
		}
	}
	if pool_feature == nil {
		return fmt.Errorf("no pool feature")
	}

	level := uint(15)
	err = platform.WriteTile(ctx, w, organization, level, pool_feature.Location.Latitude, pool_feature.Location.Longitude)
	if err != nil {
		return fmt.Errorf("write tile at %d, %f %f: %w", level, pool_feature.Location.Longitude, pool_feature.Location.Latitude, err)
	}
	w.WriteHeader(http.StatusOK)
	return nil
}
