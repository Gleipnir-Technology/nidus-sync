package sync

import (
	"html/template"
	"time"

	"github.com/Gleipnir-Technology/nidus-sync/notification"
	"github.com/google/uuid"
	"github.com/uber/h3-go/v4"
)

type BreedingSourceSummary struct {
	ID            uuid.UUID
	Type          string
	LastInspected *time.Time
	LastTreated   *time.Time
}

type MapMarker struct {
	LatLng h3.LatLng
}
type ComponentMap struct {
	Center      h3.LatLng
	GeoJSON     interface{}
	MapboxToken string
	Markers     []MapMarker
	Zoom        int
}
type ContentAuthenticatedPlaceholder struct {
	User User
}
type ContentCell struct {
	BreedingSources []BreedingSourceSummary
	CellBoundary    h3.CellBoundary
	Inspections     []Inspection
	MapData         ComponentMap
	Treatments      []Treatment
	User            User
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
type ContentMock struct {
	DistrictName string
	URLs         ContentMockURLs
}
type ContentReportDetail struct {
	NextURL   string
	UpdateURL string
}
type ContentReportDiagnostic struct {
}
type ContentDashboard struct {
	CountInspections     int
	CountMosquitoSources int
	CountServiceRequests int
	Geo                  template.JS
	IsSyncOngoing        bool
	LastSync             *time.Time
	MapData              ComponentMap
	Org                  string
	RecentRequests       []ServiceRequestSummary
	User                 User
}

type ContentDashboardLoading struct {
	User User
}

type ContentPlaceholder struct {
}
type ContentSignin struct {
	InvalidCredentials bool
}
type ContentSignup struct{}
type ContentSource struct {
	Inspections []Inspection
	MapData     ComponentMap
	Source      *BreedingSourceDetail
	Traps       []TrapNearby
	Treatments  []Treatment
	//TreatmentCadence TreatmentCadence
	TreatmentModels []TreatmentModel
	User            User
}
type Inspection struct {
	Action     string
	Date       *time.Time
	Notes      string
	Location   string
	LocationID uuid.UUID
}
type Link struct {
	Href  string
	Title string
}
type ServiceRequestSummary struct {
	Date     time.Time
	Location string
	Status   string
}
type User struct {
	DisplayName   string
	Initials      string
	Notifications []notification.Notification
	Username      string
}
