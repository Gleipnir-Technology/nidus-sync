package rmo

import (
	//"github.com/Gleipnir-Technology/nidus-sync/html"
	"github.com/Gleipnir-Technology/nidus-sync/static"
	"github.com/gorilla/mux"
)

func Router(r *mux.Router) {
	/*
		r.HandleFunc("/submit-complete", getSubmitComplete).Methods("GET")

		r.HandleFunc("/district", getDistrictList).Methods("GET")
		r.HandleFunc("/district/{slug}", getRootDistrict).Methods("GET")
		r.HandleFunc("/district/{slug}/compliance", getDistrictCompliance).Methods("GET")
		r.HandleFunc("/district/{slug}/compliance/address", getDistrictComplianceAddress).Methods("GET")
		r.HandleFunc("/district/{slug}/compliance/complete", getDistrictComplianceComplete).Methods("GET")
		r.HandleFunc("/district/{slug}/compliance/concern", getDistrictComplianceConcern).Methods("GET")
		r.HandleFunc("/district/{slug}/compliance/contact", getDistrictComplianceContact).Methods("GET")
		r.HandleFunc("/district/{slug}/compliance/evidence", getDistrictComplianceEvidence).Methods("GET")
		r.HandleFunc("/district/{slug}/compliance/permission", getDistrictCompliancePermission).Methods("GET")
		r.HandleFunc("/district/{slug}/compliance/process", getDistrictComplianceProcess).Methods("GET")
		r.HandleFunc("/district/{slug}/compliance/submit", getDistrictComplianceSubmit).Methods("GET")
		r.HandleFunc("/district/{slug}/nuisance", getNuisanceDistrict).Methods("GET")
		//r.HandleFunc("/district/{slug}/nuisance-submit-complete", renderMock(mockNuisanceSubmitCompleteT)).Methods("GET")
		//r.HandleFunc("/district/{slug}/status", renderMock(mockStatusT)).Methods("GET")
		r.HandleFunc("/district/{slug}/water", getWaterDistrict).Methods("GET")
		//r.HandleFunc("/district/{slug}/water", postWaterDistrict).Methods("POST")
		r.HandleFunc("/error", getError).Methods("GET")

		r.HandleFunc("/privacy", getPrivacy).Methods("GET")
		r.HandleFunc("/robots.txt", getRobots).Methods("GET")
		r.HandleFunc("/email/render/{code}", getEmailByCode).Methods("GET")
		r.HandleFunc("/email/confirm", getEmailConfirm).Methods("GET")
		r.HandleFunc("/email/confirm", postEmailConfirm).Methods("POST")
		r.HandleFunc("/email/confirm/complete", getEmailConfirmComplete).Methods("GET")
		r.HandleFunc("/email/unsubscribe", getEmailUnsubscribe).Methods("GET")
		r.HandleFunc("/email/unsubscribe/report/{report_id}", getEmailReportUnsubscribe).Methods("GET")
		r.HandleFunc("/image/{uuid}", getImageByUUID).Methods("GET")
		r.HandleFunc("/mailer/{public_id}", html.MakeGet(getMailer)).Methods("GET")
		r.HandleFunc("/mailer/{public_id}/confirm", html.MakePost(postMailerConfirm)).Methods("POST")
		r.HandleFunc("/mailer/{public_id}/contribute", html.MakeGet(getMailerContribute)).Methods("GET")
		r.HandleFunc("/mailer/{public_id}/evidence", html.MakeGet(getMailerEvidence)).Methods("GET")
		r.HandleFunc("/mailer/{public_id}/schedule", html.MakeGet(getMailerSchedule)).Methods("GET")
		r.HandleFunc("/mailer/{public_id}/update", html.MakeGet(getMailerUpdate)).Methods("GET")
		r.HandleFunc("/register-notifications", postRegisterNotifications).Methods("POST")
		r.HandleFunc("/register-notifications-complete", getRegisterNotificationsComplete).Methods("GET")
		r.HandleFunc("/report/suggest", getReportSuggestion).Methods("GET")
		r.HandleFunc("/scss/*", getScssDebug).Methods("GET")
		r.HandleFunc("/status", getStatus).Methods("GET")
		r.HandleFunc("/status/{report_id}", getStatusByID).Methods("GET")
		r.HandleFunc("/terms-of-service", getTerms).Methods("GET")
	*/
	static.AddStaticRoute(r, "/static")
	r.PathPrefix("/").Handler(static.SinglePageApp("static/gen/rmo"))
}
