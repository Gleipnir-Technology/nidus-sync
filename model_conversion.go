package main

import (
	"errors"
	"fmt"
	"time"

	"github.com/Gleipnir-Technology/nidus-sync/models"
	"github.com/Gleipnir-Technology/nidus-sync/sql"
	"github.com/aarondl/opt/null"
	"github.com/uber/h3-go/v4"
)

type BreedingSourceDetail struct {
	// Basic Information
	OrganizationID int32  `json:"organizationId"`
	Name           string `json:"name"`
	Description    string `json:"description"`
	LocationNumber int64  `json:"locationNumber"`
	ObjectID       int32  `json:"objectId"`
	GlobalID       string `json:"globalId"`
	ExternalID     string `json:"externalId"`

	// Status Information
	Active           bool   `json:"active"`
	DeactivateReason string `json:"deactivateReason"`
	SourceStatus     string `json:"sourceStatus"`
	Priority         string `json:"priority"`
	ScalarPriority   int64  `json:"scalarPriority"`

	// Classification
	SourceType  string `json:"sourceType"`
	Habitat     string `json:"habitat"`
	UseType     string `json:"useType"`
	WaterOrigin string `json:"waterOrigin"`
	Symbology   string `json:"symbology"`

	// Geographical Data
	LatLng            h3.LatLng `json:"latlng"`
	Zone              string    `json:"zone"`
	Zone2             string    `json:"zone2"`
	Jurisdiction      string    `json:"jurisdiction"`
	AccessDescription string    `json:"accessDescription"`

	// Inspection Data
	LarvaeInspectInterval       int16     `json:"larvaeInspectInterval"`
	LastInspectionDate          time.Time `json:"lastInspectionDate"`
	LastInspectionActivity      string    `json:"lastInspectionActivity"`
	LastInspectionActionTaken   string    `json:"lastInspectionActionTaken"`
	LastInspectionAverageLarvae float64   `json:"lastInspectionAverageLarvae"`
	LastInspectionAveragePupae  float64   `json:"lastInspectionAveragePupae"`
	LastInspectionBreeding      string    `json:"lastInspectionBreeding"`
	LastInspectionConditions    string    `json:"lastInspectionConditions"`
	LastInspectionFieldSpecies  string    `json:"lastInspectionFieldSpecies"`
	LastInspectionLifeStages    string    `json:"lastInspectionLifeStages"`

	// Treatment Data
	LastTreatmentDate         time.Time `json:"lastTreatmentDate"`
	LastTreatmentActivity     string    `json:"lastTreatmentActivity"`
	LastTreatmentProduct      string    `json:"lastTreatmentProduct"`
	LastTreatmentQuantity     float64   `json:"lastTreatmentQuantity"`
	LastTreatmentQuantityUnit string    `json:"lastTreatmentQuantityUnit"`

	// Assignment & Schedule
	AssignedTechnician      string    `json:"assignedTechnician"`
	NextActionScheduledDate time.Time `json:"nextActionScheduledDate"`

	// Metadata
	Created  time.Time `json:"created"`
	Creator  string    `json:"creator"`
	EditedAt time.Time `json:"editedAt"`
	Editor   string    `json:"editor"`
	Updated  time.Time `json:"updated"`
	Comments string    `json:"comments"`
}

type TrapNearby struct {
	Counts   []*TrapCount
	Distance string
	ID       string
}

type TrapCount struct {
	Ended   time.Time
	Females int
	ID      string
	Males   int
	Total   int
}

