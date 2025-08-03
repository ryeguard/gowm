package onecall

import (
	"time"

	"github.com/martinlindhe/unit"
)

type dailyResponseCommons struct {
	MoonPhase float64  `json:"moon_phase"` // Moon phase. 0 and 1 are 'new moon', 0.25 is 'first quarter moon', 0.5 is 'full moon' and 0.75 is 'last quarter moon'. The periods in between are called 'waxing crescent', 'waxing gibbous', 'waning gibbous', and 'waning crescent', respectively. Moon phase calculation algorithm: if the moon phase values between the start of the day and the end of the day have a round value (0, 0.25, 0.5, 0.75, 1.0), then this round value is taken, otherwise the average of moon phases for the start of the day and the end of the day is taken
	Summary   string   `json:"summary"`    // Human-readable description of the weather conditions for the day
	Temp      Temp     `json:"temp"`       // Units – default: kelvin, metric: Celsius, imperial: Fahrenheit.
	Pressure  int      `json:"pressure"`   // Atmospheric pressure on the sea level, hPa
	Humidity  int      `json:"humidity"`   // Humidity, %
	DewPoint  float64  `json:"dew_point"`  // Atmospheric temperature (varying according to pressure and humidity) below which water droplets begin to condense and dew can form. Units – default: kelvin, metric: Celsius, imperial: Fahrenheit.
	WindSpeed float64  `json:"wind_speed"` // Wind speed. Units – default: metre/sec, metric: metre/sec, imperial: miles/hour.
	WindGust  *float64 `json:"wind_gust"`  // (where available) Wind gust. Units – default: metre/sec, metric: metre/sec, imperial: miles/hour.
	WindDeg   int      `json:"wind_deg"`   //  Wind direction, degrees (meteorological)
	Clouds    int      `json:"clouds"`     // Cloudiness, %
	UVI       float64  `json:"uvi"`        // The maximum value of UV index for the day
	Pop       float64  `json:"pop"`        // Probability of precipitation. The values of the parameter vary between 0 and 1, where 0 is equal to 0%, 1 is equal to 100%
}

type DailyResponseRaw struct {
	dailyResponseCommons

	Dt        int64        // Time of the forecasted data, Unix, UTC
	Sunrise   int64        // Sunrise time, Unix, UTC. For polar areas in midnight sun and polar night periods this parameter is not returned in the response
	Sunset    int64        // Sunset time, Unix, UTC. For polar areas in midnight sun and polar night periods this parameter is not returned in the response
	Moonrise  int64        // The time of when the moon rises for this day, Unix, UTC
	Moonset   int64        //  The time of when the moon rises for this day, Unix, UTC
	FeelsLike FeelsLikeRaw `json:"feels_like"` // This accounts for the human perception of weather. Units – default: kelvin, metric: Celsius, imperial: Fahrenheit.
	Weather   []WeatherRaw
}

type DailyResponse struct {
	dailyResponseCommons

	Dt        time.Time
	Sunrise   time.Time
	Sunset    time.Time
	Moonrise  time.Time
	Moonset   time.Time
	FeelsLike FeelsLike `json:"feels_like"` // This accounts for the human perception of weather. Units – default: kelvin, metric: Celsius, imperial: Fahrenheit.
	Weather   []WeatherCondition

	// internal fields
	units *Unit
}

type Temp struct {
	Morn  float64 `json:"morn"`  // Morning temperature.
	Day   float64 `json:"day"`   // Day temperature.
	Eve   float64 `json:"eve"`   // Evening temperature.
	Night float64 `json:"night"` // Night temperature.
	Min   float64 `json:"min"`   //  Min daily temperature.
	Max   float64 `json:"max"`   //  Max daily temperature.
}

type FeelsLikeRaw struct {
	Morn  float64 `json:"morn"`  // Morning temperature.
	Day   float64 `json:"day"`   // Day temperature.
	Eve   float64 `json:"eve"`   // Evening temperature.
	Night float64 `json:"night"` // Night temperature.
}

type FeelsLike struct {
	Morn  unit.Temperature `json:"morn"`  // Morning temperature.
	Day   unit.Temperature `json:"day"`   // Day temperature.
	Eve   unit.Temperature `json:"eve"`   // Evening temperature.
	Night unit.Temperature `json:"night"` // Night temperature.
}

// daily.rain (where available) Precipitation volume, mm. Please note that only mm as units of measurement are available for this parameter
// daily.snow (where available) Snow volume, mm. Please note that only mm as units of measurement are available for this parameter
// daily.weather
// daily.weather.id Weather condition id
// daily.weather.main Group of weather parameters (Rain, Snow etc.)
// daily.weather.description Weather condition within the group (full list of weather conditions). Get the output in your language
// daily.weather.icon Weather icon id.

type dailyResponsesRaw []DailyResponseRaw
type dailyResponses []DailyResponse

func (r dailyResponsesRaw) Parse(units *Unit) []DailyResponse {
	var out []DailyResponse
	for _, d := range r {
		out = append(out, DailyResponse{
			dailyResponseCommons: d.dailyResponseCommons,
			Dt:                   time.Unix(d.Dt, 0),
			Sunrise:              time.Unix(d.Sunrise, 0),
			Sunset:               time.Unix(d.Sunset, 0),
			Moonrise:             time.Unix(d.Moonrise, 0),
			Moonset:              time.Unix(d.Moonset, 0),
			FeelsLike: FeelsLike{
				Morn:  convertValueToTemp(d.FeelsLike.Morn, units),
				Day:   convertValueToTemp(d.FeelsLike.Day, units),
				Eve:   convertValueToTemp(d.FeelsLike.Eve, units),
				Night: convertValueToTemp(d.FeelsLike.Night, units),
			},
			Weather: weathersRaw(d.Weather).convert(),

			units: units,
		})
	}
	return out
}

func (r dailyResponses) convert() []DailyResponseRaw {
	var out []DailyResponseRaw
	for _, d := range r {
		out = append(out, DailyResponseRaw{
			dailyResponseCommons: d.dailyResponseCommons,
			Dt:                   d.Dt.Unix(),
			Sunrise:              d.Sunrise.Unix(),
			Sunset:               d.Sunset.Unix(),
			Moonrise:             d.Moonrise.Unix(),
			Moonset:              d.Moonset.Unix(),
			FeelsLike: FeelsLikeRaw{
				Morn:  convertTempToValue(d.FeelsLike.Morn, d.units),
				Day:   convertTempToValue(d.FeelsLike.Day, d.units),
				Eve:   convertTempToValue(d.FeelsLike.Eve, d.units),
				Night: convertTempToValue(d.FeelsLike.Night, d.units),
			},
			Weather: weatherConditions(d.Weather).convert(),
		})
	}
	return out
}

func convertValueToTemp(value float64, owmUnits *Unit) unit.Temperature {
	switch owmUnits {
	case &Units.METRIC:
		return unit.FromCelsius(value)
	case &Units.IMPERIAL:
		return unit.FromFahrenheit(value)
	default:
		return unit.FromKelvin(value)
	}
}

func convertTempToValue(value unit.Temperature, owmUnits *Unit) float64 {
	switch owmUnits {
	case &Units.METRIC:
		return value.Celsius()
	case &Units.IMPERIAL:
		return value.Fahrenheit()
	default:
		return value.Kelvin()
	}
}
