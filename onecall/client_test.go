package onecall

import (
	"net/url"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewClient(t *testing.T) {
	client := NewClient(&ClientOptions{
		AppID: "TEST",
		Units: Units.METRIC,
	})
	_, err := client.CurrentAndForecast(0, 0, nil)
	require.Error(t, err) // 401 Unauthorized
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
