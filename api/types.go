package api

import (
	"net/http"
	"sort"
	"time"
	
	"github.com/Gleipnir-Technology/nidus-sync/db"
	"github.com/go-chi/render"
	"github.com/google/uuid"
)

type H3Cell uint64

type hasCreated interface {
	getCreated() string
}

type FS_Geometry struct {
	X float64 `db:"X"`
	Y float64 `db:"Y"`
}

func (geo FS_Geometry) Latitude() float64 {
	return geo.Y
}
func (geo FS_Geometry) Longitude() float64 {
	return geo.X
}

type FS_InspectionSample struct {
	Geometry     FS_Geometry `db:"geometry"`
	CreationDate string      `db:"creationdate"`
	Creator      string      `db:"creator"`
	EditDate     string      `db:"editdate"`
	Editor       string      `db:"editor"`
	IDByTech     string      `db:"idbytech"`
	InspectionID string      `db:"insp_id"`
	Processed    int         `db:"processed"`
	SampleID     string      `db:"sampleid"`
}

type FS_MosquitoInspection struct {
	ActionTaken     *string `db:"actiontaken"`
	Comments        *string `db:"comments"`
	Condition       *string `db:"sitecond"`
	EndDateTime     string  `db:"enddatetime"`
	FieldTech       *string `db:"fieldtech"`
	GlobalID        string  `db:"globalid"`
	LocationName    *string `db:"locationname"`
	PointLocationID string  `db:"pointlocid"`
	SiteCond        *string `db:"sitecond"`
	Zone            *string `db:"zone"`
}

type FS_PointLocation struct {
	Access                  *string     `db:"accessdesc"`
	Active                  *int        `db:"active"`
	Comments                *string     `db:"comments"`
	CreationDate            *int64      `db:"creationdate"`
	Description             *string     `db:"description"`
	Geometry                FS_Geometry `db:"geometry"`
	GlobalID                string      `db:"globalid"`
	Habitat                 *string     `db:"habitat"`
	Inspections             MosquitoInspectionSlice
	LastInspectDate         *int64  `db:"lastinspectdate"`
	Name                    *string `db:"name"`
	NextActionDateScheduled *int64  `db:"nextactiondatescheduled"`
	Treatments              []MosquitoTreatment
	UseType                 *string `db:"usetype"`
	WaterOrigin             *string `db:"waterorigin"`
	Zone                    *string `db:"zone"`
}

type FS_ServiceRequest struct {
	AssignedTech *string     `db:"assignedtech"`
	CreationDate *int64      `db:"creationdate"`
	City         *string     `db:"reqcity"`
	Dog          *int        `db:"dog"`
	Geometry     FS_Geometry `db:"geometry"`
	GlobalID     string      `db:"globalid"`
	Priority     *string     `db:"priority"`
	RecDateTime  *int64      `db:"recdatetime"`
	ReqAddr1     *string     `db:"reqaddr1"`
	ReqTarget    *string     `db:"reqtarget"`
	ReqZip       *string     `db:"reqzip"`
	Source       *string     `db:"source"`
	Spanish      *int        `db:"spanish"`
	Status       *string     `db:"status"`
}
type FS_TrapLocation struct {
	Access       *string     `db:"accessdesc"`
	CreationDate *int64      `db:"creationdate"`
	Description  *string     `db:"description"`
	Geometry     FS_Geometry `db:"geometry"`
	GlobalID     string      `db:"globalid"`
	ObjectID     int         `db:"objectid"`
	Name         *string     `db:"name"`
}

type FS_Treatment struct {
	Comments        *string  `db:"comments"`
	EndDateTime     *int64   `db:"enddatetime"`
	FieldTech       *string  `db:"fieldtech"`
	GlobalID        string   `db:"globalid"`
	Habitat         *string  `db:"habitat"`
	PointLocationID string   `db:"pointlocid"`
	Product         *string  `db:"product"`
	Quantity        float64  `db:"qty"`
	QuantityUnit    *string  `db:"qtyunit"`
	SiteCondition   *string  `db:"sitecond"`
	TreatAcres      *float64 `db:"treatacres"`
	TreatHectares   *float64 `db:"treathectares"`
}

/*
type User struct {
	DisplayName      string `db:"display_name"`
	ID               int    `db:"id"`
	PasswordHashType string `db:"password_hash_type"`
	PasswordHash     string `db:"password_hash"`
	Username         string `db:"username"`
}
*/

