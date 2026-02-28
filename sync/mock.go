package sync

import (
	"fmt"
	"github.com/Gleipnir-Technology/nidus-sync/html"
	"github.com/go-chi/chi/v5"
	"net/http"
	//"github.com/rs/zerolog/log"
)

// Unauthenticated pages
/*
	admin                           = buildTemplate("admin", "base")
	dataEntry                       = buildTemplate("data-entry", "base")
	dataEntryBad                    = buildTemplate("data-entry-bad", "base")
	dispatch                        = buildTemplate("dispatch", "base")
	dispatchResults                 = buildTemplate("dispatch-results", "base")
	mockRoot                        = buildTemplate("mock-root", "base")
	reportPage                      = buildTemplate("report", "base")
	reportConfirmation              = buildTemplate("report-confirmation", "base")
	reportContribute                = buildTemplate("report-contribute", "base")
	reportDetail                    = buildTemplate("report-detail", "base")
	reportEvidence                  = buildTemplate("report-evidence", "base")
	reportSchedule                  = buildTemplate("report-schedule", "base")
	reportUpdate                    = buildTemplate("report-update", "base")
	serviceRequest                  = buildTemplate("service-request", "base")
	serviceRequestDetail            = buildTemplate("service-request-detail", "base")
	serviceRequestLocation          = buildTemplate("service-request-location", "base")
	serviceRequestMosquito          = buildTemplate("service-request-mosquito", "base")
	serviceRequestPool              = buildTemplate("service-request-pool", "base")
	serviceRequestQuick             = buildTemplate("service-request-quick", "base")
	serviceRequestQuickConfirmation = buildTemplate("service-request-quick-confirmation", "base")
	serviceRequestUpdates           = buildTemplate("service-request-updates", "base")
	settingRoot                     = buildTemplate("setting-mock", "base")
	settingPesticide                = buildTemplate("setting-pesticide", "base")
	settingPesticideAdd             = buildTemplate("setting-pesticide-add", "base")
	settingUsers                    = buildTemplate("setting-user", "base")
	settingUsersAdd                 = buildTemplate("setting-user-add", "base")
*/

type mock struct {
	Path     string
	template string
}

var mocks = []mock{}

func addMock(r chi.Router, path string, template string) {
	mocks = append(mocks, mock{
		Path:     path,
		template: template,
	})
	r.Get(path, renderMock(template))
}
func renderMock(template_name string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		code := chi.URLParam(r, "code")
		if code == "" {
			code = "abc-123"
		}
		data := ContentMock{
			DistrictName: "Delta MVCD",
			URLs: ContentMockURLs{
				Dispatch:            "/mock/dispatch",
				DispatchResults:     "/mock/dispatch-results",
				ReportConfirmation:  fmt.Sprintf("/mock/report/%s/confirm", code),
				ReportDetail:        fmt.Sprintf("/mock/report/%s", code),
				ReportContribute:    fmt.Sprintf("/mock/report/%s/contribute", code),
				ReportEvidence:      fmt.Sprintf("/mock/report/%s/evidence", code),
				ReportSchedule:      fmt.Sprintf("/mock/report/%s/schedule", code),
				ReportUpdate:        fmt.Sprintf("/mock/report/%s/update", code),
				Root:                "/mock",
				Setting:             "/mock/setting",
				SettingIntegration:  "/mock/setting/integration",
				SettingPesticide:    "/mock/setting/pesticide",
				SettingPesticideAdd: "/mock/setting/pesticide/add",
				SettingUser:         "/mock/setting/user",
				SettingUserAdd:      "/mock/setting/user/add",
			},
		}
		html.RenderOrError(w, template_name, data)
	}
}

type contentMockList struct {
	Mocks []mock
}

func renderMockList(w http.ResponseWriter, r *http.Request) {
	data := contentMockList{
		Mocks: mocks,
	}
	html.RenderOrError(w, "sync/mock/root.html", data)
}
