package resource

import (
	"context"
	"net/http"

	"github.com/Gleipnir-Technology/nidus-sync/auth"
	"github.com/Gleipnir-Technology/nidus-sync/config"
	"github.com/Gleipnir-Technology/nidus-sync/html"
	nhttp "github.com/Gleipnir-Technology/nidus-sync/http"
	"github.com/Gleipnir-Technology/nidus-sync/platform"
	"github.com/Gleipnir-Technology/nidus-sync/platform/types"
)

type sessionR struct {
	router *router
}

func Session(r *router) *sessionR {
	return &sessionR{
		router: r,
	}
}

type organization struct {
	ID          int32              `json:"id"`
	ServiceArea *types.ServiceArea `json:"service_area"`
}

type session struct {
	Impersonating      *string                   `json:"impersonating"`
	NotificationCounts sessionNotificationCounts `json:"notification_counts"`
	Organization       organization              `json:"organization"`
	Self               user                      `json:"self"`
	URLs               sessionURL                `json:"urls"`
}
type sessionNotificationCounts struct {
	Communications uint `json:"communication"`
	Home           uint `json:"home"`
	Review         uint `json:"review"`
}

type sessionURL struct {
	API    sessionURLAPI `json:"api"`
	Tegola string        `json:"tegola"`
	Tile   string        `json:"tile"`
}
type sessionURLAPI struct {
	Avatar              string `json:"avatar"`
	Communication       string `json:"communication"`
	Impersonation       string `json:"impersonation"`
	PublicreportMessage string `json:"publicreport_message"`
	ReviewTask          string `json:"review_task"`
	ServiceRequest      string `json:"service_request"`
	Signal              string `json:"signal"`
	Site                string `json:"site"`
	Sync                string `json:"sync"`
	Upload              string `json:"upload"`
	User                string `json:"user"`
}

func (res *sessionR) Get(ctx context.Context, r *http.Request, user platform.User, query QueryParams) (*session, *nhttp.ErrorWithStatus) {
	urls := html.NewContentURL()
	counts, err := platform.NotificationCountsForUser(ctx, user)
	if err != nil {
		return nil, nhttp.NewError("get counst: %w", err)
	}
	usr := User(res.router)
	u, err := usr.response(&user)
	if err != nil {
		return nil, nhttp.NewError("create user: %w", err)
	}
	var impersonating *string
	impersonating_id := auth.ImpersonatedUser(ctx)
	if impersonating_id != nil {
		i, err := res.router.IDToURI("user.ByIDGet", int(*impersonating_id))
		if err != nil {
			return nil, nhttp.NewError("create impersonating uri: %w", err)
		}
		impersonating = &i
	}
	return &session{
		Impersonating: impersonating,
		NotificationCounts: sessionNotificationCounts{
			Communications: counts.Communications,
			Home:           counts.Home,
			Review:         counts.Review,
		},
		Organization: organization{
			ID:          user.Organization.ID,
			ServiceArea: user.Organization.ServiceArea,
		},
		Self: *u,
		URLs: sessionURL{
			API: sessionURLAPI{
				Avatar:              config.MakeURLNidus("/api/avatar"),
				Communication:       urls.API.Communication,
				Impersonation:       config.MakeURLNidus("/api/impersonation"),
				PublicreportMessage: urls.API.Publicreport.Message,
				ReviewTask:          config.MakeURLNidus("/api/review-task"),
				ServiceRequest:      config.MakeURLNidus("/api/service-request"),
				Signal:              config.MakeURLNidus("/api/signal"),
				Site:                config.MakeURLNidus("/api/site"),
				Sync:                config.MakeURLNidus("/api/sync"),
				Upload:              config.MakeURLNidus("/api/upload"),
				User:                config.MakeURLNidus("/api/user"),
			},
			Tegola: urls.Tegola,
			Tile:   config.MakeURLNidus("/api/tile/{z}/{y}/{x}"),
		},
	}, nil
}
