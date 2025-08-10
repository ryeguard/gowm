package onecall

import "time"

type hourlyResponseCommons struct {
	Temp       float64  `json:"temp"`       // Units – default: kelvin, metric: Celsius, imperial: Fahrenheit.
	FeelsLike  float64  `json:"feels_like"` // This accounts for the human perception of weather. Units – default: kelvin, metric: Celsius, imperial: Fahrenheit.
	Pressure   int      `json:"pressure"`   // Atmospheric pressure on the sea level, hPa
	Humidity   int      `json:"humidity"`   // Humidity, %
	DewPoint   float64  `json:"dew_point"`  // Atmospheric temperature (varying according to pressure and humidity) below which water droplets begin to condense and dew can form. Units – default: kelvin, metric: Celsius, imperial: Fahrenheit.
	UVI        float64  `json:"uvi"`        // The maximum value of UV index for the day
	Clouds     int      `json:"clouds"`     // Cloudiness, %
	Visibility int      `json:"visibility"`
	WindSpeed  float64  `json:"wind_speed"` // Wind speed. Units – default: metre/sec, metric: metre/sec, imperial: miles/hour.
	WindGust   *float64 `json:"wind_gust"`  // (where available) Wind gust. Units – default: metre/sec, metric: metre/sec, imperial: miles/hour.
	WindDeg    int      `json:"wind_deg"`   //  Wind direction, degrees (meteorological)
	Pop        float64  `json:"pop"`        // Probability of precipitation. The values of the parameter vary between 0 and 1, where 0 is equal to 0%, 1 is equal to 100%
}

type HourlyResponseRaw struct {
	hourlyResponseCommons

	Dt int64 // Time of the forecasted data, Unix, UTC

	Rain    *RainRaw `json:"rain,omitempty"` // (where available) Precipitation volume, mm. Please note that only mm as units of measurement are available for this parameter
	Snow    *SnowRaw `json:"snow,omitempty"` // (where available) Snow volume, mm. Please note that only mm as units of measurement are available for this parameter
	Weather []WeatherRaw
}

type HourlyResponse struct {
	hourlyResponseCommons

	Dt      time.Time
	Weather []Weather

	Rain1H *float64
	Snow1H *float64
}

type hourlyResponsesRaw []HourlyResponseRaw
type hourlyResponses []HourlyResponse

func (r hourlyResponsesRaw) Parse() []HourlyResponse {
	var out []HourlyResponse
	for _, h := range r {
		out = append(out, HourlyResponse{
			hourlyResponseCommons: h.hourlyResponseCommons,
			Dt:                    time.Unix(h.Dt, 0),
			Rain1H:                h.Rain.Parse(),
			Snow1H:                h.Snow.Parse(),
			Weather:               weathersRaw(h.Weather).convert(),
		})
	}
	return out
}

func (r hourlyResponses) convert() []HourlyResponseRaw {
	var out []HourlyResponseRaw
	for _, h := range r {
		out = append(out, HourlyResponseRaw{
			hourlyResponseCommons: h.hourlyResponseCommons,
			Dt:                    h.Dt.Unix(),
			Rain:                  convertToRain(h.Rain1H),
			Snow:                  convertToSnow(h.Snow1H),
			Weather:               weathers(h.Weather).convert(),
		})
	}
	return out
}