type TrapData struct {
	// Basic Identifiers
	OrganizationID int32  `json:"organizationId"`
	ObjectID       int32  `json:"objectId"`
	GlobalID       string `json:"globalId"`
	LocationName   string `json:"locationName"`
	LocationID     string `json:"locationId"`
	SRID           string `json:"srid"`
	Field          int64  `json:"field"`

	// Trap Information
	TrapType         string `json:"trapType"`
	TrapCondition    string `json:"trapCondition"`
	TrapActivityType string `json:"trapActivityType"`
	TrapNights       int16  `json:"trapNights"`
	Lure             string `json:"lureType"`

	// Personnel
	FieldTechnician        string `json:"fieldTechnician"`
	IdentifiedByTechnician string `json:"identifiedByTechnician"`
	SortedByTechnician     string `json:"sortedByTechnician"`

	// Timing
	StartDateTime time.Time `json:"startDateTime"`
	EndDateTime   time.Time `json:"endDateTime"`

	// Environmental Conditions
	AverageTemperature float64 `json:"averageTemperature"`
	Rainfall           float64 `json:"rainfall"`
	WindDirection      string  `json:"windDirection"`
	WindSpeed          float64 `json:"windSpeed"`
	SiteCondition      string  `json:"siteCondition"`

	// Status and Processing
	Processed     bool      `json:"processed"`
	RecordStatus  int16     `json:"recordStatus"`
	Reviewed      bool      `json:"reviewed"`
	ReviewedBy    string    `json:"reviewedBy"`
	ReviewedDate  time.Time `json:"reviewedDate"`
	GatewaySynced bool      `json:"gatewaySynced"`
	LR            bool      `json:"laboratoryReported"`
	Voltage       float64   `json:"voltage"`

	// Location Data
	GeometryX float64 `json:"geometryX"`
	GeometryY float64 `json:"geometryY"`
	Zone      string  `json:"zone"`
	Zone2     string  `json:"zone2"`

	// Vector Survey IDs
	VectorSurveyTrapDataID     string `json:"vectorSurveyTrapDataId"`
	VectorSurveyTrapLocationID string `json:"vectorSurveyTrapLocationId"`

	// Metadata
	Created        time.Time `json:"created"`
	Creator        string    `json:"creator"`
	CreatedByUser  string    `json:"createdByUser"`
	CreatedDateAlt time.Time `json:"createdDateAlt"`
	Edited         time.Time `json:"edited"`
	Editor         string    `json:"editor"`
	LastEditedDate time.Time `json:"lastEditedDate"`
	LastEditedUser string    `json:"lastEditedUser"`
	Updated        time.Time `json:"updated"`
	Comments       string    `json:"comments"`
}

type Treatment struct {
	CadenceDelta time.Duration
	Date         time.Time
	LocationID   string
	Notes        string
	Product      string
}

func toTemplateTraps(locations []sql.TrapLocationBySourceIDRow, trap_data models.FSTrapdatumSlice, counts []sql.TrapCountByLocationIDRow) ([]TrapNearby, error) {
	results := make([]TrapNearby, 0)
	count_by_trap_data_id := make(map[string]*sql.TrapCountByLocationIDRow)
	for _, c := range counts {
		count_by_trap_data_id[c.TrapdataGlobalid] = &c
	}
	counts_by_location_id := make(map[string][]*TrapCount)
	for _, td := range trap_data {
		c, ok := count_by_trap_data_id[td.Globalid]
		if !ok {
			return results, errors.New(fmt.Sprintf("Failed to find trap count for %s", td.Globalid))
		}
		if td.LocID.IsNull() {
			return results, errors.New("Got a trap data with no location ID")
		}
		loc_id := td.LocID.MustGet()
		count := &TrapCount{
			Ended:   fsToTime(td.Enddatetime),
			Females: int(c.TotalFemales.IntPart()),
			ID:      td.Globalid,
			Males:   int(c.TotalMales),
			Total:   int(c.Total.IntPart()),
		}
		counts, ok := counts_by_location_id[loc_id]
		if !ok {
			counts = []*TrapCount{count}
		} else {
			counts = append(counts, count)
		}
		counts_by_location_id[loc_id] = counts
	}
	for _, location := range locations {
		counts, ok := counts_by_location_id[location.TrapLocationGlobalid]
		if !ok {
			return results, errors.New(fmt.Sprintf("Failed to find counts for %s", location.TrapLocationGlobalid))
		}
		trap := TrapNearby{
			Counts:   counts,
			Distance: location.Distance,
			ID:       location.TrapLocationGlobalid,
		}
		results = append(results, trap)
	}
	return results, nil
}

