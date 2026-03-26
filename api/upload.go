package api
import (
	"context"
	"net/http"
	"strconv"

	nhttp "github.com/Gleipnir-Technology/nidus-sync/http"
	"github.com/Gleipnir-Technology/nidus-sync/platform"
	"github.com/go-chi/chi/v5"
	//"github.com/rs/zerolog/log"
)
func getUploadByID(ctx context.Context, r *http.Request, u platform.User, query queryParams) (*platform.UploadPoolDetail, *nhttp.ErrorWithStatus) {
	file_id_str := chi.URLParam(r, "id")
	file_id_, err := strconv.ParseInt(file_id_str, 10, 32)
	if err != nil {
		return nil, nhttp.NewError("Failed to parse file_id: %w", err)
	}
	file_id := int32(file_id_)
	detail, err := platform.GetUploadDetail(ctx, u.Organization.ID, file_id)
	if err != nil {
		return nil, nhttp.NewError("Failed to get pool: %w", err)
	}
	return detail, nil
}
