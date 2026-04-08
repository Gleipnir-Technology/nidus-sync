package resource

import (
	"context"
	nhttp "github.com/Gleipnir-Technology/nidus-sync/http"
	"github.com/Gleipnir-Technology/nidus-sync/platform"
	"github.com/Gleipnir-Technology/nidus-sync/platform/text"
	"github.com/Gleipnir-Technology/nidus-sync/platform/types"
	"github.com/rs/zerolog/log"
	"net/http"
)

type publicreportNotificationR struct {
	router *router
}

type publicreportNotification struct {
	Consent      bool   `json:"consent"`
	Email        string `json:"email"`
	Name         string `json:"name"`
	Notification bool   `json:"notification"`
	Phone        string `json:"phone"`
	ReportID     string `json:"report_id"`
	Subscription bool   `json:"subscription"`
}

func PublicreportNotification(r *router) *publicreportNotificationR {
	return &publicreportNotificationR{
		router: r,
	}
}

func (res *publicreportNotificationR) Create(ctx context.Context, r *http.Request, n publicreportNotification) (*publicreportNotification, *nhttp.ErrorWithStatus) {
	var err error
	var phone *types.E164
	if n.Phone != "" {
		phone, err = text.ParsePhoneNumber(n.Phone)
		if err != nil {
			return nil, nhttp.NewBadRequest("can't parse phone: %w", err)
		}
	}
	err = platform.PublicreportNotificationCreate(ctx, platform.PublicreportNotification{
		Consent:      n.Consent,
		Email:        n.Email,
		Name:         n.Name,
		Notification: n.Notification,
		Phone:        phone,
		ReportID:     n.ReportID,
		Subscription: n.Subscription,
	})
	if err != nil {
		return nil, nhttp.NewError("create notification: %w", err)
	}
	log.Info().Str("name", n.Name).Str("email", n.Email).Str("phone", n.Phone).Str("report_id", n.ReportID).Msg("Added reporter data")
	return &n, nil
}
