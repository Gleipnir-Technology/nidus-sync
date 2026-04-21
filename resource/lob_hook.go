package resource

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"time"

	/*
		"github.com/aarondl/opt/omitnull"
		"github.com/Gleipnir-Technology/nidus-sync/db/enums"
		"github.com/Gleipnir-Technology/nidus-sync/db/models"
		"github.com/Gleipnir-Technology/nidus-sync/html"
	*/
	nhttp "github.com/Gleipnir-Technology/nidus-sync/http"
	"github.com/aarondl/opt/omit"
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

/*
"id": "redacted",

	"description": "redacted",
	"name": "redacted",
	"address_line1": "redacted",
	"address_line2": "redacted",
	"address_city": "redacted",
	"address_state": "redacted",
	"address_zip": "redacted",
	"address_country": "redacted",
	"metadata": {},
	"date_created": "2026-04-21T21:43:44.819Z",
	"date_modified": "2026-04-21T21:43:44.819Z",
	"object": "redacted"
*/
type LobEventBody struct {
	AddressLine1   omit.Val[string]          `json:"address_line1"`
	AddressLine2   omit.Val[string]          `json:"address_line2"`
	AddressCity    omit.Val[string]          `json:"address_city"`
	AddressState   omit.Val[string]          `json:"address_state"`
	AddressZip     omit.Val[string]          `json:"address_zip"`
	AddressCountry omit.Val[string]          `json:"address_country"`
	Description    omit.Val[string]          `json:"description"`
	ID             omit.Val[string]          `json:"id"`
	Metadata       omit.Val[json.RawMessage] `json:"metadata"`
	DateCreated    omit.Val[time.Time]       `json:"date_created"`
	DateModified   omit.Val[time.Time]       `json:"date_modified"`
	Name           omit.Val[string]          `json:"name"`
	Object         omit.Val[string]          `json:"object"`
	Status         omit.Val[string]          `json:"status"`
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