type Bounds struct {
	East  float64
	North float64
	South float64
	West  float64
}

func NewBounds() Bounds {
	return Bounds{
		East:  180,
		North: 180,
		South: -180,
		West:  -180,
	}
}

type MosquitoInspection struct {
	data *FS_MosquitoInspection
}

func (mi MosquitoInspection) ActionTaken() string {
	if mi.data.ActionTaken == nil {
		return ""
	}
	return *mi.data.ActionTaken
}

func (mi MosquitoInspection) Comments() string {
	if mi.data.Comments == nil {
		return ""
	}
	return *mi.data.Comments
}

func (mi MosquitoInspection) Condition() string {
	if mi.data.Condition == nil {
		return ""
	}
	return *mi.data.Condition
}

func (mi MosquitoInspection) Created() time.Time {
	return parseTime(mi.data.EndDateTime)
}

func (mi MosquitoInspection) FieldTechnician() string {
	if mi.data.FieldTech == nil {
		return ""
	}
	return *mi.data.FieldTech
}

func (mi MosquitoInspection) ID() string {
	return mi.data.GlobalID
}

func (mi MosquitoInspection) LocationName() string {
	if mi.data.LocationName == nil {
		return ""
	}
	return *mi.data.LocationName
}

func (mi MosquitoInspection) SiteCondition() string {
	if mi.data.SiteCond == nil {
		return ""
	}
	return *mi.data.SiteCond
}

func NewMosquitoInspections(inspections []*FS_MosquitoInspection) []MosquitoInspection {
	results := make([]MosquitoInspection, 0)
	for _, t := range inspections {
		results = append(results, MosquitoInspection{data: t})
	}
	MosquitoInspectionSlice(results).Sort()

	return results
}

type MosquitoInspectionSlice []MosquitoInspection
type ByCreatedMI []MosquitoInspection

func (a ByCreatedMI) Len() int           { return len(a) }
func (a ByCreatedMI) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByCreatedMI) Less(i, j int) bool { return a[i].Created().After(a[j].Created()) }

func (inspections MosquitoInspectionSlice) Sort() {
	sort.Sort(ByCreatedMI(inspections))
}

type MosquitoSource struct {
	location    *FS_PointLocation
	Inspections []MosquitoInspection
	Treatments  []MosquitoTreatment
}

func (s MosquitoSource) Access() string {
	if s.location.Access == nil {
		return ""
	}
	return *s.location.Access
}

func (s MosquitoSource) Active() *bool {
	var result bool
	if s.location.Active == nil {
		return nil
	} else if *s.location.Active == 0 {
		result = false
	} else {
		result = true
	}
	return &result
}

func (s MosquitoSource) Comments() string {
	if s.location.Comments == nil {
		return ""
	}
	return *s.location.Comments
}

func (s MosquitoSource) Created() time.Time {
	if s.location.CreationDate == nil {
		return time.UnixMilli(0)
	}
	return time.UnixMilli(*s.location.CreationDate)
}

func (s MosquitoSource) Description() string {
	if s.location.Description == nil {
		return ""
	}
	return *s.location.Description
}

func (s MosquitoSource) ID() uuid.UUID {
	return uuid.MustParse(s.location.GlobalID)
}
func (s MosquitoSource) Habitat() string {
	if s.location.Habitat == nil {
		return ""
	}
	return *s.location.Habitat
}

func (s MosquitoSource) LastInspectionDate() time.Time {
	if s.location.LastInspectDate == nil {
		return time.UnixMilli(0)
	}
	return time.UnixMilli(*s.location.LastInspectDate)
}

func (s MosquitoSource) Location() LatLong {
	return s.location.Geometry
}

func (s MosquitoSource) Name() string {
	if s.location.Name == nil {
		return ""
	}
	return *s.location.Name
}

func (s MosquitoSource) NextActionDateScheduled() time.Time {
	if s.location.NextActionDateScheduled == nil {
		return time.UnixMilli(0)
	}
	return time.UnixMilli(*s.location.NextActionDateScheduled)
}

func (s MosquitoSource) UseType() string {
	if s.location.UseType == nil {
		return ""
	}
	return *s.location.UseType
}
func (s MosquitoSource) WaterOrigin() string {
	if s.location.WaterOrigin == nil {
		return ""
	}
	return *s.location.WaterOrigin
}
func (s MosquitoSource) Zone() string {
	if s.location.Zone == nil {
		return ""
	}
	return *s.location.Zone
}
func NewMosquitoSource(location *FS_PointLocation, inspections []*FS_MosquitoInspection, treatments []*FS_Treatment) MosquitoSource {
	return MosquitoSource{
		location:    location,
		Inspections: NewMosquitoInspections(inspections),
		Treatments:  NewMosquitoTreatments(treatments),
	}
}

