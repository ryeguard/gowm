package owm

import (
	"fmt"
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

	c.OneCall = onecall.NewClient(opts)
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

	c.Geo = geo.NewClient(opts)
	return c
}

type GeoDirectOneCallResponse struct {
	GeoDirect *geo.GeoData
	OneCall   *onecall.OneCallResponse
}

// GetWeather is an opinionated convenience function that first performs geocoding and then gets the One Call forecast.
func (c *Client) GetWeather(query string, opts *onecall.OneCallOptions) (*GeoDirectOneCallResponse, error) {
	if c.OneCall == nil || c.Geo == nil {
		return nil, fmt.Errorf("both onecall and geo clients are needed")
	}

	geo, err := c.Geo.Direct(query, &geo.GeoOptions{Limit: 1})
	if err != nil {
		return nil, fmt.Errorf("geo: %w", err)
	}

	if len(geo.Data) == 0 {
		return nil, fmt.Errorf("no result matching '%v'", query)
	}

	onecall, err := c.OneCall.CurrentAndForecast(geo.Data[0].Lat, geo.Data[0].Lon, opts)
	if err != nil {
		return nil, fmt.Errorf("onecall: %w", err)
	}

	return &GeoDirectOneCallResponse{
		GeoDirect: &geo.Data[0],
		OneCall:   onecall,
	}, nil
}
