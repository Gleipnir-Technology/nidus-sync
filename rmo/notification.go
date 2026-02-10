package rmo

import (
	"fmt"
	"net/http"

	"github.com/Gleipnir-Technology/nidus-sync/db"
	"github.com/Gleipnir-Technology/nidus-sync/platform/report"
	"github.com/Gleipnir-Technology/nidus-sync/platform/text"
	"github.com/rs/zerolog/log"
)

func postRegisterNotifications(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		respondError(w, "Failed to parse form", err, http.StatusBadRequest)
		return
	}
	has_consent := boolFromForm(r, "consent")
	has_notification := boolFromForm(r, "notification")
	has_subscribe := boolFromForm(r, "subscribe")
	email := r.PostFormValue("email")
	name := r.PostFormValue("name")
	phone_str := r.PostFormValue("phone")
	report_id := r.PostFormValue("report_id")

	var phone *text.E164
	if phone_str != "" {
		phone, err = text.ParsePhoneNumber(phone_str)
		if err != nil {
			http.Redirect(w, r, fmt.Sprintf("/error?code=invalid-phone&report=%s", report_id), http.StatusFound)
			return
		}
	}

	ctx := r.Context()
	txn, err := db.PGInstance.BobDB.BeginTx(ctx, nil)
	if err != nil {
		log.Error().Err(err).Msg("Failed to begin transaction")
		http.Redirect(w, r, fmt.Sprintf("/error?code=transaction-failed&report=%s", report_id), http.StatusFound)
		return
	}
	defer txn.Rollback(ctx)
	e := report.SaveReporter(ctx, txn, report_id, name, email, phone, has_consent)
	if e != nil {
		log.Error().Err(e).Str("name", name).Msg("Failed to save reporter")
		http.Redirect(w, r, fmt.Sprintf("/error?code=%s&report=%s", e.Code(), report_id), http.StatusFound)
		return
	}
	if email != "" {
		if has_subscribe {
			e := report.RegisterSubscriptionEmail(ctx, txn, email)
			if e != nil {
				log.Error().Err(e).Str("email", email).Msg("Failed to register subscription email")
				http.Redirect(w, r, fmt.Sprintf("/error?code=%s&report=%s", e.Code(), report_id), http.StatusFound)
			}
		}
		if has_notification {
			e := report.RegisterNotificationEmail(ctx, txn, report_id, email)
			if e != nil {
				log.Error().Err(e).Str("email", email).Msg("Failed to register notification email")
				http.Redirect(w, r, fmt.Sprintf("/error?code=%s&report=%s", e.Code(), report_id), http.StatusFound)
			}
		}
	}
	if phone != nil {
		if has_subscribe {
			e := report.RegisterSubscriptionPhone(ctx, txn, *phone)
			if e != nil {
				log.Error().Err(e).Str("phone", phone_str).Msg("Failed to register subscription phone")
				http.Redirect(w, r, fmt.Sprintf("/error?code=%s&report=%s", e.Code(), report_id), http.StatusFound)
			}
		}
		if has_notification {
			e := report.RegisterNotificationPhone(ctx, txn, report_id, *phone)
			if e != nil {
				log.Error().Err(e).Str("phone", phone_str).Msg("Failed to register notification phone")
				http.Redirect(w, r, fmt.Sprintf("/error?code=%s&report=%s", e.Code(), report_id), http.StatusFound)
			}
		}
	}
	txn.Commit(ctx)
	http.Redirect(w, r, fmt.Sprintf("/register-notifications-complete?report=%s", report_id), http.StatusFound)
}
