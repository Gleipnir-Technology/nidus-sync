package api

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	nhttp "github.com/Gleipnir-Technology/nidus-sync/http"
	"github.com/Gleipnir-Technology/nidus-sync/platform"
)

type createReviewPool struct {
	Status  string               `json:"status"`
	TaskID  int32                `json:"task_id"`
	Updates *platform.PoolUpdate `json:"updates"`
}

func postReviewPool(ctx context.Context, r *http.Request, user platform.User, req createReviewPool) (string, *nhttp.ErrorWithStatus) {
	id, err := platform.ReviewPoolCreate(ctx, user, req.TaskID, req.Status, req.Updates)

	if err != nil {
		if errors.As(err, &platform.ErrorNotFound{}) {
			return "", nhttp.NewErrorStatus(http.StatusNotFound, "review task %d not found", req.TaskID)
		}
		return "", nhttp.NewError("failed to set review: %w", err)
	}
	return fmt.Sprintf("/review/%d", id), nil
}
