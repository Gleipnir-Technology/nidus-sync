package stadia

import (
	"crypto/tls"
	"os"
	"resty.dev/v3"
	//"github.com/rs/zerolog/log"
)

type StadiaMaps struct {
	APIKey string

	client  *resty.Client
	urlBase string
}

func NewStadiaMaps(api_key string) *StadiaMaps {
	//logger := NewLogger(log.Logger)
	//r := resty.New().SetLogger(logger).SetDebug(true)
	//r := resty.New().SetDebug(true)
	r := resty.New()
	if os.Getenv("STADIA_INSECURE_SKIP_VERIFY") != "" {
		r.SetTLSClientConfig(&tls.Config{
			InsecureSkipVerify: true,
		})
	}
	return &StadiaMaps{
		APIKey:  api_key,
		client:  r,
		urlBase: "api.stadiamaps.com",
	}
}

func (s *StadiaMaps) Close() {
	s.client.Close()
}