func toTemplateTrapData(trap_data models.FSTrapdatumSlice) ([]TrapData, error) {
	var results []TrapData
	for _, r := range trap_data {
		results = append(results, TrapData{
			// Basic Identifiers
			OrganizationID: r.OrganizationID,
			ObjectID:       r.Objectid,
			GlobalID:       r.Globalid,
			LocationName:   r.Locationname.GetOr(""),
			LocationID:     r.LocID.GetOr(""),
			SRID:           r.Srid.GetOr(""),
			Field:          r.Field.GetOr(0),

			// Trap Information
			TrapType:         r.Traptype.GetOr(""),
			TrapCondition:    r.Trapcondition.GetOr(""),
			TrapActivityType: r.Trapactivitytype.GetOr(""),
			TrapNights:       r.Trapnights.GetOr(0),
			Lure:             r.Lure.GetOr(""),

			// Personnel
			FieldTechnician:        r.Fieldtech.GetOr(""),
			IdentifiedByTechnician: r.Idbytech.GetOr(""),
			SortedByTechnician:     r.Sortbytech.GetOr(""),

			// Timing
			StartDateTime: fsToTime(r.Startdatetime),
			EndDateTime:   fsToTime(r.Enddatetime),

			// Environmental Conditions
			AverageTemperature: r.Avetemp.GetOr(0),
			Rainfall:           r.Raingauge.GetOr(0),
			WindDirection:      r.Winddir.GetOr(""),
			WindSpeed:          r.Windspeed.GetOr(0),
			SiteCondition:      r.Sitecond.GetOr(""),

			// Status and Processing
			Processed:     fsIntToBool(r.Processed),
			RecordStatus:  r.Recordstatus.GetOr(0),
			Reviewed:      fsIntToBool(r.Reviewed),
			ReviewedBy:    r.Reviewedby.GetOr(""),
			ReviewedDate:  fsToTime(r.Revieweddate),
			GatewaySynced: fsIntToBool(r.Gatewaysync),
			LR:            fsIntToBool(r.LR),
			Voltage:       r.Voltage.GetOr(0),

			// Location Data
			GeometryX: r.GeometryX.GetOr(0),
			GeometryY: r.GeometryY.GetOr(0),
			Zone:      r.Zone.GetOr(""),
			Zone2:     r.Zone2.GetOr(""),

			// Vector Survey IDs
			VectorSurveyTrapDataID:     r.Vectorsurvtrapdataid.GetOr(""),
			VectorSurveyTrapLocationID: r.Vectorsurvtraplocationid.GetOr(""),

			// Metadata
			Created:        fsToTime(r.Creationdate),
			Creator:        r.Creator.GetOr(""),
			CreatedByUser:  r.CreatedUser.GetOr(""),
			CreatedDateAlt: fsToTime(r.CreatedDate),
			Edited:         fsToTime(r.Editdate),
			Editor:         r.Editor.GetOr(""),
			LastEditedDate: fsToTime(r.LastEditedDate),
			LastEditedUser: r.LastEditedUser.GetOr(""),
			Updated:        r.Updated,
			Comments:       r.Comments.GetOr(""),
		})
	}
	return results, nil
}
func toTemplateTreatment(rows models.FSTreatmentSlice) ([]Treatment, error) {
	var results []Treatment
	for _, r := range rows {
		results = append(results, Treatment{
			Date:       *fsTimestampToTime(r.Enddatetime),
			LocationID: r.Pointlocid.GetOr("none"),
			Notes:      r.Comments.GetOr("none"),
			Product:    r.Product.GetOr("none"),
		})
	}
	return results, nil
}

func toTemplateInspection(rows models.FSMosquitoinspectionSlice) ([]Inspection, error) {
	var results []Inspection
	for _, r := range rows {
		results = append(results, Inspection{
			Action:     r.Actiontaken.GetOr("none"),
			Date:       *fsTimestampToTime(r.Enddatetime),
			Notes:      r.Comments.GetOr("none"),
			Location:   r.Locationname.GetOr("none"),
			LocationID: r.Pointlocid.GetOr(""),
		})
	}
	return results, nil
}

