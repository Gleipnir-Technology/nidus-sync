package api

import (
	"context"
	"net/http"
	"time"

	"github.com/Gleipnir-Technology/bob"
	"github.com/Gleipnir-Technology/bob/dialect/psql"
	"github.com/Gleipnir-Technology/bob/dialect/psql/sm"
	"github.com/Gleipnir-Technology/nidus-sync/db"
	"github.com/Gleipnir-Technology/nidus-sync/db/models"
	nhttp "github.com/Gleipnir-Technology/nidus-sync/http"
	"github.com/Gleipnir-Technology/nidus-sync/platform"
	"github.com/Gleipnir-Technology/nidus-sync/platform/types"
	//"github.com/aarondl/opt/null"
	"github.com/stephenafamo/scan"
)

type reviewTaskPool struct {
	Address   types.Address  `json:"address"`
	Condition string         `json:"condition"`
	Created   time.Time      `json:"created"`
	Creator   platform.User  `json:"creator"`
	ID        int32          `json:"id"`
	Location  Location       `json:"location"`
	Reviewed  *time.Time     `json:"addressed"`
	Reviewer  *platform.User `json:"addressor"`
}
type contentListReviewTaskPool struct {
	Tasks []reviewTaskPool `json:"tasks"`
}

func listReviewTaskPool(ctx context.Context, r *http.Request, org *models.Organization, user *models.User, query queryParams) (*contentListReviewTaskPool, *nhttp.ErrorWithStatus) {
	type _Row struct {
		Address    types.Address `db:"address"`
		Condition  string        `db:"condition"`
		Created    time.Time     `db:"created"`
		CreatorID  int32         `db:"creator_id"`
		ID         int32         `db:"id"`
		Latitude   float64       `db:"latitude"`
		Longitude  float64       `db:"longitude"`
		Reviewed   *time.Time    `db:"reviewed"`
		ReviewerID *int32        `db:"reviewer_id"`
		Species    *string       `db:"species"`
		Title      string        `db:"title"`
		Type       string        `db:"type"`
	}
	limit := 20
	if query.Limit != nil {
		limit = *query.Limit
	}
	rows, err := bob.All(ctx, db.PGInstance.BobDB, psql.Select(
		sm.Columns(
			"feature_pool.condition AS condition",
			"review_task.created AS created",
			"review_task.creator_id AS creator_id",
			"review_task.id AS id",
			"review_task.reviewed AS reviewed",
			"review_task.reviewer_id AS reviewer_id",
			"address.country AS \"address.country\"",
			"address.locality AS \"address.locality\"",
			"address.number_ AS \"address.number\"",
			"address.postal_code AS \"address.postal_code\"",
			"address.region AS \"address.region\"",
			"address.street AS \"address.street\"",
			"address.unit AS \"address.unit\"",
			"ST_Y(address.location) AS latitude",
			"ST_X(address.location) AS longitude",
		),
		sm.From("review_task_pool"),
		sm.InnerJoin("feature_pool").OnEQ(
			psql.Quote("review_task_pool", "feature_pool_id"),
			psql.Quote("feature_pool", "feature_id"),
		),
		sm.InnerJoin("review_task").OnEQ(
			psql.Quote("review_task_pool", "review_task_id"),
			psql.Quote("review_task", "id"),
		),
		sm.InnerJoin("feature").OnEQ(
			psql.Quote("feature_pool", "feature_id"),
			psql.Quote("feature", "id"),
		),
		sm.InnerJoin("site").On(
			psql.And(
				psql.Quote("feature", "site_id").EQ(psql.Quote("site", "id")),
				psql.Quote("feature", "site_version").EQ(psql.Quote("site", "version")),
			),
		),
		sm.InnerJoin("address").OnEQ(
			psql.Quote("site", "address_id"),
			psql.Quote("address", "id"),
		),
		sm.Where(psql.Quote("review_task", "organization_id").EQ(psql.Arg(org.ID))),
		sm.Where(psql.Quote("review_task", "reviewed").IsNull()),
		sm.Limit(limit),
	), scan.StructMapper[_Row]())
	if err != nil {
		return nil, nhttp.NewError("failed to get signals: %w", err)
	}
	users_by_id, err := platform.UsersByID(ctx, org)
	if err != nil {
		return nil, nhttp.NewError("users by id: %w", err)
	}
	tasks := make([]reviewTaskPool, len(rows))
	for i, row := range rows {
		tasks[i] = reviewTaskPool{
			Address:   row.Address,
			Condition: row.Condition,
			Created:   row.Created,
			Creator:   *users_by_id[row.CreatorID],
			ID:        row.ID,
			Location: Location{
				Latitude:  row.Latitude,
				Longitude: row.Longitude,
			},
			Reviewed: row.Reviewed,
			Reviewer: userOrNil(users_by_id, row.ReviewerID),
		}
	}
	return &contentListReviewTaskPool{
		Tasks: tasks,
	}, nil
}
