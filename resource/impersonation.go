package resource

import (
	"context"
	"net/http"

	"github.com/Gleipnir-Technology/nidus-sync/auth"
	nhttp "github.com/Gleipnir-Technology/nidus-sync/http"
	"github.com/Gleipnir-Technology/nidus-sync/platform"
	"github.com/aarondl/opt/omit"
	"github.com/rs/zerolog/log"
)

func Impersonation(r *router) *impersonationR {
	return &impersonationR{
		router: r,
	}
}

type impersonationR struct {
	router *router
}
type impersonation struct {
	UserID omit.Val[int] `json:"id"`
}

func (res *impersonationR) Create(ctx context.Context, r *http.Request, u platform.User, i impersonation) (*impersonation, *nhttp.ErrorWithStatus) {
	if i.UserID.IsUnset() {
		return nil, nhttp.NewBadRequest("you must provide an 'id'")
	}
	target_id := i.UserID.MustGet()
	l, err := platform.ImpersonationCreate(ctx, u, target_id)
	if err != nil {
		return nil, nhttp.NewError("create impersonation: %w", err)
	}
	auth.ImpersonateUser(ctx, target_id)
	log.Info().Int("user.id", u.ID).Str("username", u.Username).Int("target.id", target_id).Int32("log.id", l.ID).Msg("Impersonation begins")
	return &impersonation{
		UserID: i.UserID,
	}, nil
}
func (res *impersonationR) Delete(ctx context.Context, r *http.Request, u platform.User) *nhttp.ErrorWithStatus {
	if auth.ImpersonatedUser == nil {
		return nhttp.NewBadRequest("not impersonating")
	}
	real_user_id := auth.ImpersonatorID(ctx)
	if real_user_id == nil {
		return nhttp.NewError("No impersonator ID")
	}
	err := platform.ImpersonationEnd(ctx, u, *real_user_id)
	if err != nil {
		return nhttp.NewError("end impersonation: %w", err)
	}
	auth.ImpersonateEnd(ctx)
	return nil
}
