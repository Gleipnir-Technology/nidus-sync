package api

import (
	"net/http"
	"time"
	
	"github.com/Gleipnir-Technology/nidus-sync/db/models"
	"github.com/Gleipnir-Technology/nidus-sync/platform"
	"github.com/go-chi/render"
)

type H3Cell uint64

type hasCreated interface {
	getCreated() string
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
}
type MosquitoTreatment struct {
}

type MosquitoSource struct {
	//location    *FS_PointLocation
	Inspections []MosquitoInspection
	Treatments  []MosquitoTreatment
}

type LatLong interface {
	Latitude() float64
	Longitude() float64
}

type ServiceRequest struct {
}

type TrapData struct {
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


type ResponseFieldseeker struct {
	MosquitoSources []ResponseMosquitoSource `json:"sources"`
	ServiceRequests []ResponseServiceRequest `json:"requests"`
	TrapData        []ResponseTrapData       `json:"traps"`
}
// ResponseErr renderer type for handling all sorts of errors.
type ResponseClientIos struct {
	Fieldseeker     ResponseFieldseeker      `json:"fieldseeker"`
}

func (i ResponseClientIos) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
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
func NewResponseMosquitoInspection(i models.FieldseekerMosquitoinspection) ResponseMosquitoInspection {
	return ResponseMosquitoInspection{
		ActionTaken:   i.Actiontaken.GetOr(""),
		Comments:      i.Comments.GetOr(""),
		Condition:     i.Sitecond.GetOr(""),
		Created:       i.Creationdate.MustGet().Format("2006-01-02T15:04:05.000Z"),
		ID:            i.Globalid.MustGet().String(),
		LocationName:  i.Locationname.GetOr(""),
		SiteCondition: i.Sitecond.GetOr(""),
	}
}
func NewResponseMosquitoInspections(inspections []models.FieldseekerMosquitoinspection) []ResponseMosquitoInspection {
	results := make([]ResponseMosquitoInspection, 0)
	for _, i := range inspections {
		results = append(results, NewResponseMosquitoInspection(i))
	}
	return results
}

func (rtd ResponseMosquitoSource) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func NewResponseMosquitoSource(ms platform.MosquitoSource) ResponseMosquitoSource {

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
func NewResponseMosquitoSources(sources []platform.MosquitoSource) []ResponseMosquitoSource {
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
func NewResponseMosquitoTreatment(i platform.MosquitoTreatment) ResponseMosquitoTreatment {
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
func NewResponseMosquitoTreatments(treatments []platform.MosquitoTreatment) []ResponseMosquitoTreatment {
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

func NewResponseServiceRequest(sr platform.ServiceRequest) ResponseServiceRequest {
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
func NewResponseServiceRequests(requests []platform.ServiceRequest) []ResponseServiceRequest {
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
func NewResponseTrapDatum(td platform.TrapData) ResponseTrapData {
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
func NewResponseTrapData(data []platform.TrapData) []ResponseTrapData {
	results := make([]ResponseTrapData, 0)
	for _, i := range data {
		results = append(results, NewResponseTrapDatum(i))
	}
	return results
}

func toResponseFieldseeker(csync platform.ClientSync) ResponseFieldseeker {
	return ResponseFieldseeker{
	}
}
