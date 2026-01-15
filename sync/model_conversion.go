package sync

import (
	"errors"
	"fmt"
	"time"

	"github.com/Gleipnir-Technology/nidus-sync/db/models"
	"github.com/Gleipnir-Technology/nidus-sync/db/sql"
	"github.com/Gleipnir-Technology/nidus-sync/h3utils"
	"github.com/aarondl/opt/null"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"github.com/uber/h3-go/v4"
)

type BreedingSourceDetail struct {
	// Basic Information
	OrganizationID int32     `json:"organizationId"`
	Name           string    `json:"name"`
	Description    string    `json:"description"`
	LocationNumber int64     `json:"locationNumber"`
	ObjectID       int64     `json:"objectId"`
	GlobalID       uuid.UUID `json:"globalId"`
	ExternalID     string    `json:"externalId"`

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
	H3Cell            h3.Cell `json:"h3cell"`
	Zone              string  `json:"zone"`
	Zone2             string  `json:"zone2"`
	Jurisdiction      string  `json:"jurisdiction"`
	AccessDescription string  `json:"accessDescription"`

	// Inspection Data
	LarvaeInspectInterval       int16      `json:"larvaeInspectInterval"`
	LastInspectionDate          *time.Time `json:"lastInspectionDate"`
	LastInspectionActivity      string     `json:"lastInspectionActivity"`
	LastInspectionActionTaken   string     `json:"lastInspectionActionTaken"`
	LastInspectionAverageLarvae float64    `json:"lastInspectionAverageLarvae"`
	LastInspectionAveragePupae  float64    `json:"lastInspectionAveragePupae"`
	LastInspectionBreeding      string     `json:"lastInspectionBreeding"`
	LastInspectionConditions    string     `json:"lastInspectionConditions"`
	LastInspectionFieldSpecies  string     `json:"lastInspectionFieldSpecies"`
	LastInspectionLifeStages    string     `json:"lastInspectionLifeStages"`

	// Treatment Data
	LastTreatmentDate         *time.Time `json:"lastTreatmentDate"`
	LastTreatmentActivity     string     `json:"lastTreatmentActivity"`
	LastTreatmentProduct      string     `json:"lastTreatmentProduct"`
	LastTreatmentQuantity     float64    `json:"lastTreatmentQuantity"`
	LastTreatmentQuantityUnit string     `json:"lastTreatmentQuantityUnit"`

	// Assignment & Schedule
	AssignedTechnician      string     `json:"assignedTechnician"`
	NextActionScheduledDate *time.Time `json:"nextActionScheduledDate"`

	// Metadata
	Created  *time.Time `json:"created"`
	Creator  string     `json:"creator"`
	EditedAt *time.Time `json:"editedAt"`
	Editor   string     `json:"editor"`
	Comments string     `json:"comments"`
}

type TrapNearby struct {
	Counts   []*TrapCount
	Distance string
	ID       uuid.UUID
}

type TrapCount struct {
	Ended   time.Time
	Females int
	ID      uuid.UUID
	Males   int
	Total   int
}

type TrapData struct {
	// Basic Identifiers
	OrganizationID int32     `json:"organizationId"`
	ObjectID       int64     `json:"objectId"`
	GlobalID       uuid.UUID `json:"globalId"`
	LocationName   string    `json:"locationName"`
	LocationID     uuid.UUID `json:"locationId"`
	SRID           uuid.UUID `json:"srid"`
	Field          int64     `json:"field"`

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
	StartDateTime *time.Time `json:"startDateTime"`
	EndDateTime   *time.Time `json:"endDateTime"`

	// Environmental Conditions
	AverageTemperature float64 `json:"averageTemperature"`
	Rainfall           float64 `json:"rainfall"`
	WindDirection      string  `json:"windDirection"`
	WindSpeed          float64 `json:"windSpeed"`
	SiteCondition      string  `json:"siteCondition"`

	// Status and Processing
	Processed     bool       `json:"processed"`
	RecordStatus  int16      `json:"recordStatus"`
	Reviewed      bool       `json:"reviewed"`
	ReviewedBy    string     `json:"reviewedBy"`
	ReviewedDate  *time.Time `json:"reviewedDate"`
	GatewaySynced bool       `json:"gatewaySynced"`
	LR            bool       `json:"laboratoryReported"`
	Voltage       float64    `json:"voltage"`

	// Location Data
	H3Cell h3.Cell `json:"h3cell"`
	Zone   string  `json:"zone"`
	Zone2  string  `json:"zone2"`

	// Vector Survey IDs
	VectorSurveyTrapDataID     string `json:"vectorSurveyTrapDataId"`
	VectorSurveyTrapLocationID string `json:"vectorSurveyTrapLocationId"`

	// Metadata
	Created        *time.Time `json:"created"`
	Creator        string     `json:"creator"`
	CreatedByUser  string     `json:"createdByUser"`
	CreatedDateAlt *time.Time `json:"createdDateAlt"`
	Edited         *time.Time `json:"edited"`
	Editor         string     `json:"editor"`
	LastEditedDate *time.Time `json:"lastEditedDate"`
	LastEditedUser string     `json:"lastEditedUser"`
	Comments       string     `json:"comments"`
}

