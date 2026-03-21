package platform

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/Gleipnir-Technology/bob"
	"github.com/Gleipnir-Technology/bob/dialect/psql"
	"github.com/Gleipnir-Technology/bob/dialect/psql/um"
	"github.com/Gleipnir-Technology/nidus-sync/db"
	"github.com/Gleipnir-Technology/nidus-sync/db/enums"
	"github.com/Gleipnir-Technology/nidus-sync/db/models"
	nhttp "github.com/Gleipnir-Technology/nidus-sync/http"
	"github.com/Gleipnir-Technology/nidus-sync/platform/event"
	"github.com/aarondl/opt/omit"
	"github.com/aarondl/opt/omitnull"
	"github.com/rs/zerolog/log"
)

type PoolUpdate struct {
	Condition *string  `json:"condition"`
	Latitude  *float32 `json:"latitude"`
	Longitude *float32 `json:"longitude"`
}

func ReviewPoolCreate(ctx context.Context, user User, task_id int32, status string, update *PoolUpdate) (int32, error) {
	txn, err := db.PGInstance.BobDB.BeginTx(ctx, nil)
	if err != nil {
		return 0, nhttp.NewError("start txn: %w", err)
	}
	defer txn.Rollback(ctx)
	review_task, err := models.ReviewTasks.Query(
		models.SelectWhere.ReviewTasks.ID.EQ(task_id),
		models.SelectWhere.ReviewTasks.OrganizationID.EQ(user.Organization.ID()),
	).One(ctx, txn)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			return 0, newNotFound("review task %d", task_id)
		}
		return 0, fmt.Errorf("find review task: %d", task_id)
	}
	var resolution enums.Reviewtaskresolutiontype
	err = resolution.Scan(status)
	if err != nil {
		return 0, fmt.Errorf("status '%s' is not recognized: %w", status, err)
	}
	review_task.Update(ctx, txn, &models.ReviewTaskSetter{
		Resolution: omitnull.From(resolution),
		Reviewed:   omitnull.From(time.Now()),
		ReviewerID: omitnull.From(int32(user.ID)),
	})
	review_task_pool, err := models.ReviewTaskPools.Query(
		models.SelectWhere.ReviewTaskPools.ReviewTaskID.EQ(review_task.ID),
	).One(ctx, txn)
	switch status {
	case "discarded":
		// Nothing to do, we already discarded it above
	case "committed":
		err = commitReviewPool(ctx, txn, user, review_task_pool, update)
	default:
		return 0, nhttp.NewErrorStatus(http.StatusBadRequest, "unrecognized status %s", status)
	}
	if err != nil {
		return 0, err
	}
	event.Updated(event.TypeReviewTask, user.Organization.ID(), strconv.Itoa(int(review_task.ID)))
	txn.Commit(ctx)
	log.Info().Int32("id", review_task.ID).Str("status", status).Msg("review completed")
	return review_task.ID, err
}
func commitReviewPool(ctx context.Context, txn bob.Tx, user User, review_task_pool *models.ReviewTaskPool, update *PoolUpdate) error {
	if update == nil {
		return nil
	}
	feature_pool, err := models.FindFeaturePool(ctx, txn, review_task_pool.FeaturePoolID)
	if err != nil {
		return nhttp.NewError("find feature pool: %w", err)
	}
	condition := feature_pool.Condition
	if update.Condition != nil {
		err := condition.Scan(*update.Condition)
		if err != nil {
			return nhttp.NewErrorStatus(http.StatusBadRequest, "unrecognized condition %s", update.Condition)
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
	if update.Latitude != nil || update.Longitude != nil {
		if update.Latitude == nil || update.Longitude == nil {
			return nhttp.NewErrorStatus(http.StatusBadRequest, "you have to specify lat and lng together")
		}
		_, err = psql.Update(
			um.Table("review_task_pool"),
			um.SetCol("location").To(
				psql.F("ST_SetSRID",
					psql.F("ST_MakePoint",
						psql.Arg(update.Longitude),
						psql.Arg(update.Latitude),
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
						psql.Arg(*update.Longitude),
						psql.Arg(*update.Latitude),
					), psql.Arg(4326),
				),
			),
			um.Where(psql.Quote("feature", "id").EQ(psql.Arg(review_task_pool.FeaturePoolID))),
		).Exec(ctx, txn)
		if err != nil {
			return nhttp.NewError("save feature: %w", err)
		}
	}
	log.Debug().Str("condition", string(condition)).Int32("id", review_task_pool.ReviewTaskID).Msg("checking")
	// if the pool is either murkey or green, immediately create a signal from it
	if condition == enums.PoolconditiontypeGreen || condition == enums.PoolconditiontypeMurky {
		feature, err := models.FindFeature(ctx, txn, feature_pool.FeatureID)
		if err != nil {
			return nhttp.NewError("find feature %d: %w", feature_pool.FeatureID, err)
		}
		signal, err := models.Signals.Insert(&models.SignalSetter{
			Addressed:            omitnull.FromPtr[time.Time](nil),
			Addressor:            omitnull.FromPtr[int32](nil),
			Created:              omit.From(time.Now()),
			Creator:              omit.From[int32](int32(user.ID)),
			FeaturePoolFeatureID: omitnull.From(feature_pool.FeatureID),
			//ID: omit.Val[int32],
			OrganizationID: omit.From(user.Organization.ID()),
			ReportID:       omitnull.FromPtr[int32](nil),
			Species:        omitnull.FromPtr[enums.Mosquitospecies](nil),
			Type:           omit.From(enums.SignaltypeFlyoverPool),
			SiteID:         omitnull.From(feature.SiteID),
			Location:       omit.From[string](feature.Location.GetOr("")),
		}).One(ctx, txn)
		if err != nil {
			return nhttp.NewError("create signal: %w", err)
		}
		event.Created(event.TypeSignal, user.Organization.ID(), strconv.Itoa(int(signal.ID)))
		log.Debug().Int32("id", signal.ID).Msg("created pool signal")
	}
	return nil
}
