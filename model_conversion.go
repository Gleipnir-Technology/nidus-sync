package main

import (
	"github.com/Gleipnir-Technology/nidus-sync/models"
	"github.com/aarondl/opt/null"
	"time"
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
	X                 float64 `json:"x"`
	Y                 float64 `json:"y"`
	GeometryX         float64 `json:"geometryX"`
	GeometryY         float64 `json:"geometryY"`
	Zone              string  `json:"zone"`
	Zone2             string  `json:"zone2"`
	Jurisdiction      string  `json:"jurisdiction"`
	AccessDescription string  `json:"accessDescription"`

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

// ConvertToDisplayModel transforms the DB model into the display model
func ConvertToDisplayModel(source *models.FSPointlocation) *BreedingSourceDetail {
	// Helper function to convert unix timestamp to time.Time
	toTime := func(val null.Val[int64]) time.Time {
		v, ok := val.Get()
		if !ok {
			return time.UnixMilli(0)
		}
		t := time.UnixMilli(v)
		return t
	}

	// Helper function to convert int16 to bool
	toBool := func(val null.Val[int16]) bool {
		if !val.IsValue() {
			return false
		}
		b := val.MustGet() != 0
		return b
	}

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
		Active:           toBool(source.Active),
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
		X:                 source.X.GetOr(0),
		Y:                 source.Y.GetOr(0),
		GeometryX:         source.GeometryX,
		GeometryY:         source.GeometryY,
		Zone:              source.Zone.GetOr(""),
		Zone2:             source.Zone2.GetOr(""),
		Jurisdiction:      source.Jurisdiction.GetOr(""),
		AccessDescription: source.Accessdesc.GetOr(""),

		// Inspection Data
		LarvaeInspectInterval:       source.Larvinspectinterval.GetOr(0),
		LastInspectionDate:          toTime(source.Lastinspectdate),
		LastInspectionActivity:      source.Lastinspectactivity.GetOr(""),
		LastInspectionActionTaken:   source.Lastinspectactiontaken.GetOr(""),
		LastInspectionAverageLarvae: source.Lastinspectavglarvae.GetOr(0),
		LastInspectionAveragePupae:  source.Lastinspectavgpupae.GetOr(0),
		LastInspectionBreeding:      source.Lastinspectbreeding.GetOr(""),
		LastInspectionConditions:    source.Lastinspectconditions.GetOr(""),
		LastInspectionFieldSpecies:  source.Lastinspectfieldspecies.GetOr(""),
		LastInspectionLifeStages:    source.Lastinspectlstages.GetOr(""),

		// Treatment Data
		LastTreatmentDate:         toTime(source.Lasttreatdate),
		LastTreatmentActivity:     source.Lasttreatactivity.GetOr(""),
		LastTreatmentProduct:      source.Lasttreatproduct.GetOr(""),
		LastTreatmentQuantity:     source.Lasttreatqty.GetOr(0),
		LastTreatmentQuantityUnit: source.Lasttreatqtyunit.GetOr(""),

		// Assignment & Schedule
		AssignedTechnician:      source.Assignedtech.GetOr(""),
		NextActionScheduledDate: toTime(source.Nextactiondatescheduled),

		// Metadata
		Created:  toTime(source.Creationdate),
		Creator:  source.Creator.GetOr(""),
		EditedAt: toTime(source.Editdate),
		Editor:   source.Editor.GetOr(""),
		Updated:  source.Updated,
		Comments: source.Comments.GetOr(""),
	}
}
