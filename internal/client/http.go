package client

import (
	"net/http"
	"time"

	"github.com/homocodian/404-hunter/internal/config"
)

func NewHttpClient(config *config.Config) *http.Client {
	return &http.Client{
		Timeout: 3 * time.Second,
		Transport: &http.Transport{
			ResponseHeaderTimeout: 2 * time.Second,
			TLSHandshakeTimeout:   2 * time.Second,
			MaxIdleConns:          1000,
			MaxIdleConnsPerHost:   config.Workers,
			IdleConnTimeout:       90 * time.Second,
		},
	}
}
