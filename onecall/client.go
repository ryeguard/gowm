package onecall

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
)

const (
	baseURL = "https://api.openweathermap.org/data/3.0/onecall"
)

type Client struct {
	baseURL    string
	appID      string
	httpClient *http.Client
	unit       Unit
}

type ClientOptions struct {
	HttpClient *http.Client
	AppID      string // Your OpenWeather API key. May also be set as environment variable.
	Units      Unit   // Units to use for the client. Overruled by unit option explicitly passed to client calls.
}

func NewClient(opts *ClientOptions) (*Client, error) {
	client := &Client{
		baseURL: baseURL,
	}
	if opts == nil || opts.HttpClient == nil {
		client.httpClient = http.DefaultClient
	}
	// Attempt to load from env var(s) if not set
	if opts == nil || opts.AppID == "" {
		if apiID, ok := loadEnvVar(); ok {
			client.appID = apiID
		} else {
			return nil, fmt.Errorf("app id is required as client options or set as environment variable")
		}
	}
	if opts != nil && opts.Units.IsValid() {
		client.unit = opts.Units
	}
	return client, nil
}

func loadEnvVar() (string, bool) {
	for _, key := range []string{"OWM_APP_ID", "OWM_API_KEY"} {
		appId, ok := os.LookupEnv(key)
		if ok {
			return appId, true
		}
	}
	return "", false
}

type OneCallOptions struct {
	Exclude []Exclude
	Units   Unit
	Lang    Lang
}

func (c *Client) OneCallRaw(lat, lon float64, opts *OneCallOptions) (*OneCallResponseRaw, error) {
	if lat < -90 || lat > 90 {
		return nil, fmt.Errorf("TODO")
	}
	if lon < -180 || lon > 180 {
		return nil, fmt.Errorf("TODO")
	}

	u, err := url.Parse(c.baseURL)
	if err != nil {
		return nil, fmt.Errorf("parse url: %w", err)
	}

	q := u.Query()
	q.Set("lat", fmt.Sprintf("%.2f", lat))
	q.Set("lon", fmt.Sprintf("%.2f", lon))
	q.Set("appid", c.appID)

	if opts != nil && len(opts.Exclude) > 0 {
		q.Set("exclude", ExcludeList(opts.Exclude).String())
	}

	if opts != nil && opts.Units.IsValid() {
		q.Set("unit", opts.Units.String())
	} else if c.unit.IsValid() {
		q.Set("unit", c.unit.String())
	}

	if opts != nil && opts.Lang.IsValid() {
		q.Set("lang", opts.Lang.String())
	}

	u.RawQuery = q.Encode()

	fmt.Println(u.String())

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

	var oneCallResp OneCallResponseRaw
	if err := json.Unmarshal(bodyBytes, &oneCallResp); err != nil {
		return nil, fmt.Errorf("failed to decode one call response JSON: %w", err)
	}

	return &oneCallResp, nil
}

func (c *Client) OneCall(lat, lon float64, opts *OneCallOptions) (*OneCallResponse, error) {
	raw, err := c.OneCallRaw(lat, lon, opts)
	if err != nil {
		return nil, err
	}

	return raw.Parse(), nil

}
