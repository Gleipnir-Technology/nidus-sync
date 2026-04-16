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

type Address struct{}
type ResponseAddressList struct {
	Addresses  []Address `json:"data"`
	Count      int       `json:"count"`
	CountTotal int       `json:"total_count"`
}

func (l *Lob) AddressList(ctx context.Context) ([]*Address, error) {
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
	return []*Address{}, nil
}