type MosquitoTreatment struct {
	data *FS_Treatment
}

func (t MosquitoTreatment) Comments() string {
	if t.data.Comments == nil {
		return ""
	}
	return *t.data.Comments
}
func (t MosquitoTreatment) Created() time.Time {
	if t.data.EndDateTime == nil {
		return time.UnixMilli(0)
	}
	return time.UnixMilli(*t.data.EndDateTime)
}
func (t MosquitoTreatment) FieldTechnician() string {
	if t.data.FieldTech == nil {
		return ""
	}
	return *t.data.FieldTech
}
func (mi MosquitoTreatment) ID() string {
	return mi.data.GlobalID
}
func (t MosquitoTreatment) Habitat() string {
	if t.data.Habitat == nil {
		return ""
	}
	return *t.data.Habitat
}
func (t MosquitoTreatment) Product() string {
	if t.data.Product == nil {
		return ""
	}
	return *t.data.Product
}
func (t MosquitoTreatment) Quantity() float64 {
	return t.data.Quantity
}
func (t MosquitoTreatment) QuantityUnit() string {
	if t.data.QuantityUnit == nil {
		return ""
	}
	return *t.data.QuantityUnit
}
func (t MosquitoTreatment) SiteCondition() string {
	if t.data.SiteCondition == nil {
		return ""
	}
	return *t.data.SiteCondition
}
func (t MosquitoTreatment) TreatAcres() float64 {
	if t.data.TreatAcres == nil {
		return 0
	}
	return *t.data.TreatAcres
}
func (t MosquitoTreatment) TreatHectares() float64 {
	if t.data.TreatHectares == nil {
		return 0
	}
	return *t.data.TreatHectares
}
func NewMosquitoTreatments(treatments []*FS_Treatment) []MosquitoTreatment {
	results := make([]MosquitoTreatment, 0)
	for _, t := range treatments {
		results = append(results, MosquitoTreatment{data: t})
	}
	MosquitoTreatmentSlice(results).Sort()
	return results
}

type MosquitoTreatmentSlice []MosquitoTreatment
type ByCreatedMT []MosquitoTreatment

func (a ByCreatedMT) Len() int           { return len(a) }
func (a ByCreatedMT) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByCreatedMT) Less(i, j int) bool { return a[i].Created().After(a[j].Created()) }

func (inspections MosquitoTreatmentSlice) Sort() {
	sort.Sort(ByCreatedMT(inspections))
}

type LatLong interface {
	Latitude() float64
	Longitude() float64
}

type ServiceRequest struct {
	data *FS_ServiceRequest
}

func (sr ServiceRequest) Address() string {
	if sr.data.ReqAddr1 == nil {
		return ""
	}
	return *sr.data.ReqAddr1
}
func (sr ServiceRequest) AssignedTechnician() string {
	if sr.data.AssignedTech == nil {
		return ""
	}
	return *sr.data.AssignedTech
}
func (sr ServiceRequest) City() string {
	if sr.data.City == nil {
		return ""
	}
	return *sr.data.City
}
func (sr ServiceRequest) Created() time.Time {
	if sr.data.CreationDate == nil {
		return time.UnixMilli(0)
	}
	return time.UnixMilli(*sr.data.CreationDate)
}
func (sr ServiceRequest) HasDog() *bool {
	var result bool
	if sr.data.Dog == nil {
		return nil
	} else if *sr.data.Dog == 0 {
		result = false
	} else {
		result = true
	}
	return &result
}
func (sr ServiceRequest) HasSpanishSpeaker() *bool {
	var result bool
	if sr.data.Spanish == nil {
		return nil
	} else if *sr.data.Spanish == 0 {
		result = false
	} else {
		result = true
	}
	return &result
}
func (sr ServiceRequest) ID() uuid.UUID {
	return uuid.MustParse(sr.data.GlobalID)
}
func (sr ServiceRequest) Location() LatLong {
	return sr.data.Geometry
}
func (sr ServiceRequest) Priority() string {
	if sr.data.Priority == nil {
		return ""
	}
	return *sr.data.Priority
}
func (sr ServiceRequest) RecDateTime() time.Time {
	if sr.data.RecDateTime == nil {
		return time.UnixMilli(0)
	}
	return time.UnixMilli(*sr.data.RecDateTime)
}
func (sr ServiceRequest) Status() string {
	if sr.data.Status == nil {
		return ""
	}
	return *sr.data.Status
}
func (sr ServiceRequest) Source() string {
	if sr.data.Source == nil {
		return ""
	}
	return *sr.data.Source
}
func (sr ServiceRequest) Target() string {
	if sr.data.ReqTarget == nil {
		return ""
	}
	return *sr.data.ReqTarget
}
func (sr ServiceRequest) UseType() string {
	return ""
}
func (sr ServiceRequest) WaterOrigin() string {
	return ""
}
func (sr ServiceRequest) Zip() string {
	if sr.data.ReqZip == nil {
		return ""
	}
	return *sr.data.ReqZip
}
func NewServiceRequest(data *FS_ServiceRequest) ServiceRequest {
	return ServiceRequest{data: data}
}

