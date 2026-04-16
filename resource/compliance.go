package resource

import (
	"context"
	nhttp "github.com/Gleipnir-Technology/nidus-sync/http"
	"github.com/Gleipnir-Technology/nidus-sync/platform"
	"net/http"
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
