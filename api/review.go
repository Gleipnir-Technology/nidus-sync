package api

import (
	"context"
	"errors"
	"net/http"

	nhttp "github.com/Gleipnir-Technology/nidus-sync/http"
	"github.com/Gleipnir-Technology/nidus-sync/platform"
)

type createReviewPool struct {
	Status  string               `json:"status"`
	TaskID  int32                `json:"task_id"`
	Updates *platform.PoolUpdate `json:"updates"`
}
type createdReviewPool struct{}

func postReviewPool(ctx context.Context, r *http.Request, user platform.User, req createReviewPool) (*createdReviewPool, *nhttp.ErrorWithStatus) {
	_, err := platform.ReviewPoolCreate(ctx, user, req.TaskID, req.Status, req.Updates)

	if err != nil {
		if errors.As(err, &platform.ErrorNotFound{}) {
			return nil, nhttp.NewErrorStatus(http.StatusNotFound, "review task %d not found", req.TaskID)
		}
		return nil, nhttp.NewError("failed to set review: %w", err)
	}
	return &createdReviewPool{}, nil
}
