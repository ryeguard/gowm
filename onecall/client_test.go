package onecall

import (
	"log/slog"
	"net/http"
	"net/url"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestNewClient(t *testing.T) {
	var tests = []struct {
		name string
		opts *ClientOptions
	}{
		{
			name: "nil options",
			opts: nil,
		},
		{
			name: "empty options",
			opts: &ClientOptions{},
		},
		{
			name: "custom options",
			opts: &ClientOptions{
				HttpClient: &http.Client{Timeout: time.Second},
				Logger:     slog.New(&slog.TextHandler{}),
				AppID:      "123",
				Units:      Units.IMPERIAL,
			},
		},
	}

	for _, tc := range tests {
		c := NewClient(tc.opts)
		require.NotNil(t, c)
	}
}

func TestBuildURL(t *testing.T) {

	var tests = []struct {
		name         string
		lat, lon     float64
		opts         *OneCallOptions
		wantRawQuery string
	}{
		{
			name:         "no input",
			wantRawQuery: "appid=abc&lat=0.000000&lon=0.000000",
		},
		{
			name:         "lat,lon",
			lat:          59.33,
			lon:          18.06,
			wantRawQuery: "appid=abc&lat=59.330000&lon=18.060000",
		},
		{
			name: "exclude option",
			lat:  59.33,
			lon:  18.06,
			opts: &OneCallOptions{
				Exclude: []Part{Parts.HOURLY},
			},
			wantRawQuery: "appid=abc&exclude=hourly&lat=59.330000&lon=18.060000",
		},
		{
			name: "include option",
			lat:  59.33,
			lon:  18.06,
			opts: &OneCallOptions{
				Include: []Part{Parts.CURRENT, Parts.MINUTELY, Parts.HOURLY, Parts.DAILY},
			},
			wantRawQuery: "appid=abc&exclude=alerts&lat=59.330000&lon=18.060000",
		},
		{
			name: "exclude (1) overwrites include (all) option",
			lat:  59.33,
			lon:  18.06,
			opts: &OneCallOptions{
				Exclude: []Part{Parts.ALERTS},
				Include: Parts.allSlice(),
			},
			wantRawQuery: "appid=abc&exclude=alerts&lat=59.330000&lon=18.060000",
		},
		{
			name: "excluded overwrites include option",
			lat:  59.33,
			lon:  18.06,
			opts: &OneCallOptions{
				Exclude: []Part{Parts.ALERTS, Parts.CURRENT},
				Include: []Part{Parts.ALERTS, Parts.CURRENT, Parts.MINUTELY, Parts.HOURLY, Parts.DAILY},
			},
			wantRawQuery: "appid=abc&exclude=alerts%2Ccurrent&lat=59.330000&lon=18.060000",
		},
		{
			name: "unit and lang option",
			lat:  59.33,
			lon:  18.06,
			opts: &OneCallOptions{
				Units: Units.METRIC,
				Lang:  Langs.SWEDISH,
			},
			wantRawQuery: "appid=abc&lang=sv&lat=59.330000&lon=18.060000&units=metric",
		},
	}

	client := NewClient(&ClientOptions{AppID: "abc"})
	templateURL := url.URL{Scheme: "https", Host: "api.openweathermap.org", Path: "/data/3.0/onecall"}

	for _, tc := range tests {
		got, err := client.buildURL(tc.lat, tc.lon, tc.opts)
		wantURL := templateURL
		wantURL.RawQuery = tc.wantRawQuery
		require.NoError(t, err)
		require.Equal(t, &wantURL, got, tc.name)
	}
}

func TestBuildHistoricalURL(t *testing.T) {
	var tests = []struct {
		name         string
		lat, lon     float64
		dt           int64
		opts         *OneCallOptions
		wantRawQuery string
	}{
		{
			name:         "basic historical request",
			lat:          59.33,
			lon:          18.06,
			dt:           1609459200, // 2021-01-01 00:00:00 UTC
			wantRawQuery: "appid=abc&dt=1609459200&lat=59.330000&lon=18.060000",
		},
		{
			name: "historical with units and lang",
			lat:  59.33,
			lon:  18.06,
			dt:   1609459200,
			opts: &OneCallOptions{
				Units: Units.METRIC,
				Lang:  Langs.SWEDISH,
			},
			wantRawQuery: "appid=abc&dt=1609459200&lang=sv&lat=59.330000&lon=18.060000&units=metric",
		},
		{
			name: "historical with imperial units",
			lat:  40.7128,
			lon:  -74.0060,
			dt:   1643803200, // 2022-02-02 12:00:00 UTC
			opts: &OneCallOptions{
				Units: Units.IMPERIAL,
			},
			wantRawQuery: "appid=abc&dt=1643803200&lat=40.712800&lon=-74.006000&units=imperial",
		},
	}

	client := NewClient(&ClientOptions{AppID: "abc"})
	templateURL := url.URL{Scheme: "https", Host: "api.openweathermap.org", Path: "/data/3.0/onecall/timemachine"}

	for _, tc := range tests {
		got, err := client.buildHistoricalURL(tc.lat, tc.lon, tc.dt, tc.opts)
		wantURL := templateURL
		wantURL.RawQuery = tc.wantRawQuery
		require.NoError(t, err)
		require.Equal(t, &wantURL, got, tc.name)
	}
}

func TestHistoricalRaw_Validation(t *testing.T) {
	client := NewClient(&ClientOptions{AppID: "test-key"})

	tests := []struct {
		name    string
		lat     float64
		lon     float64
		dt      int64
		wantErr string
	}{
		{
			name:    "invalid lat - too high",
			lat:     91,
			lon:     0,
			dt:      1609459200,
			wantErr: "lat argument must be in range (-90; 90), is 91",
		},
		{
			name:    "invalid lat - too low",
			lat:     -91,
			lon:     0,
			dt:      1609459200,
			wantErr: "lat argument must be in range (-90; 90), is -91",
		},
		{
			name:    "invalid lon - too high",
			lat:     0,
			lon:     181,
			dt:      1609459200,
			wantErr: "lon argument must be in range (-180; 180), is 181",
		},
		{
			name:    "invalid lon - too low",
			lat:     0,
			lon:     -181,
			dt:      1609459200,
			wantErr: "lon argument must be in range (-180; 180), is -181",
		},
		{
			name:    "invalid timestamp - zero",
			lat:     59.33,
			lon:     18.06,
			dt:      0,
			wantErr: "dt (timestamp) must be a positive Unix timestamp, is 0",
		},
		{
			name:    "invalid timestamp - negative",
			lat:     59.33,
			lon:     18.06,
			dt:      -100,
			wantErr: "dt (timestamp) must be a positive Unix timestamp, is -100",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			_, err := client.HistoricalRaw(tc.lat, tc.lon, tc.dt, nil)
			require.Error(t, err)
			require.Contains(t, err.Error(), tc.wantErr)
		})
	}
}