// Helper function to convert unix timestamp to time.Time
func fsToTime(val null.Val[int64]) time.Time {
	v, ok := val.Get()
	if !ok {
		return time.UnixMilli(0)
	}
	t := time.UnixMilli(v)
	return t
}

// Helper function to convert int16 to bool
func fsIntToBool(val null.Val[int16]) bool {
	if !val.IsValue() {
		return false
	}
	b := val.MustGet() != 0
	return b
}

// toTemplateBreedingSource transforms the DB model into the display model
func toTemplateBreedingSource(source *models.FSPointlocation) *BreedingSourceDetail {
	return &BreedingSourceDetail{
		// Basic Information
		OrganizationID: source.OrganizationID,
		Name:           source.Name.MustGet(),
		Description:    source.Description.MustGet(),
		LocationNumber: source.Locationnumber.GetOr(0),
		ObjectID:       source.Objectid,
		GlobalID:       source.Globalid,
		ExternalID:     source.Externalid.GetOr(""),

		// Status Information
		Active:           fsIntToBool(source.Active),
		DeactivateReason: source.DeactivateReason.GetOr(""),
		SourceStatus:     source.Sourcestatus.GetOr(""),
		Priority:         source.Priority.GetOr(""),
		ScalarPriority:   source.Scalarpriority.GetOr(0),

		// Classification
		SourceType:  source.Stype.GetOr(""),
		Habitat:     source.Habitat.GetOr(""),
		UseType:     source.Usetype.GetOr(""),
		WaterOrigin: source.Waterorigin.GetOr(""),
		Symbology:   source.Symbology.GetOr(""),

		// Geographical Data
		LatLng: h3.LatLng{
			Lat: source.GeometryY,
			Lng: source.GeometryX,
		},
		Zone:              source.Zone.GetOr(""),
		Zone2:             source.Zone2.GetOr(""),
		Jurisdiction:      source.Jurisdiction.GetOr(""),
		AccessDescription: source.Accessdesc.GetOr(""),

		// Inspection Data
		LarvaeInspectInterval:       source.Larvinspectinterval.GetOr(0),
		LastInspectionDate:          fsToTime(source.Lastinspectdate),
		LastInspectionActivity:      source.Lastinspectactivity.GetOr(""),
		LastInspectionActionTaken:   source.Lastinspectactiontaken.GetOr(""),
		LastInspectionAverageLarvae: source.Lastinspectavglarvae.GetOr(0),
		LastInspectionAveragePupae:  source.Lastinspectavgpupae.GetOr(0),
		LastInspectionBreeding:      source.Lastinspectbreeding.GetOr(""),
		LastInspectionConditions:    source.Lastinspectconditions.GetOr(""),
		LastInspectionFieldSpecies:  source.Lastinspectfieldspecies.GetOr(""),
		LastInspectionLifeStages:    source.Lastinspectlstages.GetOr(""),

		// Treatment Data
		LastTreatmentDate:         fsToTime(source.Lasttreatdate),
		LastTreatmentActivity:     source.Lasttreatactivity.GetOr(""),
		LastTreatmentProduct:      source.Lasttreatproduct.GetOr(""),
		LastTreatmentQuantity:     source.Lasttreatqty.GetOr(0),
		LastTreatmentQuantityUnit: source.Lasttreatqtyunit.GetOr(""),

		// Assignment & Schedule
		AssignedTechnician:      source.Assignedtech.GetOr(""),
		NextActionScheduledDate: fsToTime(source.Nextactiondatescheduled),

		// Metadata
		Created:  fsToTime(source.Creationdate),
		Creator:  source.Creator.GetOr(""),
		EditedAt: fsToTime(source.Editdate),
		Editor:   source.Editor.GetOr(""),
		Updated:  source.Updated,
		Comments: source.Comments.GetOr(""),
	}
}
