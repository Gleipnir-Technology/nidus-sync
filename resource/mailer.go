package resource

import (
	"context"
	"net/http"
	"strconv"

	nhttp "github.com/Gleipnir-Technology/nidus-sync/http"
	"github.com/Gleipnir-Technology/nidus-sync/platform"
	"github.com/Gleipnir-Technology/nidus-sync/platform/types"
	//"github.com/aarondl/opt/null"
	"github.com/gorilla/mux"
)

type mailerR struct {
	router *router
}

func Mailer(r *router) *mailerR {
	return &mailerR{
		router: r,
	}
}

func (res *mailerR) ByIDGet(ctx context.Context, r *http.Request, user platform.User, query QueryParams) (*types.Mailer, *nhttp.ErrorWithStatus) {
	vars := mux.Vars(r)
	id_str := vars["id"]
	id, err := strconv.Atoi(id_str)
	if err != nil {
		return nil, nhttp.NewBadRequest("'%s' is not a valid mailer ID: %w", id_str, err)
	}
	mailer, err := platform.MailerByID(ctx, user, int32(id))
	if err != nil {
		return nil, nhttp.NewError("mailer by id: %w", err)
	}
	return mailer, nil
}
func (res *mailerR) List(ctx context.Context, r *http.Request, user platform.User, query QueryParams) ([]*types.Mailer, *nhttp.ErrorWithStatus) {
	limit := 1000
	if query.Limit != nil {
		limit = *query.Limit
	}
	mailers, err := platform.MailerList(ctx, user, limit)
	if err != nil {
		return nil, nhttp.NewError("list signals: %w", err)
	}
	for _, mailer := range mailers {
		uri, err := res.router.IDToURI("mailer.ByIDGet", int(mailer.ID))
		if err != nil {
			return nil, nhttp.NewError("set uri: %w", err)
		}
		mailer.URI = uri
	}
	return mailers, nil
}