type TrapData struct {
	data *FS_TrapLocation
}

func (tl TrapData) Access() string {
	if tl.data.Access == nil {
		return ""
	}
	return *tl.data.Access
}
func (tl TrapData) Created() time.Time {
	if tl.data.CreationDate == nil {
		return time.UnixMilli(0)
	}
	return time.UnixMilli(*tl.data.CreationDate)
}
func (tl TrapData) Description() string {
	if tl.data.Description == nil {
		return ""
	}
	return *tl.data.Description
}
func (tl TrapData) ID() uuid.UUID {
	return uuid.MustParse(tl.data.GlobalID)
}
func (tl TrapData) Location() LatLong {
	return tl.data.Geometry
}
func (tl TrapData) Name() string {
	if tl.data.Name == nil {
		return ""
	}
	return *tl.data.Name
}
func NewTrapData(data *FS_TrapLocation) TrapData {
	return TrapData{data: data}
}

type Location struct {
	Latitude  float64
	Longitude float64
}

type NoteImagePayload struct {
	UUID    string    `json:"uuid"`
	Cell    H3Cell    `json:"cell"`
	Created time.Time `json:"created"`
}

type NoteAudio struct {
	UUID                    string `db:"uuid"`
	Breadcrumbs             []NoteAudioBreadcrumbPayload
	Created                 time.Time  `db:"created"`
	Creator                 int        `db:"creator"`
	Deleted                 *time.Time `db:"deleted"`
	Duration                int        `db:"duration"`
	IsAudioNormalized       bool       `db:"is_audio_normalized"`
	IsTranscodedeToOgg      bool       `db:"is_transcoded_to_ogg"`
	Transcription           *string    `db:"transcription"`
	TranscriptionUserEdited bool       `db:"transcription_user_edited"`
	Version                 int        `db:"version"`
}

type NoteAudioPayload struct {
	UUID                    string                       `json:"uuid"`
	Breadcrumbs             []NoteAudioBreadcrumbPayload `json:"breadcrumbs"`
	Created                 time.Time                    `json:"created"`
	Duration                int                          `json:"duration"`
	Transcription           *string                      `json:"transcription"`
	TranscriptionUserEdited bool                         `json:"transcriptionUserEdited"`
	Version                 int                          `json:"version"`
}

type ResponseMosquitoSource struct {
	Access                  string                       `json:"access"`
	Active                  *bool                        `json:"active"`
	Comments                string                       `json:"comments"`
	Created                 string                       `json:"created"`
	Description             string                       `json:"description"`
	ID                      string                       `json:"id"`
	LastInspectionDate      string                       `json:"last_inspection_date"`
	Location                ResponseLocation             `json:"location"`
	Habitat                 string                       `json:"habitat"`
	Inspections             []ResponseMosquitoInspection `json:"inspections"`
	Name                    string                       `json:"name"`
	NextActionDateScheduled string                       `json:"next_action_date_scheduled"`
	Treatments              []ResponseMosquitoTreatment  `json:"treatments"`
	UseType                 string                       `json:"use_type"`
	WaterOrigin             string                       `json:"water_origin"`
	Zone                    string                       `json:"zone"`
}


type NoteAudioBreadcrumbPayload struct {
	Cell             H3Cell    `json:"cell"`
	Created          time.Time `json:"created"`
	ManuallySelected bool      `json:"manuallySelected"`
}

type NidusNotePayload struct {
	UUID      string    `json:"uuid"`
	Timestamp time.Time `json:"timestamp"`
	Images    []string  `json:"images"`
	Location  Location  `json:"location"`
	Text      string    `json:"text"`
}


