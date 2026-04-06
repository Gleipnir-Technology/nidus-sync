package stadia

import (
	"crypto/tls"
	"github.com/rs/zerolog/log"
	"os"
	"resty.dev/v3"
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
		log.Warn().Msg("Using insecure TLS verification settings")
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

func (s *StadiaMaps) AddResponseMiddleware(m resty.ResponseMiddleware) {
	s.client.SetResponseBodyUnlimitedReads(true)
	s.client.AddResponseMiddleware(m)
}
func (s *StadiaMaps) Close() {
	s.client.Close()
}
