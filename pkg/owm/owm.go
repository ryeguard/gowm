package owm

import (
	"log/slog"
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
	logger     *slog.Logger
}

type ClientOptions struct {
	HttpClient *http.Client
	AppID      string
	Logger     *slog.Logger
}

func NewClient(opts *ClientOptions) *Client {
	if opts == nil {
		opts = &ClientOptions{}
	}

	if opts.HttpClient == nil {
		opts.HttpClient = http.DefaultClient
	}

	if opts.Logger == nil {
		opts.Logger = slog.Default()
	}

	if opts.AppID == "" {
		if apiID, ok := internal.LoadEnvVar(); ok {
			opts.AppID = apiID
		} else {
			opts.Logger.Warn("owm client app id needed for auth is not set")
		}
	}

	return &Client{
		httpClient: opts.HttpClient,
		appID:      opts.AppID,
		logger:     opts.Logger,
	}
}

func (c *Client) WithOneCall(opts *onecall.ClientOptions) *Client {
	if opts == nil {
		opts = &onecall.ClientOptions{}
	}

	if opts.HttpClient == nil {
		opts.HttpClient = c.httpClient
	}

	if opts.AppID == "" {
		opts.AppID = c.appID
	}

	if opts.Logger == nil {
		opts.Logger = c.logger
	}

	oc, err := onecall.NewClient(opts)
	if err != nil {
		c.logger.Error("unable to create onecall client", "error", err)
		panic(err)
	}
	c.OneCall = oc
	return c
}

func (c *Client) WithGeo(opts *geo.ClientOptions) *Client {
	if opts == nil {
		opts = &geo.ClientOptions{}
	}
	if opts.HttpClient == nil {
		opts.HttpClient = c.httpClient
	}

	if opts.AppID == "" {
		opts.AppID = c.appID
	}

	if opts.Logger == nil {
		opts.Logger = c.logger
	}

	geo, err := geo.NewClient(opts)
	if err != nil {
		c.logger.Error("unable to create geo client", "error", err)
		panic(err)
	}
	c.Geo = geo
	return c
}