// ResponseErr renderer type for handling all sorts of errors.
type ResponseClientIos struct {
	MosquitoSources []ResponseMosquitoSource `json:"sources"`
	ServiceRequests []ResponseServiceRequest `json:"requests"`
	TrapData        []ResponseTrapData       `json:"traps"`
}

func (i ResponseClientIos) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}
func NewResponseClientIos(sources []db.MosquitoSource, requests []db.ServiceRequest, traps []db.TrapData) ResponseClientIos {
	return ResponseClientIos{
		MosquitoSources: NewResponseMosquitoSources(sources),
		ServiceRequests: NewResponseServiceRequests(requests),
		TrapData:        NewResponseTrapData(traps),
	}
}

// In the best case scenario, the excellent github.com/pkg/errors package
// helps reveal information on the error, setting it on Err, and in the Render()
// method, using it to set the application-specific error code in AppCode.
type ResponseErr struct {
	Error          error `json:"-"` // low-level runtime error
	HTTPStatusCode int   `json:"-"` // http response status code

	StatusText string `json:"status"`          // user-level status message
	AppCode    int64  `json:"code,omitempty"`  // application-specific error code
	ErrorText  string `json:"error,omitempty"` // application-level error message, for debugging
}

func (e *ResponseErr) Render(w http.ResponseWriter, r *http.Request) error {
	render.Status(r, e.HTTPStatusCode)
	return nil
}

type ResponseLocation struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

func (rtd ResponseLocation) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func NewResponseLocation(l LatLong) ResponseLocation {
	return ResponseLocation{
		Latitude:  l.Latitude(),
		Longitude: l.Longitude(),
	}
}

type ResponseMosquitoInspection struct {
	ActionTaken     string `json:"action_taken"`
	Comments        string `json:"comments"`
	Condition       string `json:"condition"`
	Created         string `json:"created"`
	EndDateTime     string `json:"end_date_time"`
	FieldTechnician string `json:"field_technician"`
	ID              string `json:"id"`
	LocationName    string `json:"location_name"`
	SiteCondition   string `json:"site_condition"`
}

func (rtd ResponseMosquitoInspection) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}
func NewResponseMosquitoInspection(i MosquitoInspection) ResponseMosquitoInspection {
	return ResponseMosquitoInspection{
		ActionTaken:   i.ActionTaken(),
		Comments:      i.Comments(),
		Condition:     i.Condition(),
		Created:       i.Created().Format("2006-01-02T15:04:05.000Z"),
		ID:            i.ID(),
		LocationName:  i.LocationName(),
		SiteCondition: i.SiteCondition(),
	}
}
func NewResponseMosquitoInspections(inspections []MosquitoInspection) []ResponseMosquitoInspection {
	results := make([]ResponseMosquitoInspection, 0)
	for _, i := range inspections {
		results = append(results, NewResponseMosquitoInspection(i))
	}
	return results
}

func (rtd ResponseMosquitoSource) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func NewResponseMosquitoSource(ms db.MosquitoSource) ResponseMosquitoSource {

	return ResponseMosquitoSource{
		/*
		Active:                  ms.Active(),
		Access:                  ms.Access(),
		Comments:                ms.Comments(),
		Created:                 ms.Created().Format("2006-01-02T15:04:05.000Z"),
		Description:             ms.Description(),
		ID:                      ms.ID().String(),
		LastInspectionDate:      ms.LastInspectionDate().Format("2006-01-02T15:04:05.000Z"),
		Location:                NewResponseLocation(ms.Location()),
		Habitat:                 ms.Habitat(),
		Inspections:             NewResponseMosquitoInspections(ms.Inspections),
		Name:                    ms.Name(),
		NextActionDateScheduled: ms.NextActionDateScheduled().Format("2006-01-02T15:04:05.000Z"),
		Treatments:              NewResponseMosquitoTreatments(ms.Treatments),
		UseType:                 ms.UseType(),
		WaterOrigin:             ms.WaterOrigin(),
		Zone:                    ms.Zone(),
		*/
	}
}
func NewResponseMosquitoSources(sources []db.MosquitoSource) []ResponseMosquitoSource {
	results := make([]ResponseMosquitoSource, 0)
	for _, i := range sources {
		results = append(results, NewResponseMosquitoSource(i))
	}
	return results
}

