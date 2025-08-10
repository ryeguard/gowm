package onecall

import "time"

type currentResponseCommon struct {
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
}

// Current weather data API response
type CurrentResponseRaw struct {
	currentResponseCommon

	// Current time, Unix, UTC
	Dt int64 `json:"dt"`

	// Sunrise time, Unix, UTC. For polar areas in midnight sun and polar night periods this parameter is not returned in the response
	Sunrise int64 `json:"sunrise"`

	// Sunset time, Unix, UTC. For polar areas in midnight sun and polar night periods this parameter is not returned in the response
	Sunset int64 `json:"sunset"`

	Weather []WeatherRaw `json:"weather"`
}

type CurrentResponse struct {
	currentResponseCommon

	// Current time
	Dt time.Time

	// Sunrise time. For polar areas in midnight sun and polar night periods this parameter is not returned in the response.
	// Use time.IsZero() to distinguish if a value was returned or not.
	Sunrise time.Time

	// Sunset time. For polar areas in midnight sun and polar night periods this parameter is not returned in the response.
	// Use time.IsZero() to distinguish if a value was returned or not.
	Sunset time.Time

	Weather []Weather
}

func (c *CurrentResponseRaw) Parse() *CurrentResponse {
	return &CurrentResponse{
		currentResponseCommon: c.currentResponseCommon,
		Dt:                    time.Unix(c.Dt, 0),
		Sunrise:               time.Unix(c.Sunrise, 0),
		Sunset:                time.Unix(c.Sunset, 0),
		Weather:               weathersRaw(c.Weather).convert(),
	}
}

func (c CurrentResponse) Parse() CurrentResponseRaw {
	return CurrentResponseRaw{
		currentResponseCommon: c.currentResponseCommon,
		Dt:                    c.Dt.Unix(),
		Sunrise:               c.Sunrise.Unix(),
		Sunset:                c.Sunset.Unix(),
		Weather:               weathers(c.Weather).convert(),
	}
}
