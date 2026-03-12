package sync

import (
	"context"
	"net/http"

	"github.com/Gleipnir-Technology/nidus-sync/html"
	nhttp "github.com/Gleipnir-Technology/nidus-sync/http"
	"github.com/Gleipnir-Technology/nidus-sync/platform"
)

type contentRadar struct {
	Organization platform.Organization
}

func getRadar(ctx context.Context, r *http.Request, user platform.User) (*html.Response[contentRadar], *nhttp.ErrorWithStatus) {
	data := contentRadar{
		Organization: user.Organization,
	}
	return html.NewResponse("sync/radar.html", data), nil
}
