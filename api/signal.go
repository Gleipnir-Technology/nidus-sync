package api

import (
	"context"
	"net/http"
	"time"

	"github.com/Gleipnir-Technology/bob"
	"github.com/Gleipnir-Technology/bob/dialect/psql"
	"github.com/Gleipnir-Technology/bob/dialect/psql/sm"
	"github.com/Gleipnir-Technology/nidus-sync/db"
	nhttp "github.com/Gleipnir-Technology/nidus-sync/http"
	"github.com/Gleipnir-Technology/nidus-sync/platform"
	"github.com/Gleipnir-Technology/nidus-sync/platform/types"
	//"github.com/aarondl/opt/null"
	"github.com/stephenafamo/scan"
)

type signal struct {
	Address   types.Address  `json:"address"`
	Addressed *time.Time     `json:"addressed"`
	Addressor *platform.User `json:"addressor"`
	Created   time.Time      `json:"created"`
	Creator   platform.User  `json:"creator"`
	ID        int32          `json:"id"`
	Location  Location       `json:"location"`
	Species   string         `json:"species"`
	Title     string         `json:"title"`
	Type      string         `json:"type"`
}
type contentListSignal struct {
	Signals []signal `json:"signals"`
}

func listSignal(ctx context.Context, r *http.Request, user platform.User, query queryParams) (*contentListSignal, *nhttp.ErrorWithStatus) {
	type _Row struct {
		Address   types.Address `db:"address"`
		Addressed *time.Time    `db:"addressed"`
		Addressor *int32        `db:"addressor"`
		Created   time.Time     `db:"created"`
		Creator   int32         `db:"creator_id"`
		ID        int32         `db:"id"`
		Latitude  float64       `db:"latitude"`
		Longitude float64       `db:"longitude"`
		Location  Location      `db:"location"`
		Species   *string       `db:"species"`
		Title     string        `db:"title"`
		Type      string        `db:"type"`
	}
	limit := 20
	if query.Limit != nil {
		limit = *query.Limit
	}
	rows, err := bob.All(ctx, db.PGInstance.BobDB, psql.Select(
		sm.Columns(
			"signal.addressed AS addressed",
			"signal.addressor AS addressor",
			"signal.created AS created",
			"signal.creator AS creator_id",
			"signal.id AS id",
			"signal.species AS species",
			"signal.title AS title",
			"signal.type_ AS type",
			"address.country AS \"address.country\"",
			"address.locality AS \"address.locality\"",
			"address.number_ AS \"address.number\"",
			"address.postal_code AS \"address.postal_code\"",
			"address.region AS \"address.region\"",
			"address.street AS \"address.street\"",
			"address.unit AS \"address.unit\"",
			"ST_Y(address.geom) AS latitude",
			"ST_X(address.geom) AS longitude",
		),
		sm.From("signal"),
		sm.InnerJoin("signal_pool").OnEQ(
			psql.Quote("signal", "id"),
			psql.Quote("signal_pool", "signal_id"),
		),
		sm.InnerJoin("pool").OnEQ(
			psql.Quote("signal_pool", "pool_id"),
			psql.Quote("pool", "id"),
		),
		sm.InnerJoin("site").On(
			psql.Quote("pool", "site_id").EQ(psql.Quote("site", "id")),
		),
		sm.InnerJoin("address").OnEQ(
			psql.Quote("site", "address_id"),
			psql.Quote("address", "id"),
		),
		sm.Where(psql.Quote("signal", "organization_id").EQ(psql.Arg(user.Organization.ID()))),
		sm.Where(psql.Quote("signal", "addressed").IsNull()),
		sm.Limit(limit),
	), scan.StructMapper[_Row]())

	/*
		rows, err := models.Signals.Query(
			models.SelectWhere.Signals.OrganizationID.EQ(org.ID),
			sm.OrderBy("created").Desc(),
		).All(ctx, db.PGInstance.BobDB)
	*/
	if err != nil {
		return nil, nhttp.NewError("failed to get signals: %w", err)
	}
	users_by_id, err := platform.UsersByOrg(ctx, user.Organization)
	if err != nil {
		return nil, nhttp.NewError("users by id: %w", err)
	}
	signals := make([]signal, len(rows))
	for i, row := range rows {
		var species string = ""
		if row.Species != nil {
			species = *row.Species
		}
		signals[i] = signal{
			Address:   row.Address,
			Addressed: row.Addressed,
			Addressor: userOrNil(users_by_id, row.Addressor),
			Created:   row.Created,
			Creator:   *users_by_id[row.Creator],
			ID:        row.ID,
			Location: Location{
				Latitude:  row.Latitude,
				Longitude: row.Longitude,
			},
			Species: species,
			Title:   row.Title,
			Type:    row.Type,
		}
	}
	return &contentListSignal{
		Signals: signals,
	}, nil
}

func userOrNil(usersByID map[int32]*platform.User, id *int32) *platform.User {
	if id == nil {
		return nil
	}
	u, ok := usersByID[*id]
	if !ok {
		return nil
	}
	return u
}
