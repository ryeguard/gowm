package owm

import (
	"net/http"

	"github.com/ryeguard/gowm/onecall"
)

type Client struct {
	OneCall *onecall.Client
}

func NewClient(httpClient *http.Client, apiKey string, opts *onecall.ClientOptions) *Client {
	return &Client{
		OneCall: onecall.NewClient(httpClient, apiKey, opts),
	}
}
