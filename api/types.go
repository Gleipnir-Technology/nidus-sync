package api

import (
	"net/http"
	"time"

	"github.com/Gleipnir-Technology/nidus-sync/db/models"
	"github.com/Gleipnir-Technology/nidus-sync/h3utils"
	"github.com/Gleipnir-Technology/nidus-sync/platform"
	"github.com/aarondl/opt/null"
	"github.com/go-chi/render"
	"github.com/rs/zerolog/log"
)

type H3Cell uint64

type hasCreated interface {
	getCreated() string
}

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

/* not sure if used
type Location struct {
	Latitude  float64
	Longitude float64
}
*/

type NoteImagePayload struct {
	UUID      string     `json:"uuid"`
	Cell      H3Cell     `json:"cell"`
	Created   time.Time  `json:"created"`
	CreatorID int        `db:"creator_id"`
	Deleted   *time.Time `json:"deleted"`
	DeletorID *int32     `json:"deletor_id"`
	Version   int32      `json:"version"`
}

type NoteAudioPayload struct {
	UUID                    string                       `json:"uuid"`
	Breadcrumbs             []NoteAudioBreadcrumbPayload `json:"breadcrumbs"`
	Created                 time.Time                    `json:"created"`
	CreatorID               int                          `json:"creator_id"`
	Deleted                 *time.Time                   `json:"deleted"`
	DeletorID               *int32                       `json:"deletor_id"`
	Duration                float32                      `json:"duration"`
	Transcription           *string                      `json:"transcription"`
	TranscriptionUserEdited bool                         `json:"transcriptionUserEdited"`
	Version                 int32                        `json:"version"`
}

type ResponseDistrict struct {
	Agency  string `json:"agency"`
	Manager string `json:"manager"`
	Phone   string `json:"phone"`
	Website string `json:"website"`
}

