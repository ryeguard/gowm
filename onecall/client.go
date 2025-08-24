package onecall

import (
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"os"

	"github.com/ryeguard/gowm/internal"
)

const (
	baseURL      = "https://api.openweathermap.org/data/3.0/onecall"
	latParam     = "lat"
	lonParam     = "lon"
	appIDParam   = "appid"
	excludeParam = "exclude"
	unitsParam   = "units"
	langParam    = "lang"
)

type Client struct {
	baseURL    string
	appID      string
	httpClient *http.Client
	logger     *slog.Logger
	unit       Unit
}

type ClientOptions struct {
	HttpClient *http.Client
	Logger     *slog.Logger
	AppID      string // Your OpenWeather API key. May also be set as environment variable.
	Units      Unit   // Units to use for the client. Overruled by unit option explicitly passed to client calls.
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
		if appID, ok := internal.LoadEnvVar(); ok {
			opts.AppID = appID
		}
	}

	client := &Client{
		baseURL:    baseURL,
		appID:      opts.AppID,
		httpClient: opts.HttpClient,
		logger:     opts.Logger,
		unit:       opts.Units,
	}

	if opts.Units.IsValid() {
		client.unit = opts.Units
	}
	return client
}

type OneCallOptions struct {
	Exclude    []Part // Exclude parts of the weather data from the API response. Supersedes Include'd parts.
	Include    []Part // Include parts of the weather data from the API response. Superseded by Exclude'd parts.
	Units      Unit
	Lang       Lang
	SaveAsJson string
}

func (c *Client) CurrentAndForecastRaw(lat, lon float64, opts *OneCallOptions) (*OneCallResponseRaw, error) {
	if lat < -90 || lat > 90 {
		return nil, fmt.Errorf("lat argument must be in range (-90; 90), is %v", lat)
	}
	if lon < -180 || lon > 180 {
		return nil, fmt.Errorf("lon argument must be in range (-180; 180), is %v", lon)
	}

	url, err := c.buildURL(lat, lon, opts)
	if err != nil {
		return nil, fmt.Errorf("build URL: %w", err)
	}

	resp, err := c.httpClient.Get(url.String())
	if err != nil {
		return nil, fmt.Errorf("get: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized {
		return nil, fmt.Errorf("unauthorized, make sure AppID/API key is set")
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status: %s", resp.Status)
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	if opts != nil && opts.SaveAsJson != "" {
		f, err := os.Create(opts.SaveAsJson)
		if err != nil {
			return nil, fmt.Errorf("failed to create file: %w", err)
		}
		defer f.Close()

		_, err = f.Write(bodyBytes)
		if err != nil {
			return nil, fmt.Errorf("failed to write to file: %w", err)
		}
	}

	var oneCallResp OneCallResponseRaw
	if err := json.Unmarshal(bodyBytes, &oneCallResp); err != nil {
		return nil, fmt.Errorf("failed to decode one call response JSON: %w", err)
	}

	return &oneCallResp, nil
}

func (c *Client) CurrentAndForecast(lat, lon float64, opts *OneCallOptions) (*OneCallResponse, error) {
	raw, err := c.CurrentAndForecastRaw(lat, lon, opts)
	if err != nil {
		return nil, err
	}

	return raw.Parse(), nil
}

func (c *Client) buildURL(lat, lon float64, opts *OneCallOptions) (*url.URL, error) {
	u, err := url.Parse(c.baseURL)
	if err != nil {
		return nil, fmt.Errorf("parse url: %w", err)
	}

	q := u.Query()
	q.Set(latParam, fmt.Sprintf("%f", lat))
	q.Set(lonParam, fmt.Sprintf("%f", lon))
	q.Set(appIDParam, c.appID)

	if opts != nil {
		if len(opts.Exclude) > 0 && len(opts.Include) > 0 {
			toExclude := PartList(opts.Include).Invert().Add(opts.Exclude)
			if len(toExclude) > 0 {
				q.Set(excludeParam, toExclude.String())
			}
		} else if len(opts.Exclude) > 0 {
			q.Set(excludeParam, PartList(opts.Exclude).String())
		} else if len(opts.Include) > 0 {
			q.Set(excludeParam, PartList(opts.Include).Invert().String())
		}
	}

	if opts != nil && opts.Units.IsValid() {
		q.Set(unitsParam, opts.Units.String())
	} else if c.unit.IsValid() {
		q.Set(unitsParam, c.unit.String())
	}

	if opts != nil && opts.Lang.IsValid() {
		q.Set(langParam, opts.Lang.String())
	}

	u.RawQuery = q.Encode()

	return u, nil
}
