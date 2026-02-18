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
	api_key := os.Getenv("TOMTOM_API_KEY")
	r := resty.New()
	return &TomTom{
		APIKey:  api_key,
		client:  r,
		urlBase: "api.tomtom.com",
	}
}

func (s *TomTom) Close() {
	s.client.Close()
}

func (s *TomTom) SetDebug(enabled bool) {
	s.client.Close()
	if enabled {
		s.client = resty.New().SetDebug(true)
	} else {
		s.client = resty.New().SetDebug(false)
	}
}