type Trap struct {
	Active      bool
	Comments    string
	Description string
	GlobalID    uuid.UUID
}

type Treatment struct {
	CadenceDelta time.Duration
	Date         *time.Time
	LocationID   uuid.UUID
	Notes        string
	Product      string
}

func toTemplateTrap(traps models.FieldseekerTraplocationSlice) (results []Trap, err error) {
	for _, t := range traps {
		results = append(results, Trap{
			Active:      toBool16Or(t.Active, false),
			Comments:    t.Comments.GetOr(""),
			Description: t.Description.GetOr(""),
			GlobalID:    t.Globalid,
		})
	}
	return results, err
}
func toTemplateTrapsNearby(locations []sql.TrapLocationBySourceIDRow, trap_data []sql.TrapDataByLocationIDRecentRow, counts []sql.TrapCountByLocationIDRow) ([]TrapNearby, error) {
	results := make([]TrapNearby, 0)
	count_by_trap_data_id := make(map[uuid.UUID]*sql.TrapCountByLocationIDRow)
	for _, c := range counts {
		count_by_trap_data_id[c.TrapdataGlobalid] = &c
	}
	counts_by_location_id := make(map[uuid.UUID][]*TrapCount)
	for _, td := range trap_data {
		c, ok := count_by_trap_data_id[td.Globalid]
		if !ok {
			return results, errors.New(fmt.Sprintf("Failed to find trap count for %s", td.Globalid))
		}
		loc_id := td.LocID
		count := &TrapCount{
			Ended:   td.Enddatetime,
			Females: int(c.TotalFemales),
			ID:      td.Globalid,
			Males:   int(c.TotalMales),
			Total:   int(c.Total),
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

func toTemplateTrapData(trap_data models.FieldseekerTrapdatumSlice) ([]TrapData, error) {
	var results []TrapData
	for _, r := range trap_data {
		if r.H3cell.IsNull() {
			continue
		}
		cell, err := h3utils.ToCell(r.H3cell.MustGet())
		if err != nil {
			log.Error().Err(err).Msg("Failed to get location for trap data")
			continue
		}
		results = append(results, TrapData{
			// Basic Identifiers
			OrganizationID: r.OrganizationID,
			ObjectID:       r.Objectid,
			GlobalID:       r.Globalid,
			LocationName:   r.Locationname.GetOr(""),
			LocationID:     r.LocID.GetOr(uuid.UUID{}),
			SRID:           r.Srid.GetOr(uuid.UUID{}),
			Field:          int64(r.Field.GetOr(0)),

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
			StartDateTime: getTimeOrNull(r.Startdatetime),
			EndDateTime:   getTimeOrNull(r.Enddatetime),

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
			ReviewedDate:  getTimeOrNull(r.Revieweddate),
			GatewaySynced: fsIntToBool(r.Gatewaysync),
			LR:            fsIntToBool(r.LR),
			Voltage:       r.Voltage.GetOr(0),

			// Location Data
			H3Cell: cell,
			Zone:   r.Zone.GetOr(""),
			Zone2:  r.Zone2.GetOr(""),

			// Vector Survey IDs
			VectorSurveyTrapDataID:     r.Vectorsurvtrapdataid.GetOr(""),
			VectorSurveyTrapLocationID: r.Vectorsurvtraplocationid.GetOr(""),

			// Metadata
			Created:        getTimeOrNull(r.Creationdate),
			Creator:        r.Creator.GetOr(""),
			CreatedByUser:  r.CreatedUser.GetOr(""),
			CreatedDateAlt: getTimeOrNull(r.CreatedDate),
			Edited:         getTimeOrNull(r.Editdate),
			Editor:         r.Editor.GetOr(""),
			LastEditedDate: getTimeOrNull(r.LastEditedDate),
			LastEditedUser: r.LastEditedUser.GetOr(""),
			Comments:       r.Comments.GetOr(""),
		})
	}
	return results, nil
}
func toTemplateTreatment(rows models.FieldseekerTreatmentSlice) ([]Treatment, error) {
	var results []Treatment
	for _, r := range rows {
		results = append(results, Treatment{
			Date:       getTimeOrNull(r.Enddatetime),
			LocationID: r.Pointlocid.GetOr(uuid.UUID{}),
			Notes:      r.Comments.GetOr("none"),
			Product:    r.Product.GetOr("none"),
		})
	}
	return results, nil
}

func toTemplateInspection(rows models.FieldseekerMosquitoinspectionSlice) ([]Inspection, error) {
	var results []Inspection
	for _, r := range rows {
		results = append(results, Inspection{
			Action:     r.Actiontaken.GetOr("none"),
			Date:       getTimeOrNull(r.Enddatetime),
			Notes:      r.Comments.GetOr("none"),
			Location:   r.Locationname.GetOr("none"),
			LocationID: r.Pointlocid.GetOr(uuid.UUID{}),
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
func toTemplateBreedingSource(source *models.FieldseekerPointlocation) *BreedingSourceDetail {
	if source.H3cell.IsNull() {
		log.Error().Msg("h3 cell is null")
		return nil
	}
	cell, err := h3utils.ToCell(source.H3cell.MustGet())
	if err != nil {
		log.Error().Err(err).Msg("Failed to get h3 cell from point location")
		return nil
	}
	return &BreedingSourceDetail{
		// Basic Information
		OrganizationID: source.OrganizationID,
		Name:           source.Name.MustGet(),
		Description:    source.Description.MustGet(),
		LocationNumber: int64(source.Locationnumber.GetOr(0)),
		ObjectID:       source.Objectid,
		GlobalID:       source.Globalid,
		ExternalID:     source.Externalid.GetOr(""),

		// Status Information
		Active:           fsIntToBool(source.Active),
		DeactivateReason: source.DeactivateReason.GetOr(""),
		SourceStatus:     source.Sourcestatus.GetOr(""),
		Priority:         source.Priority.GetOr(""),
		ScalarPriority:   int64(source.Scalarpriority.GetOr(0)),

		// Classification
		SourceType:  source.Stype.GetOr(""),
		Habitat:     source.Habitat.GetOr(""),
		UseType:     source.Usetype.GetOr(""),
		WaterOrigin: source.Waterorigin.GetOr(""),
		Symbology:   source.Symbology.GetOr(""),

		// Geographical Data
		H3Cell:            cell,
		Zone:              source.Zone.GetOr(""),
		Zone2:             source.Zone2.GetOr(""),
		Jurisdiction:      source.Jurisdiction.GetOr(""),
		AccessDescription: source.Accessdesc.GetOr(""),

		// Inspection Data
		LarvaeInspectInterval:       source.Larvinspectinterval.GetOr(0),
		LastInspectionDate:          getTimeOrNull(source.Lastinspectdate),
		LastInspectionActivity:      source.Lastinspectactivity.GetOr(""),
		LastInspectionActionTaken:   source.Lastinspectactiontaken.GetOr(""),
		LastInspectionAverageLarvae: source.Lastinspectavglarvae.GetOr(0),
		LastInspectionAveragePupae:  source.Lastinspectavgpupae.GetOr(0),
		LastInspectionBreeding:      source.Lastinspectbreeding.GetOr(""),
		LastInspectionConditions:    source.Lastinspectconditions.GetOr(""),
		LastInspectionFieldSpecies:  source.Lastinspectfieldspecies.GetOr(""),
		LastInspectionLifeStages:    source.Lastinspectlstages.GetOr(""),

		// Treatment Data
		LastTreatmentDate:         getTimeOrNull(source.Lasttreatdate),
		LastTreatmentActivity:     source.Lasttreatactivity.GetOr(""),
		LastTreatmentProduct:      source.Lasttreatproduct.GetOr(""),
		LastTreatmentQuantity:     source.Lasttreatqty.GetOr(0),
		LastTreatmentQuantityUnit: source.Lasttreatqtyunit.GetOr(""),

		// Assignment & Schedule
		AssignedTechnician:      source.Assignedtech.GetOr(""),
		NextActionScheduledDate: getTimeOrNull(source.Nextactiondatescheduled),

		// Metadata
		Created:  getTimeOrNull(source.Creationdate),
		Creator:  source.Creator.GetOr(""),
		EditedAt: getTimeOrNull(source.Editdate),
		Editor:   source.Editor.GetOr(""),
		Comments: source.Comments.GetOr(""),
	}
}

func getTimeOrNull(v null.Val[time.Time]) *time.Time {
	if v.IsNull() {
		return nil
	}
	val := v.MustGet()
	return &val
}

func toBool16Or(t null.Val[int16], def bool) bool {
	if t.IsNull() {
		return def
	}
	val := t.MustGet()
	var b bool
	if val == 0 {
		b = false
	} else {
		b = true
	}
	return b
}
