package api

import (
	"context"
	"net/http"
	"time"

	"github.com/Gleipnir-Technology/bob"
	"github.com/Gleipnir-Technology/bob/dialect/psql"
	"github.com/Gleipnir-Technology/bob/dialect/psql/um"
	"github.com/Gleipnir-Technology/nidus-sync/db"
	"github.com/Gleipnir-Technology/nidus-sync/db/enums"
	"github.com/Gleipnir-Technology/nidus-sync/db/models"
	nhttp "github.com/Gleipnir-Technology/nidus-sync/http"
	"github.com/aarondl/opt/omit"
	"github.com/aarondl/opt/omitnull"
	"github.com/rs/zerolog/log"
	/*
	   "github.com/Gleipnir-Technology/bob/dialect/psql/sm"
	   "github.com/Gleipnir-Technology/nidus-sync/platform/geom"
	   "github.com/aarondl/opt/omit"
	   "github.com/aarondl/opt/omitnull"
	   "github.com/stephenafamo/scan"
	*/)

type reviewPoolUpdate struct {
	Condition *string  `json:"condition"`
	Latitude  *float32 `json:"latitude"`
	Longitude *float32 `json:"longitude"`
}
type createReviewPool struct {
	Status  string            `json:"status"`
	TaskID  int32             `json:"task_id"`
	Updates *reviewPoolUpdate `json:"updates"`
}
type createdReviewPool struct{}

func postReviewPool(ctx context.Context, r *http.Request, org *models.Organization, user *models.User, req createReviewPool) (*createdReviewPool, *nhttp.ErrorWithStatus) {
	txn, err := db.PGInstance.BobDB.BeginTx(ctx, nil)
	if err != nil {
		return nil, nhttp.NewError("start txn: %w", err)
	}
	defer txn.Rollback(ctx)
	review_task, err := models.ReviewTasks.Query(
		models.SelectWhere.ReviewTasks.ID.EQ(req.TaskID),
		models.SelectWhere.ReviewTasks.OrganizationID.EQ(org.ID),
	).One(ctx, txn)
	if err != nil {
		return nil, nhttp.NewErrorStatus(http.StatusNotFound, "review task %d not found", req.TaskID)
	}
	var resolution enums.Reviewtaskresolutiontype
	err = resolution.Scan(req.Status)
	if err != nil {
		return nil, nhttp.NewErrorStatus(http.StatusNotFound, "status '%s' is not recognized", req.Status)
	}
	review_task.Update(ctx, txn, &models.ReviewTaskSetter{
		Resolution: omitnull.From(resolution),
		Reviewed:   omitnull.From(time.Now()),
		ReviewerID: omitnull.From(user.ID),
	})
	review_task_pool, err := models.ReviewTaskPools.Query(
		models.SelectWhere.ReviewTaskPools.ReviewTaskID.EQ(review_task.ID),
	).One(ctx, txn)
	var e *nhttp.ErrorWithStatus
	switch req.Status {
	case "discarded":
		e = discardReviewPool(ctx, txn, user, req, review_task_pool)
	case "committed":
		e = commitReviewPool(ctx, txn, user, req, review_task_pool)
	default:
		return nil, nhttp.NewErrorStatus(http.StatusBadRequest, "unrecognized status %s", req.Status)
	}
	if e != nil {
		return nil, e
	}
	txn.Commit(ctx)
	log.Info().Int32("id", review_task.ID).Str("status", req.Status).Msg("committed")
	return &createdReviewPool{}, e
}
func discardReviewPool(ctx context.Context, txn bob.Tx, user *models.User, req createReviewPool, review_task_pool *models.ReviewTaskPool) *nhttp.ErrorWithStatus {
	return nil
}
func commitReviewPool(ctx context.Context, txn bob.Tx, user *models.User, req createReviewPool, review_task_pool *models.ReviewTaskPool) *nhttp.ErrorWithStatus {
	if req.Updates == nil {
		return nil
	}
	up := *req.Updates
	feature_pool, err := models.FindFeaturePool(ctx, txn, review_task_pool.FeaturePoolID)
	if err != nil {
		return nhttp.NewError("find feature pool: %w", err)
	}
	if up.Condition != nil {
		var condition enums.Poolconditiontype
		err := condition.Scan(up.Condition)
		if err != nil {
			return nhttp.NewErrorStatus(http.StatusBadRequest, "unrecognized condition %s", up.Condition)
		}
		err = review_task_pool.Update(ctx, txn, &models.ReviewTaskPoolSetter{
			Condition: omitnull.From(condition),
		})
		if err != nil {
			return nhttp.NewError("update rewiew task: %w", err)
		}
		err = feature_pool.Update(ctx, txn, &models.FeaturePoolSetter{
			Condition: omit.From(condition),
		})
		if err != nil {
			return nhttp.NewError("update feature_pool: %w", err)
		}
	}
	if up.Latitude != nil || up.Longitude != nil {
		if up.Latitude == nil || up.Longitude == nil {
			return nhttp.NewErrorStatus(http.StatusBadRequest, "you have to specify lat and lng together")
		}
		_, err = psql.Update(
			um.Table("review_task_pool"),
			um.SetCol("location").To(
				psql.F("ST_SetSRID",
					psql.F("ST_MakePoint",
						psql.Arg(*up.Longitude),
						psql.Arg(*up.Latitude),
					), psql.Arg(4326),
				),
			),
			um.Where(psql.Quote("review_task_pool", "review_task_id").EQ(psql.Arg(review_task_pool.ReviewTaskID))),
		).Exec(ctx, txn)
		if err != nil {
			return nhttp.NewError("save task: %w", err)
		}
		_, err = psql.Update(
			um.Table("feature"),
			um.SetCol("location").To(
				psql.F("ST_SetSRID",
					psql.F("ST_MakePoint",
						psql.Arg(*up.Longitude),
						psql.Arg(*up.Latitude),
					), psql.Arg(4326),
				),
			),
			um.Where(psql.Quote("feature", "id").EQ(psql.Arg(review_task_pool.FeaturePoolID))),
		).Exec(ctx, txn)
		if err != nil {
			return nhttp.NewError("save feature: %w", err)
		}
	}
	return nil
}
