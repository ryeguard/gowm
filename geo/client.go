package geo

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
	baseURL    = "http://api.openweathermap.org/geo/1.0"
	appIDParam = "appid"
	limitParam = "limit"
	qParam     = "q"
	latParam   = "lat"
	lonParam   = "lon"
	maxLimit   = 5
)

type Client struct {
	baseURL    string
	appID      string
	httpClient *http.Client
	logger     *slog.Logger
}

type ClientOptions struct {
	HttpClient *http.Client
	AppID      string // Your OpenWeather API key. May also be set as environment variable.
	Logger     *slog.Logger
}

func NewClient(opts *ClientOptions) *Client {

	// Defaults if opts are not provided
	if opts == nil {
		opts = &ClientOptions{}

	}

	if opts.AppID == "" {
		if apiID, ok := internal.LoadEnvVar(); ok {
			opts.AppID = apiID
		}
	}

	if opts.HttpClient == nil {
		opts.HttpClient = http.DefaultClient
	}

	if opts.Logger == nil {
		opts.Logger = slog.Default()
	}

	return &Client{
		baseURL:    baseURL,
		appID:      opts.AppID,
		httpClient: opts.HttpClient,
		logger:     opts.Logger,
	}
}

type GeoOptions struct {
	Limit      int // Number of the locations in the API response.
	SaveAsJson string
}

// Direct returns coordinates by location name
// The query should be the city name, state code (only for the US) and country code divided by comma. Please use ISO 3166 country codes.
func (c *Client) Direct(query string, opts *GeoOptions) (*GeoResponse, error) {
	if query == "" {
		return nil, fmt.Errorf("query argument is required")
	}

	u, err := url.Parse(c.baseURL + "/direct")
	if err != nil {
		return nil, fmt.Errorf("parse url: %w", err)
	}

	q := u.Query()
	q.Set(qParam, query)
	q.Set(appIDParam, c.appID)

	if opts != nil {
		if opts.Limit > maxLimit || opts.Limit < 0 {
			c.logger.Warn("limit out of range", "limit", opts.Limit)
		}
		q.Set(limitParam, fmt.Sprintf("%d", opts.Limit))
	}

	u.RawQuery = q.Encode()

	resp, err := c.httpClient.Get(u.String())
	if err != nil {
		return nil, fmt.Errorf("get: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status: %s", resp.Status)
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	if opts.SaveAsJson != "" {
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

	var geoResp GeoResponse
	if err := json.Unmarshal(bodyBytes, &geoResp.Data); err != nil {
		return nil, fmt.Errorf("failed to decode geo response JSON: %w", err)
	}

	return &geoResp, nil
}

// Reverse geocoding allows to get name of the location (city name or area name) by using geographical coordinates (lat, lon).
// The limit parameter in the API call allows you to cap how many location names you will see in the API response.
func (c *Client) Reverse(lat, lon float64, opts *GeoOptions) (*GeoResponse, error) {
	if lat < -90 || lat > 90 {
		return nil, fmt.Errorf("lat argument must be in range (-90; 90), is %v", lat)
	}
	if lon < -180 || lon > 180 {
		return nil, fmt.Errorf("lon argument must be in range (-180; 180), is %v", lon)
	}

	u, err := url.Parse(c.baseURL + "/reverse")
	if err != nil {
		return nil, fmt.Errorf("parse url: %w", err)
	}

	q := u.Query()
	q.Set(latParam, fmt.Sprintf("%f", lat))
	q.Set(lonParam, fmt.Sprintf("%f", lon))
	q.Set(appIDParam, c.appID)

	if opts != nil {
		if opts.Limit > maxLimit || opts.Limit < 0 {
			c.logger.Warn("limit out of range", "limit", opts.Limit)
		}
		q.Set(limitParam, fmt.Sprintf("%d", opts.Limit))
	}

	u.RawQuery = q.Encode()

	resp, err := c.httpClient.Get(u.String())
	if err != nil {
		return nil, fmt.Errorf("get: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status: %s", resp.Status)
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	if opts.SaveAsJson != "" {
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

	var geoResp GeoResponse
	if err := json.Unmarshal(bodyBytes, &geoResp.Data); err != nil {
		return nil, fmt.Errorf("failed to decode geo response JSON: %w", err)
	}

	return &geoResp, nil
}
