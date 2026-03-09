package api

import (
	"context"
	"net/http"
	"time"

	"github.com/Gleipnir-Technology/nidus-sync/config"
	"github.com/Gleipnir-Technology/nidus-sync/db/models"
	nhttp "github.com/Gleipnir-Technology/nidus-sync/http"
	"github.com/Gleipnir-Technology/nidus-sync/platform/publicreport"
	"github.com/Gleipnir-Technology/nidus-sync/platform/types"
	"github.com/google/uuid"
	//"github.com/rs/zerolog/log"
)

type reporter struct {
	HasEmail bool   `json:"has_email"`
	HasPhone bool   `json:"has_phone"`
	Name     string `json:"name"`
}
type communication struct {
	Created      time.Time          `json:"created"`
	ID           string             `json:"id"`
	PublicReport types.PublicReport `json:"public_report"`
	Type         string             `json:"type"`
}
type contentListCommunication struct {
	Communications []communication `json:"communications"`
}

func listCommunication(ctx context.Context, r *http.Request, org *models.Organization, user *models.User, query queryParams) (*contentListCommunication, *nhttp.ErrorWithStatus) {
	nreports, err := publicreport.NuisanceReportForOrganization(ctx, org.ID)
	if err != nil {
		return nil, nhttp.NewError("nuisance report query: %w", err)
	}
	wreports, err := publicreport.WaterReportForOrganization(ctx, org.ID)
	if err != nil {
		return nil, nhttp.NewError("water report query: %w", err)
	}
	comms := make([]communication, len(nreports)+len(wreports))
	for i, report := range nreports {
		comms[i] = communication{
			Created:      report.Created,
			ID:           report.PublicID,
			PublicReport: report,
			Type:         "nuisance",
		}
	}
	for i, report := range wreports {
		comms[i+len(nreports)] = communication{
			Created:      report.Created,
			ID:           report.PublicID,
			PublicReport: report,
			Type:         "water",
		}
	}
	return &contentListCommunication{
		Communications: comms,
	}, nil
}

func toImageURLs(m map[string][]uuid.UUID, id string) []string {
	uuids, ok := m[id]
	if !ok {
		return []string{}
	}
	urls := make([]string, len(uuids))
	for i, u := range uuids {
		urls[i] = config.MakeURLNidus("/api/image/%s/content", u.String())
	}
	return urls
}
