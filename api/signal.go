package api

import (
	"context"
	"net/http"
	"time"

	"github.com/Gleipnir-Technology/bob/dialect/psql/sm"
	"github.com/Gleipnir-Technology/nidus-sync/db"
	"github.com/Gleipnir-Technology/nidus-sync/db/models"
	nhttp "github.com/Gleipnir-Technology/nidus-sync/http"
	"github.com/Gleipnir-Technology/nidus-sync/platform"
	"github.com/aarondl/opt/null"
)

type signal struct {
	Addressed *time.Time     `json:"addressed"`
	Addressor *platform.User `json:"addressed"`
	Created   time.Time      `json:"created"`
	Creator   platform.User  `json:"creator"`
	ID        int32          `json:"id"`
	Species   string         `json:"species"`
	Type      string         `json:"type"`
}
type contentListSignal struct {
	Signals []signal `json:"signals"`
}

func listSignal(ctx context.Context, r *http.Request, org *models.Organization, user *models.User) (*contentListSignal, *nhttp.ErrorWithStatus) {
	rows, err := models.Signals.Query(
		models.SelectWhere.Signals.OrganizationID.EQ(org.ID),
		sm.OrderBy("created").Desc(),
	).All(ctx, db.PGInstance.BobDB)
	if err != nil {
		return nil, nhttp.NewError("failed to get signals: %w", err)
	}
	users_by_id, err := platform.UsersByID(ctx, org)
	if err != nil {
		return nil, nhttp.NewError("users by id: %w", err)
	}
	signals := make([]signal, len(rows))
	for i, row := range rows {
		var species string = ""
		if row.Species.IsValue() {
			species = row.Species.MustGet().String()
		}
		signals[i] = signal{
			Addressed: row.Addressed.Ptr(),
			Addressor: userOrNil(users_by_id, row.Addressor),
			Created:   row.Created,
			Creator:   *users_by_id[row.Creator],
			ID:        row.ID,
			Species:   species,
			Type:      row.Type.String(),
		}
	}
	return &contentListSignal{
		Signals: signals,
	}, nil
}

func userOrNil(usersByID map[int32]*platform.User, id null.Val[int32]) *platform.User {
	if id.IsNull() {
		return nil
	}
	u, ok := usersByID[id.MustGet()]
	if !ok {
		return nil
	}
	return u
}
