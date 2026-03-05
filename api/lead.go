package api

import (
	"context"
	"net/http"

	"github.com/Gleipnir-Technology/nidus-sync/db/models"
	nhttp "github.com/Gleipnir-Technology/nidus-sync/http"
	"github.com/rs/zerolog/log"
)

type formLeads struct {
	SignalIDs []int `schema:"signal_ids"`
}
type createdLead struct {
	ID int `json:"id"`
}
type contentListLead struct {
	Leads []lead `json:"leads"`
}
type lead struct {
	ID int32 `json:"id"`
}

func listLead(ctx context.Context, r *http.Request, org *models.Organization, user *models.User) (*contentListLead, *nhttp.ErrorWithStatus) {
	return &contentListLead{
		Leads: make([]lead, 0),
	}, nil
}
func postLeads(ctx context.Context, r *http.Request, org *models.Organization, user *models.User, f formLeads) (*createdLead, *nhttp.ErrorWithStatus) {
	log.Info().Ints("signal ids", f.SignalIDs).Msg("fake post leads")
	return &createdLead{
		ID: 0,
	}, nil
}
