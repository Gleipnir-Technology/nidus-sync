package sync

import (
	"github.com/uber/h3-go/v4"
)

type MapMarker struct {
	LatLng h3.LatLng
}
type ContentMockURLs struct {
	Dispatch            string
	DispatchResults     string
	ReportConfirmation  string
	ReportDetail        string
	ReportContribute    string
	ReportEvidence      string
	ReportSchedule      string
	ReportUpdate        string
	Root                string
	Setting             string
	SettingIntegration  string
	SettingPesticide    string
	SettingPesticideAdd string
	SettingUser         string
	SettingUserAdd      string
}
type ContentReportDetail struct {
	NextURL   string
	UpdateURL string
}
type ContentReportDiagnostic struct {
}

type Link struct {
	Href  string
	Title string
}
