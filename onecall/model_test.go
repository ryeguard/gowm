package onecall

import (
	"bytes"
	"encoding/json"
	"log"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestParse(t *testing.T) {
	now := time.Now()

	raw := OneCallResponseRaw{
		oneCallResponseCommon: oneCallResponseCommon{
			Lat:            1,
			Lon:            2,
			Timezone:       "timezone",
			TimezoneOffset: 3,
		},
		Current: CurrentResponseRaw{
			currentResponseCommon: currentResponseCommon{
				Temp:       4,
				FeelsLike:  5,
				Pressure:   6,
				Humidity:   7,
				DewPoint:   8,
				Clouds:     9,
				UVI:        10,
				Visibility: 11,
				WindSpeed:  12,
				WindGust:   ptr(13.0),
				WindDeg:    14,
			},
			Dt:      now.Add(15 * time.Second).Unix(),
			Sunrise: now.Add(16 * time.Second).Unix(),
			Sunset:  now.Add(17 * time.Second).Unix(),
			Weather: []WeatherRaw{
				{ID: 200, Main: "Thunderstorm", Description: "thunderstorm with light rain", Icon: "11d"},
			},
		},
	}
	parsedAndConverted := raw.Parse(nil).convert()
	require.NotNil(t, parsedAndConverted)
	require.Equal(t, *parsedAndConverted, raw)
}
func TestParseTestData(t *testing.T) {
	b, err := os.ReadFile("test_data/response_stockholm_metric.json")
	if err != nil {
		log.Fatalf("Failed to read file: %v\n", err)
	}

	var raw OneCallResponseRaw
	json.NewDecoder(bytes.NewBuffer(b)).Decode(&raw)

	parsedAndConverted := raw.Parse(nil)

	require.NotNil(t, parsedAndConverted)
	require.Equal(t, *parsedAndConverted.convert(), raw)
}
