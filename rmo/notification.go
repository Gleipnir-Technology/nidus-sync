package rmo

import (
	"fmt"
	"net/http"

	"github.com/Gleipnir-Technology/nidus-sync/platform/report"
	"github.com/Gleipnir-Technology/nidus-sync/platform/text"
)

var (
	registerNotificationsCompleteT = buildTemplate("register-notifications-complete", "base")
)

func postRegisterNotifications(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		respondError(w, "Failed to parse form", err, http.StatusBadRequest)
		return
	}
	consent := r.PostFormValue("consent")
	email := r.PostFormValue("email")
	phone_str := r.PostFormValue("phone")
	report_id := r.PostFormValue("report_id")
	subscribe := postFormBool(r, "subscribe")
	if consent != "on" {
		respondError(w, "You must consent", nil, http.StatusBadRequest)
		return
	}
	if email == "" && phone_str == "" {
		http.Redirect(w, r, fmt.Sprintf("/submit-complete?report=%s", report_id), http.StatusFound)
		return
	}
	phone, err := text.ParsePhoneNumber(phone_str)
	if err != nil {
		http.Redirect(w, r, fmt.Sprintf("/error?code=invalid-phone&report=%s", report_id), http.StatusFound)
		return
	}

	ctx := r.Context()
	if subscribe != nil && *subscribe {
		if email != "" {
			e := report.RegisterSubscriptionEmail(ctx, email)
			if e != nil {
				http.Redirect(w, r, fmt.Sprintf("/error?code=%s&report=%s", report_id, e.Code()), http.StatusFound)
			}
		}
		if phone_str != "" {
			e := report.RegisterSubscriptionPhone(ctx, phone)
			if e != nil {
				http.Redirect(w, r, fmt.Sprintf("/error?code=%s&report=%s", report_id, e.Code()), http.StatusFound)
			}
		}
	}
	e := report.RegisterNotifications(ctx, report_id, email, phone)
	if e != nil {
		http.Redirect(w, r, fmt.Sprintf("/error?code=%s&report=%s", report_id, e.Code()), http.StatusFound)
	} else {
		http.Redirect(w, r, fmt.Sprintf("/register-notifications-complete?report=%s", report_id), http.StatusFound)
	}
}
