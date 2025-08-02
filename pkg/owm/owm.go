package owm

import (
	"net/http"

	"github.com/ryeguard/gowm/geo"
	"github.com/ryeguard/gowm/internal"
	"github.com/ryeguard/gowm/onecall"
)

type Client struct {
	OneCall *onecall.Client
	Geo     *geo.Client

	httpClient *http.Client
	appID      string
}

type ClientOptions struct {
	HttpClient *http.Client
	AppID      string
}

func NewClient(opts *ClientOptions) *Client {
	client := &Client{}

	// Defaults if opts are not provided
	if opts == nil {
		if apiID, ok := internal.LoadEnvVar(); ok {
			client.appID = apiID
		}

		client.httpClient = http.DefaultClient
	} else {
		if opts.AppID == "" {
			if apiID, ok := internal.LoadEnvVar(); ok {
				client.appID = apiID
			}
		} else {
			client.appID = opts.AppID
		}

		if opts.HttpClient == nil {
			client.httpClient = http.DefaultClient
		} else {
			client.httpClient = opts.HttpClient
		}
	}

	return client
}

func (c *Client) WithOneCall(opts *onecall.ClientOptions) *Client {
	if opts == nil {
		opts = &onecall.ClientOptions{AppID: c.appID, HttpClient: http.DefaultClient}
	} else {
		if opts.AppID == "" {
			opts.AppID = c.appID
		}

		if opts.HttpClient == nil {
			opts.HttpClient = http.DefaultClient
		}
	}

	oc, err := onecall.NewClient(opts)
	if err != nil {
		panic(err)
	}
	c.OneCall = oc
	return c
}

func (c *Client) WithGeo(opts *geo.ClientOptions) *Client {
	if opts == nil {
		opts = &geo.ClientOptions{AppID: c.appID, HttpClient: http.DefaultClient}
	} else {
		if opts.AppID == "" {
			opts.AppID = c.appID
		}

		if opts.HttpClient == nil {
			opts.HttpClient = http.DefaultClient
		}
	}

	geo, err := geo.NewClient(opts)
	if err != nil {
		panic(err)
	}
	c.Geo = geo
	return c
}
