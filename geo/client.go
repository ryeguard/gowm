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
	baseURL    = "http://api.openweathermap.org/geo/1.0/direct"
	qParam     = "q"
	appIDParam = "appid"
	limitParam = "limit"
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

func NewClient(opts *ClientOptions) (*Client, error) {
	client := &Client{
		baseURL: baseURL,
	}
	// Defaults if opts are not provided
	if opts == nil {
		if apiID, ok := internal.LoadEnvVar(); ok {
			client.appID = apiID
		}

		client.httpClient = http.DefaultClient
		client.logger = slog.Default()
	} else { // Otherwise use provided values
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

		if opts.Logger == nil {
			client.logger = slog.Default()
		} else {
			client.logger = opts.Logger
		}
	}
	return client, nil
}

type GeoOptions struct {
	Limit int // Number of the locations in the API response.
}

// Direct returns coordinates by location name
// The query should be the city name, state code (only for the US) and country code divided by comma. Please use ISO 3166 country codes.
func (c *Client) Direct(query string, opts *GeoOptions) (*DirectResponse, error) {
	if query == "" {
		return nil, fmt.Errorf("query argument is required")
	}

	u, err := url.Parse(c.baseURL)
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

	// Save response body to a file
	f, err := os.Create("response.json")
	if err != nil {
		return nil, fmt.Errorf("failed to create file: %w", err)
	}
	defer f.Close()

	_, err = f.Write(bodyBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to write to file: %w", err)
	}

	var directResp DirectResponse
	if err := json.Unmarshal(bodyBytes, &directResp.Data); err != nil {
		return nil, fmt.Errorf("failed to decode direct response JSON: %w", err)
	}

	return &directResp, nil
}
