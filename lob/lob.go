package lob

import (
	"context"
	"crypto/tls"
	"fmt"
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
type ResponseAddressList struct {
	Addresses  []Address `json:"data"`
	Count      int       `json:"count"`
	CountTotal int       `json:"total_count"`
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
