package tomtom

import (
	"os"

	"resty.dev/v3"
)

type TomTom struct {
	APIKey string

	client  *resty.Client
	urlBase string
}

func NewClient() *TomTom {
	//logger := NewLogger(log.Logger)
	//r := resty.New().SetLogger(logger).SetDebug(true)
	r := resty.New().SetDebug(true)
	api_key := os.Getenv("TOMTOM_API_KEY")
	//r := resty.New()
	return &TomTom{
		APIKey:  api_key,
		client:  r,
		urlBase: "api.tomtom.com",
	}
}

func (s *TomTom) Close() {
	s.client.Close()
}
