package http_client

import (
	"net/http"
	"time"
)

type HttpRequester interface {
	Do(req *http.Request) (*http.Response, error)
}

type HttpClient struct {
	Requester           HttpRequester
	BaseURL             string
	AuthorizationHeader string
}
type HttpClientCfg struct {
	Timeout                   int    `env:"TIMEOUT" envDefault:"30" json:"timeout"` // in seconds
	BaseURL                   string `env:"BASE_URL" envDefault:"http://localhost:8080" json:"baseUrl"`
	AuthorizationHeaderSecret string `env:"AUTHORIZATION_HEADER_SECRET" envDefault:"" json:"-"`
}

func New(cfg HttpClientCfg) *HttpClient {
	return &HttpClient{
		Requester:           &http.Client{Timeout: time.Duration(cfg.Timeout) * time.Second},
		BaseURL:             cfg.BaseURL,
		AuthorizationHeader: cfg.AuthorizationHeaderSecret,
	}
}