type ResponseMosquitoTreatment struct {
	Comments        string  `json:"comments"`
	Created         string  `json:"created"`
	EndDateTime     string  `json:"end_date_time"`
	FieldTechnician string  `json:"field_technician"`
	Habitat         string  `json:"habitat"`
	ID              string  `json:"id"`
	Product         string  `json:"product"`
	Quantity        float64 `json:"quantity"`
	QuantityUnit    string  `json:"quantity_unit"`
	SiteCondition   string  `json:"site_condition"`
	TreatAcres      float64 `json:"treat_acres"`
	TreatHectares   float64 `json:"treat_hectares"`
}

func (rtd ResponseMosquitoTreatment) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}
func NewResponseMosquitoTreatment(i db.MosquitoTreatment) ResponseMosquitoTreatment {
	return ResponseMosquitoTreatment{
		/*
		Comments:        i.Comments(),
		Created:         i.Created().Format("2006-01-02T15:04:05.000Z"),
		FieldTechnician: i.FieldTechnician(),
		Habitat:         i.Habitat(),
		ID:              i.ID(),
		Product:         i.Product(),
		Quantity:        i.Quantity(),
		QuantityUnit:    i.QuantityUnit(),
		SiteCondition:   i.SiteCondition(),
		TreatAcres:      i.TreatAcres(),
		TreatHectares:   i.TreatHectares(),
		*/
	}
}
func NewResponseMosquitoTreatments(treatments []db.MosquitoTreatment) []ResponseMosquitoTreatment {
	results := make([]ResponseMosquitoTreatment, 0)
	for _, i := range treatments {
		results = append(results, NewResponseMosquitoTreatment(i))
	}
	return results
}

type ResponseNote struct {
	CategoryName string `json:"categoryName"`
	Content      string `json:"content"`

	ID        string           `json:"id"`
	Location  ResponseLocation `json:"location"`
	Timestamp string           `json:"timestamp"`
}

func (rtd ResponseNote) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

type ResponseServiceRequest struct {
	Address            string           `json:"address"`
	AssignedTechnician string           `json:"assigned_technician"`
	City               string           `json:"city"`
	Created            string           `json:"created"`
	HasDog             *bool            `json:"has_dog"`
	HasSpanishSpeaker  *bool            `json:"has_spanish_speaker"`
	ID                 string           `json:"id"`
	Location           ResponseLocation `json:"location"`
	Priority           string           `json:"priority"`
	RecordedDate       string           `json:"recorded_date"`
	Source             string           `json:"source"`
	Status             string           `json:"status"`
	Target             string           `json:"target"`
	Zip                string           `json:"zip"`
}

func (srr ResponseServiceRequest) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func NewResponseServiceRequest(sr db.ServiceRequest) ResponseServiceRequest {
	return ResponseServiceRequest{
		/*
		Address:            sr.Address(),
		AssignedTechnician: sr.AssignedTechnician(),
		City:               sr.City(),
		Created:            sr.Created().Format("2006-01-02T15:04:05.000Z"),
		HasDog:             sr.HasDog(),
		HasSpanishSpeaker:  sr.HasSpanishSpeaker(),
		ID:                 sr.ID().String(),
		Location:           NewResponseLocation(sr.Location()),
		Priority:           sr.Priority(),
		Status:             sr.Status(),
		Source:             sr.Source(),
		Target:             sr.Target(),
		Zip:                sr.Zip(),
		*/
	}
}
func NewResponseServiceRequests(requests []db.ServiceRequest) []ResponseServiceRequest {
	results := make([]ResponseServiceRequest, 0)
	for _, i := range requests {
		results = append(results, NewResponseServiceRequest(i))
	}
	return results
}

type ResponseTrapData struct {
	Created     string           `json:"created"`
	Description string           `json:"description"`
	ID          string           `json:"id"`
	Location    ResponseLocation `json:"location"`
	Name        string           `json:"name"`
}

func (rtd ResponseTrapData) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}
func NewResponseTrapDatum(td db.TrapData) ResponseTrapData {
	return ResponseTrapData{
		/*
		Created:     td.Created.Format("2006-01-02T15:04:05.000Z"),
		Description: td.Description,
		ID:          td.ID.String(),
		Location:    NewResponseLocation(td.Location),
		Name:        td.Name,
		*/
	}
}
func NewResponseTrapData(data []db.TrapData) []ResponseTrapData {
	results := make([]ResponseTrapData, 0)
	for _, i := range data {
		results = append(results, NewResponseTrapDatum(i))
	}
	return results
}
