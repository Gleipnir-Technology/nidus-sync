package api

import (
	"context"
	"net/http"
	"time"

	"github.com/Gleipnir-Technology/bob"
	"github.com/Gleipnir-Technology/bob/dialect/psql"
	"github.com/Gleipnir-Technology/bob/dialect/psql/sm"
	"github.com/Gleipnir-Technology/bob/dialect/psql/um"
	"github.com/Gleipnir-Technology/nidus-sync/db"
	"github.com/Gleipnir-Technology/nidus-sync/db/enums"
	"github.com/Gleipnir-Technology/nidus-sync/db/models"
	nhttp "github.com/Gleipnir-Technology/nidus-sync/http"
	"github.com/Gleipnir-Technology/nidus-sync/platform"
	"github.com/Gleipnir-Technology/nidus-sync/platform/geom"
	"github.com/aarondl/opt/omit"
	"github.com/aarondl/opt/omitnull"
	"github.com/rs/zerolog/log"
	"github.com/stephenafamo/scan"
)

type createLead struct {
	PoolLocations map[int]Location `json:"pool_locations"`
	SignalIDs     []int            `json:"signal_ids"`
}
type createdLead struct {
	ID int32 `json:"id"`
}
type contentListLead struct {
	Leads []lead `json:"leads"`
}
type lead struct {
	ID int32 `json:"id"`
}

func listLead(ctx context.Context, r *http.Request, user platform.User, query queryParams) (*contentListLead, *nhttp.ErrorWithStatus) {
	return &contentListLead{
		Leads: make([]lead, 0),
	}, nil
}
func postLeads(ctx context.Context, r *http.Request, user platform.User, req createLead) (*createdLead, *nhttp.ErrorWithStatus) {
	if len(req.SignalIDs) == 0 {
		return nil, nhttp.NewErrorStatus(http.StatusBadRequest, "can't make a lead with no signals")
	}
	if len(req.SignalIDs) > 1 {
		return nil, nhttp.NewErrorStatus(http.StatusBadRequest, "can't make a lead with multiple signals yet")
	}
	signal_id := req.SignalIDs[0]
	txn, err := db.PGInstance.BobDB.BeginTx(ctx, nil)
	defer txn.Rollback(ctx)

	if err != nil {
		return nil, nhttp.NewError("start transaction: %w", err)
	}
	type _Row struct {
		ID int32 `db:"site_id"`
	}
	site, err := bob.One(ctx, db.PGInstance.BobDB, psql.Select(
		sm.Columns(
			"pool.site_id AS site_id",
		),
		sm.From("signal_pool"),
		sm.InnerJoin("pool").OnEQ(
			psql.Quote("signal_pool", "pool_id"),
			psql.Quote("pool", "id"),
		),
		sm.InnerJoin("site").On(
			psql.Quote("pool", "site_id").EQ(psql.Quote("site", "id")),
		),
		sm.Where(psql.Quote("signal_pool", "signal_id").EQ(psql.Arg(signal_id))),
		sm.Where(psql.Quote("site", "organization_id").EQ(psql.Arg(user.Organization.ID()))),
	), scan.StructMapper[_Row]())
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			return nil, nhttp.NewErrorStatus(http.StatusBadRequest, "Can't make a lead from signal %d: %w", signal_id, err)
		}
		return nil, nhttp.NewError("failed getting site: %w", err)
	}

	lead, err := models.Leads.Insert(&models.LeadSetter{
		Created: omit.From(time.Now()),
		Creator: omit.From(int32(user.ID)),
		// ID
		OrganizationID: omit.From(int32(user.Organization.ID())),
		SiteID:         omitnull.From(site.ID),
		Type:           omit.From(enums.LeadtypeGreenPool),
	}).One(ctx, txn)
	if err != nil {
		return nil, nhttp.NewError("failed to create lead: %w", err)
	}
	_, err = psql.Update(
		um.Table("signal"),
		um.SetCol("addressed").ToArg(time.Now()),
		um.SetCol("addressor").ToArg(user.ID),
		um.Where(psql.Quote("id").EQ(psql.Arg(signal_id))),
	).Exec(ctx, txn)
	if err != nil {
		return nil, nhttp.NewError("failed to update signal %d: %w", signal_id, err)
	}
	pool_location, ok := req.PoolLocations[signal_id]
	if ok {
		log.Info().Float64("lat", pool_location.Latitude).Float64("lng", pool_location.Longitude).Msg("got pool location")
		geom_query := geom.PostgisPointQuery(pool_location.Longitude, pool_location.Latitude)
		_, err = psql.Update(
			um.Table("pool"),
			um.SetCol("geometry").To(geom_query),
			um.From("signal_pool"),
			um.Where(psql.Quote("signal_pool", "pool_id").EQ(psql.Quote("pool", "id"))),
			um.Where(psql.Quote("signal_pool", "signal_id").EQ(psql.Arg(signal_id))),
		).Exec(ctx, txn)
		if err != nil {
			return nil, nhttp.NewError("failed to update pool through signal %d: %w", signal_id, err)
		}
	}
	txn.Commit(ctx)

	return &createdLead{
		ID: lead.ID,
	}, nil
}
