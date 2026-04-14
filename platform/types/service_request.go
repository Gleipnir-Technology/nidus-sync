package types

import (
	"net/http"
	"time"

	"github.com/Gleipnir-Technology/nidus-sync/db/models"
	"github.com/aarondl/opt/null"
	//"github.com/google/uuid"
)

type ServiceRequest struct {
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

func ServiceRequestFromModel(sr *models.FieldseekerServicerequest) ServiceRequest {
	//log.Debug().Int32("id", m.ID).Float64("lat", m.LocationLatitude.GetOr(0.0)).Float64("lng", m.LocationLongitude.GetOr(0.0)).Msg("converting address")
	return ServiceRequest{
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
func (srr ServiceRequest) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func formatTime(t null.Val[time.Time]) string {
	if t.IsNull() {
		return ""
	}
	v := t.MustGet()
	return v.Format("2006-01-02T15:04:05.000Z")
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
