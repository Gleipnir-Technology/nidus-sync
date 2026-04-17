package lob

import (
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"os"

	"github.com/rs/zerolog/log"
	"resty.dev/v3"
)

type Lob struct {
	APIKey string

	client     *resty.Client
	urlBaseApi string
}

func NewLob(api_key string) *Lob {
	r := resty.New()
	if os.Getenv("LOB_INSECURE_SKIP_VERIFY") != "" {
		log.Warn().Msg("Using insecure TLS verification settings")
		r.SetTLSClientConfig(&tls.Config{
			InsecureSkipVerify: true,
		})
	}
	r.SetBasicAuth(api_key, "")
	l := &Lob{
		APIKey:     api_key,
		client:     r,
		urlBaseApi: "api.lob.com",
	}
	return l
}

type Address struct {
	ID             string                 `json:"id"`
	Description    string                 `json:"description"`
	Name           string                 `json:"name"`
	Company        string                 `json:"company"`
	Phone          *string                `json:"phone"`
	Email          *string                `json:"email"`
	AddressLine1   string                 `json:"address_line1"`
	AddressLine2   *string                `json:"address_line2"`
	AddressCity    string                 `json:"address_city"`
	AddressState   string                 `json:"address_state"`
	AddressZip     string                 `json:"address_zip"`
	AddressCountry string                 `json:"address_country"`
	Metadata       map[string]interface{} `json:"metadata"`
	DateCreated    string                 `json:"date_created"`
	DateModified   string                 `json:"date_modified"`
	RecipientMoved bool                   `json:"recipient_moved"`
	Object         string                 `json:"object"`
}
type Letter struct {
	ID                   string                 `json:"id"`
	Description          string                 `json:"description"`
	Metadata             map[string]interface{} `json:"metadata"`
	To                   Address                `json:"to"`
	From                 Address                `json:"from"`
	Color                bool                   `json:"color"`
	DoubleSided          bool                   `json:"double_sided"`
	AddressPlacement     string                 `json:"address_placement"`
	ReturnEnvelope       bool                   `json:"return_envelope"`
	PerforatedPage       *int                   `json:"perforated_page"`
	ExtraService         string                 `json:"extra_service"`
	CustomEnvelope       *string                `json:"custom_envelope"`
	TemplateID           string                 `json:"template_id"`
	TemplateVersionID    string                 `json:"template_version_id"`
	MailType             string                 `json:"mail_type"`
	URL                  string                 `json:"url"`
	MergeVariables       map[string]interface{} `json:"merge_variables"`
	Carrier              string                 `json:"carrier"`
	TrackingNumber       string                 `json:"tracking_number"`
	TrackingEvents       []interface{}          `json:"tracking_events"`
	Thumbnails           []interface{}          `json:"thumbnails"`
	ExpectedDeliveryDate string                 `json:"expected_delivery_date"`
	DateCreated          string                 `json:"date_created"`
	DateModified         string                 `json:"date_modified"`
	SendDate             string                 `json:"send_date"`
	UseType              string                 `json:"use_type"`
	FSC                  bool                   `json:"fsc"`
	Object               string                 `json:"object"`
}
type ResponseAddressList struct {
	Addresses  []Address `json:"data"`
	Count      int       `json:"count"`
	CountTotal int       `json:"total_count"`
}
type ResponseLetterList struct {
	Letters    []Letter `json:"data"`
	Count      int      `json:"count"`
	CountTotal int      `json:"total_count"`
}

type RequestAddressCreate struct {
	AddressLine1 string `json:"address_line1"`
	AddressCity  string `json:"address_city"`
	AddressState string `json:"address_state"`
	AddressZip   string `json:"address_zip"`
	Name         string `json:"name"`
}
type RequestLetterCreate struct {
	Color   bool
	From    string
	File    io.Reader
	To      string
	UseType string
}

func (l *Lob) AddressCreate(ctx context.Context, req RequestAddressCreate) (Address, error) {
	var result Address
	resp, err := l.client.R().
		SetBody(req).
		SetContext(ctx).
		SetContentType("application/json").
		SetResult(&result).
		SetPathParam("urlBase", l.urlBaseApi).
		Post("https://{urlBase}/v1/addresses")
	if err != nil {
		return result, fmt.Errorf("address list post: %w", err)
	}
	if !resp.IsSuccess() {
		content, err := io.ReadAll(resp.Body)
		if err != nil {
			return result, fmt.Errorf("not successful, and can't read response body")
		}
		return result, fmt.Errorf("not successful: %s", string(content))
	}
	return result, nil
}
func (l *Lob) AddressList(ctx context.Context) ([]Address, error) {
	var result ResponseAddressList

	resp, err := l.client.R().
		//SetQueryParamsFromValues(query).
		SetContext(ctx).
		SetResult(&result).
		SetPathParam("urlBase", l.urlBaseApi).
		Get("https://{urlBase}/v1/addresses")
	if err != nil {
		return nil, fmt.Errorf("address list get: %w", err)
	}
	if !resp.IsSuccess() {
		return nil, fmt.Errorf("not successful")
	}
	return result.Addresses, nil
}

func (l *Lob) LetterCreate(ctx context.Context, req RequestLetterCreate) (Letter, error) {
	var result Letter
	color_str := "false"
	if req.Color {
		color_str = "true"
	}
	resp, err := l.client.R().
		SetMultipartField(
			"file",
			"content.pdf",
			"application/pdf",
			req.File,
		).
		SetMultipartFormData(map[string]string{
			"color":    color_str,
			"from":     req.From,
			"to":       req.To,
			"use_type": req.UseType,
		}).
		SetContext(ctx).
		SetResult(&result).
		SetPathParam("urlBase", l.urlBaseApi).
		Post("https://{urlBase}/v1/letters")
	if err != nil {
		return result, fmt.Errorf("letters list post: %w", err)
	}
	if !resp.IsSuccess() {
		return result, fmt.Errorf("not successful")
	}
	return result, nil
}
func (l *Lob) LetterList(ctx context.Context) ([]Letter, error) {
	var result ResponseLetterList

	resp, err := l.client.R().
		//SetQueryParamsFromValues(query).
		SetContext(ctx).
		SetResult(&result).
		SetPathParam("urlBase", l.urlBaseApi).
		Get("https://{urlBase}/v1/letters")
	if err != nil {
		return nil, fmt.Errorf("letter list get: %w", err)
	}
	if !resp.IsSuccess() {
		return nil, fmt.Errorf("not successful")
	}
	return result.Letters, nil
}
