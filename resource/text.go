package resource

import (
	"context"
	"net/http"

	//modelpublic "github.com/Gleipnir-Technology/nidus-sync/db/gen/nidus-sync/public/model"
	//modelpublicreport "github.com/Gleipnir-Technology/nidus-sync/db/gen/nidus-sync/publicreport/model"
	nhttp "github.com/Gleipnir-Technology/nidus-sync/http"
	"github.com/Gleipnir-Technology/nidus-sync/platform"
	//"github.com/gorilla/mux"
	//"github.com/rs/zerolog/log"
)

type textR struct {
	router *router
}

func Text(r *router) *textR {
	return &textR{
		router: r,
	}
}

type textResource struct {
	ID int
}

func (res *textR) Get(ctx context.Context, r *http.Request, user platform.User, query QueryParams) (textResource, *nhttp.ErrorWithStatus) {
	text_id, error_with_status := res.router.IDFromMux(r)
	if error_with_status != nil {
		return textResource{}, error_with_status
	}
	return textResource{ID: text_id}, nil
}
