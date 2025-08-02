package onecall

import "time"

type oneCallResponseCommon struct {
	Lat            float64 `json:"lat"`
	Lon            float64 `json:"lon"`
	Timezone       string  `json:"timezone"`
	TimezoneOffset int     `json:"timezone_offset"`
}

type OneCallResponseRaw struct {
	oneCallResponseCommon
	Current CurrentResponseRaw `json:"current"`
}

type CurrentResponseRaw struct {
	// Current time, Unix, UTC
	Dt int64 `json:"dt"`

	// Sunrise time, Unix, UTC. For polar areas in midnight sun and polar night periods this parameter is not returned in the response
	Sunrise int64 `json:"sunrise"`

	// Sunset time, Unix, UTC. For polar areas in midnight sun and polar night periods this parameter is not returned in the response
	Sunset int64 `json:"sunset"`

	// Temperature. Units - default: kelvin, metric: Celsius, imperial: Fahrenheit.
	Temp float64 `json:"temp"`

	// Temperature. This temperature parameter accounts for the human perception of weather. Units – default: kelvin, metric: Celsius, imperial: Fahrenheit.
	FeelsLike float64 `json:"feels_like"`

	// Atmospheric pressure at sea level, hPa
	Pressure int64 `json:"pressure"`

	// Humidity, %
	Humidity int `json:"humidity"`

	// Atmospheric temperature (varying according to pressure and humidity) below which water droplets begin to condense and dew can form. Units – default: kelvin, metric: Celsius, imperial: Fahrenheit
	DewPoint float64 `json:"dew_point"`

	// Cloudiness, %
	Clouds int `json:"clouds"`

	// Current UV index.
	UVI float64 `json:"uvi"`

	// Average visibility, metres. The maximum value of the visibility is 10 km
	Visibility int `json:"visibility"`

	// Wind speed. Wind speed. Units – default: metre/sec, metric: metre/sec, imperial: miles/hour.
	WindSpeed float64 `json:"wind_speed"`

	// (where available) Wind gust. Units – default: metre/sec, metric: metre/sec, imperial: miles/hour.
	WindGust *float64 `json:"wind_gust"`

	// Wind direction, degrees (meteorological)
	WindDeg int `json:"wind_deg"`

	Weather []WeatherRaw `json:"weather"`
}

type WeatherRaw struct {
	ID          int64  `json:"id"`
	Main        string `json:"main"`
	Description string `json:"description"`
	Icon        string `json:"icon"`
}

type Weather struct {
	Condition WeatherCondition
}

type CurrentResponse struct {
	Dt time.Time

	// Sunrise time, Unix, UTC. For polar areas in midnight sun and polar night periods this parameter is not returned in the response.
	// Use time.IsZero() to distinguish if a value was returned or not.
	Sunrise time.Time

	// Sunset time, Unix, UTC. For polar areas in midnight sun and polar night periods this parameter is not returned in the response.
	// Use time.IsZero() to distinguish if a value was returned or not.
	Sunset time.Time

	// Temperature. Units - default: kelvin, metric: Celsius, imperial: Fahrenheit.
	Temp float64

	Weather []WeatherCondition
}

type OneCallResponse struct {
	oneCallResponseCommon
	Current CurrentResponse
	Weather []WeatherCondition
}

func (c *CurrentResponseRaw) Parse() *CurrentResponse {
	var conditions []WeatherCondition
	for _, w := range c.Weather {
		c, ok := idToWeatherCondition[w.ID]
		if !ok {
			panic("unrecognized id")
		}
		conditions = append(conditions, c)
	}
	return &CurrentResponse{
		Dt:      time.Unix(c.Dt, 0),
		Sunrise: time.Unix(c.Sunrise, 0),
		Sunset:  time.Unix(c.Sunset, 0),
		Temp:    c.Temp,
		Weather: conditions,
	}
}

func (r *OneCallResponseRaw) Parse() *OneCallResponse {
	return &OneCallResponse{
		oneCallResponseCommon: r.oneCallResponseCommon,
		Current:               *r.Current.Parse(),
	}
}
