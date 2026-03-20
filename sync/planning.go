package sync

import (
	"context"
	"fmt"
	"html/template"
	"net/http"

	"github.com/Gleipnir-Technology/nidus-sync/config"
	"github.com/Gleipnir-Technology/nidus-sync/html"
	nhttp "github.com/Gleipnir-Technology/nidus-sync/http"
	"github.com/Gleipnir-Technology/nidus-sync/platform"
	//"github.com/rs/zerolog/log"
)

type contentPlanningRoot struct {
	URLTiles template.HTMLAttr
}

func getPlanningRoot(ctx context.Context, r *http.Request, user platform.User) (*html.Response[contentPlanningRoot], *nhttp.ErrorWithStatus) {
	return html.NewResponse("sync/planning-root.html", contentPlanningRoot{
		URLTiles: template.HTMLAttr(fmt.Sprintf(`url-tiles="%s"`, config.MakeURLNidus("/api/tile/{z}/{y}/{x}"))),
	}), nil
}
