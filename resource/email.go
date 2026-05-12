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

type emailR struct {
	router *router
}

func Email(r *router) *emailR {
	return &emailR{
		router: r,
	}
}

type email struct {
	ID int
}

func (res *emailR) Get(ctx context.Context, r *http.Request, user platform.User, query QueryParams) (email, *nhttp.ErrorWithStatus) {
	email_id, error_with_status := res.router.IDFromMux(r)
	if error_with_status != nil {
		return email{}, error_with_status
	}
	return email{ID: email_id}, nil
}