type ResponseMosquitoSource struct {
	Access                  string                       `json:"access"`
	Active                  *bool                        `json:"active"`
	Comments                string                       `json:"comments"`
	Created                 string                       `json:"created"`
	Description             string                       `json:"description"`
	H3Cell                  int64                        `json:"h3cell"`
	ID                      string                       `json:"id"`
	LastInspectionDate      string                       `json:"last_inspection_date"`
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

type ResponseFieldseeker struct {
	MosquitoSources []ResponseMosquitoSource `json:"sources"`
	ServiceRequests []ResponseServiceRequest `json:"requests"`
	TrapData        []ResponseTrapData       `json:"traps"`
}

// ResponseErr renderer type for handling all sorts of errors.
type ResponseClientIos struct {
	Fieldseeker ResponseFieldseeker `json:"fieldseeker"`
	Since       time.Time           `json:"since"`
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
func NewResponseMosquitoInspection(i *models.FieldseekerMosquitoinspection) ResponseMosquitoInspection {
	return ResponseMosquitoInspection{
		ActionTaken:   i.Actiontaken.GetOr(""),
		Comments:      i.Comments.GetOr(""),
		Condition:     i.Sitecond.GetOr(""),
		Created:       i.Creationdate.MustGet().Format("2006-01-02T15:04:05.000Z"),
		ID:            i.Globalid.String(),
		LocationName:  i.Locationname.GetOr(""),
		SiteCondition: i.Sitecond.GetOr(""),
	}
}
func NewResponseMosquitoInspections(inspections models.FieldseekerMosquitoinspectionSlice) []ResponseMosquitoInspection {
	results := make([]ResponseMosquitoInspection, 0)
	for _, i := range inspections {
		results = append(results, NewResponseMosquitoInspection(i))
	}
	return results
}

func (rd ResponseDistrict) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (rtd ResponseMosquitoSource) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func NewResponseMosquitoSource(ms platform.MosquitoSource) ResponseMosquitoSource {
	pl := ms.PointLocation
	h3cell, err := h3utils.ToCell(pl.H3cell.GetOr("0"))
	if err != nil {
		log.Warn().Err(err).Msg("Failed to convert h3 cell")
		h3cell = 0
	}
	return ResponseMosquitoSource{
		Active:                  toBool16(pl.Active),
		Access:                  pl.Accessdesc.GetOr(""),
		Comments:                pl.Comments.GetOr(""),
		Created:                 formatTime(pl.Creationdate),
		Description:             pl.Description.GetOr(""),
		H3Cell:                  int64(h3cell),
		ID:                      pl.Globalid.String(),
		LastInspectionDate:      formatTime(pl.Lastinspectdate),
		Habitat:                 pl.Habitat.GetOr(""),
		Inspections:             NewResponseMosquitoInspections(ms.Inspections),
		Name:                    pl.Name.GetOr(""),
		NextActionDateScheduled: formatTime(pl.Nextactiondatescheduled),
		Treatments:              NewResponseMosquitoTreatments(ms.Treatments),
		UseType:                 pl.Usetype.GetOr(""),
		WaterOrigin:             pl.Waterorigin.GetOr(""),
		Zone:                    pl.Zone.GetOr(""),
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
func NewResponseMosquitoTreatment(i *models.FieldseekerTreatment) ResponseMosquitoTreatment {
	return ResponseMosquitoTreatment{
		Comments:        i.Comments.GetOr(""),
		Created:         formatTime(i.Creationdate),
		FieldTechnician: i.Fieldtech.GetOr(""),
		Habitat:         i.Habitat.GetOr(""),
		ID:              i.Globalid.String(),
		Product:         i.Product.GetOr(""),
		Quantity:        i.Qty.GetOr(0),
		QuantityUnit:    i.Qtyunit.GetOr(""),
		SiteCondition:   i.Sitecond.GetOr(""),
		TreatAcres:      i.Treatacres.GetOr(0.0),
		TreatHectares:   i.Treathectares.GetOr(0.0),
	}
}
func NewResponseMosquitoTreatments(treatments models.FieldseekerTreatmentSlice) []ResponseMosquitoTreatment {
	results := make([]ResponseMosquitoTreatment, 0)
	for _, i := range treatments {
		results = append(results, NewResponseMosquitoTreatment(i))
	}
	return results
}

type ResponseNote struct {
	CategoryName string `json:"categoryName"`
	Content      string `json:"content"`

	H3Cell    int64  `json:"h3cell"`
	ID        string `json:"id"`
	Timestamp string `json:"timestamp"`
}

func (rtd ResponseNote) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

type ResponseServiceRequest struct {
	Address            string `json:"address"`
	AssignedTechnician string `json:"assigned_technician"`
	City               string `json:"city"`
	Created            string `json:"created"`
	H3Cell             int64  `json:"h3cell"`
	HasDog             *bool  `json:"has_dog"`
	HasSpanishSpeaker  *bool  `json:"has_spanish_speaker"`
	ID                 string `json:"id"`
	Priority           string `json:"priority"`
	RecordedDate       string `json:"recorded_date"`
	Source             string `json:"source"`
	Status             string `json:"status"`
	Target             string `json:"target"`
	Zip                string `json:"zip"`
}

func (srr ResponseServiceRequest) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func NewResponseServiceRequest(sr *models.FieldseekerServicerequest) ResponseServiceRequest {
	return ResponseServiceRequest{
		Address:            sr.Reqaddr1.GetOr(""),
		AssignedTechnician: sr.Assignedtech.GetOr(""),
		City:               sr.Reqcity.GetOr(""),
		Created:            formatTime(sr.Creationdate),
		//H3Cell:             sr.H3Cell,
		HasDog:            toBool(sr.Dog),
		HasSpanishSpeaker: toBool(sr.Spanish),
		ID:                sr.Globalid.String(),
		Priority:          sr.Priority.GetOr(""),
		Status:            sr.Status.GetOr(""),
		Source:            sr.Source.GetOr(""),
		Target:            sr.Reqtarget.GetOr(""),
		Zip:               sr.Reqzip.GetOr(""),
	}
}
func NewResponseServiceRequests(requests models.FieldseekerServicerequestSlice) []ResponseServiceRequest {
	results := make([]ResponseServiceRequest, 0)
	for _, i := range requests {
		results = append(results, NewResponseServiceRequest(i))
	}
	return results
}

type ResponseTrapData struct {
	Created     string `json:"created"`
	Description string `json:"description"`
	H3Cell      int64  `json:"h3cell"`
	ID          string `json:"id"`
	Name        string `json:"name"`
}

func (rtd ResponseTrapData) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}
func NewResponseTrapDatum(td *models.FieldseekerTraplocation) ResponseTrapData {
	return ResponseTrapData{
		Created:     formatTime(td.Creationdate),
		Description: td.Description.GetOr(""),
		ID:          td.Globalid.String(),
		//H3Cell:    td.H3Cell,
		Name: td.Name.GetOr(""),
	}
}
func NewResponseTrapData(data models.FieldseekerTraplocationSlice) []ResponseTrapData {
	results := make([]ResponseTrapData, 0)
	for _, i := range data {
		results = append(results, NewResponseTrapDatum(i))
	}
	return results
}

func toResponseFieldseeker(sync platform.FieldseekerRecordsSync) ResponseFieldseeker {
	return ResponseFieldseeker{
		MosquitoSources: NewResponseMosquitoSources(sync.MosquitoSources),
		ServiceRequests: NewResponseServiceRequests(sync.ServiceRequests),
		TrapData:        NewResponseTrapData(sync.TrapData),
	}
}

func formatTime(t null.Val[time.Time]) string {
	if t.IsNull() {
		return ""
	}
	v := t.MustGet()
	return v.Format("2006-01-02T15:04:05.000Z")
}

func toBool16(t null.Val[int16]) *bool {
	if t.IsNull() {
		return nil
	}
	val := t.MustGet()
	var b bool
	if val == 0 {
		b = false
	} else {
		b = true
	}
	return &b
}
func toBool(t null.Val[int32]) *bool {
	if t.IsNull() {
		return nil
	}
	val := t.MustGet()
	var b bool
	if val == 0 {
		b = false
	} else {
		b = true
	}
	return &b
}
