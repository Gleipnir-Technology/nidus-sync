package resource

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"time"

	/*
		"github.com/Gleipnir-Technology/nidus-sync/db/enums"
		"github.com/Gleipnir-Technology/nidus-sync/db/models"
		"github.com/Gleipnir-Technology/nidus-sync/html"
	*/
	nhttp "github.com/Gleipnir-Technology/nidus-sync/http"
	"github.com/aarondl/opt/omit"
	"github.com/aarondl/opt/omitnull"
	/*
		"github.com/Gleipnir-Technology/nidus-sync/platform"
		"github.com/Gleipnir-Technology/nidus-sync/platform/types"
		"github.com/google/uuid"
		"github.com/gorilla/mux"
	*/
	"github.com/rs/zerolog/log"
)

func LobHook(r *router) *lobHookR {
	return &lobHookR{
		router: r,
	}
}

type lobHookR struct {
	router *router
}

type LobAddress struct {
	AddressCity    string          `json:"address_city"`
	AddressCountry string          `json:"address_country"`
	AddressLine1   string          `json:"address_line1"`
	AddressLine2   string          `json:"address_line2"`
	AddressState   string          `json:"address_state"`
	AddressZip     string          `json:"address_zip"`
	DateCreated    time.Time       `json:"date_created"`
	DateModified   time.Time       `json:"date_modified"`
	Description    string          `json:"description"`
	ID             string          `json:"id"`
	Metadata       json.RawMessage `json:"metadata"`
	Name           string          `json:"name"`
	Object         string          `json:"object"`
}
type LobEventBody struct {
	AddressCity          omit.Val[string]          `json:"address_city"`
	AddressCountry       omit.Val[string]          `json:"address_country"`
	AddressLine1         omit.Val[string]          `json:"address_line1"`
	AddressLine2         omit.Val[string]          `json:"address_line2"`
	AddressPlacement     omit.Val[string]          `json:"address_placement"`
	AddressState         omit.Val[string]          `json:"address_state"`
	AddressZip           omit.Val[string]          `json:"address_zip"`
	Carrier              omit.Val[string]          `json:"carrier"`
	Color                omit.Val[bool]            `json:"color"`
	CustomEnvelope       omitnull.Val[bool]        `json:"custom_envelope"`
	DateCreated          omit.Val[time.Time]       `json:"date_created"`
	DateModified         omit.Val[time.Time]       `json:"date_modified"`
	Description          omit.Val[string]          `json:"description"`
	DoubleSided          omit.Val[bool]            `json:"double_sided"`
	ExpectedDeliveryDate omit.Val[time.Time]       `json:"expected_delivery_date"`
	ExtraService         omitnull.Val[bool]        `json:"extra_service"`
	FailureReason        omitnull.Val[string]      `json:"failure_reason"`
	From                 omit.Val[LobAddress]      `json:"from"`
	ID                   omit.Val[string]          `json:"id"`
	IsDashboard          omit.Val[bool]            `json:"is_dashboard"`
	Metadata             omit.Val[json.RawMessage] `json:"metadata"`
	MailType             omit.Val[string]          `json:"mail_type"`
	MergeVariables       omit.Val[string]          `json:"merge_variables"`
	Name                 omit.Val[string]          `json:"name"`
	Object               omit.Val[string]          `json:"object"`
	PerforatedPage       omitnull.Val[bool]        `json:"perforated_page"`
	RawURL               omit.Val[string]          `json:"raw_url"`
	ReturnEnvelope       omit.Val[bool]            `json:"return_envelope"`
	SendDate             omit.Val[time.Time]       `json:"send_date"`
	Status               omit.Val[string]          `json:"status"`
	To                   omit.Val[LobAddress]      `json:"to"`
	TrackingNumber       omit.Val[string]          `json:"tracking_number"`
	URL                  omit.Val[string]          `json:"url"`
	USPSCampaignID       omitnull.Val[string]      `json:"usps_campaign_id"`
}
type LobEventType struct {
	ID             string `json:"id"`
	EnabledForTest bool   `json:"enabled_for_test"`
	Resource       string `json:"addresses"`
	Object         string `json:"object"`
}
type LobEvent struct {
	Body        LobEventBody `json:"body"`
	DateCreated time.Time    `json:"date_created"`
	ID          string       `json:"id"`
	Object      string       `json:"object"`
	ReferenceID string       `json:"reference_id"`
	EventType   LobEventType `json:"event_type"`
}

func (res *lobHookR) Event(ctx context.Context, w http.ResponseWriter, r *http.Request) *nhttp.ErrorWithStatus {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return nhttp.NewError("read body: %w", err)
	}
	var event LobEvent
	err = json.Unmarshal(body, &event)
	if err != nil {
		return nhttp.NewBadRequest("unmarshal json: %w", err)
	}
	log.Info().Str("method", r.Method).Str("content", string(body)).Str("id", event.ID).Msg("got lob event")
	http.Error(w, "", http.StatusNoContent)
	return nil
}
