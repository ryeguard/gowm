package onecall

type oneCallResponseCommon struct {
	Lat            float64 `json:"lat"`
	Lon            float64 `json:"lon"`
	Timezone       string  `json:"timezone"`
	TimezoneOffset int     `json:"timezone_offset"`
}

// OneCallResponseRaw is a direct mapping of what is returned from One Call API calls.
type OneCallResponseRaw struct {
	oneCallResponseCommon
	Current  CurrentResponseRaw  `json:"current"`
	Minutely []MinuteResponseRaw `json:"minutely"` // (where available)
	Hourly   []HourlyResponseRaw `json:"hourly"`
	Daily    []DailyResponseRaw  `json:"daily"`
}

// OneCallResponse is parsed from `OneCallResponseRaw` and is a more convenient, ergonomic data structure.
type OneCallResponse struct {
	oneCallResponseCommon
	Current  CurrentResponse
	Minutely []MinuteResponse
	Hourly   []HourlyResponse
	Daily    []DailyResponse
}

type RainRaw struct {
	OneH float64 `json:"1h"`
}

type SnowRaw struct {
	OneH float64 `json:"1h"`
}

func (r *RainRaw) Parse() *float64 {
	if r == nil {
		return nil
	}
	return &r.OneH
}

func convertToRain(v *float64) *RainRaw {
	if v == nil {
		return nil
	}
	return &RainRaw{OneH: *v}
}

func (s *SnowRaw) Parse() *float64 {
	if s == nil {
		return nil
	}
	return &s.OneH
}

func convertToSnow(v *float64) *SnowRaw {
	if v == nil {
		return nil
	}
	return &SnowRaw{OneH: *v}
}

type WeatherRaw struct {
	ID          int64  `json:"id"`
	Main        string `json:"main"`
	Description string `json:"description"`
	Icon        string `json:"icon"`
}

type weathersRaw []WeatherRaw

type Weather struct {
	WeatherCondition
	Icon string
}

func (w weathersRaw) convert() []Weather {
	var out []Weather
	for _, v := range w {
		wc := idToWeatherCondition[v.ID]
		out = append(out,
			Weather{
				Icon: v.Icon,
				WeatherCondition: WeatherCondition{
					Code:        wc.Code,
					Group:       wc.Group,
					Description: wc.Description,
				}},
		)
	}
	return out
}

func (r OneCallResponseRaw) Parse() *OneCallResponse {
	return &OneCallResponse{
		oneCallResponseCommon: r.oneCallResponseCommon,
		Current:               *r.Current.Parse(),
		Minutely:              minuteResponsesRaw(r.Minutely).Parse(),
		Hourly:                hourlyResponsesRaw(r.Hourly).Parse(),
		Daily:                 dailyResponsesRaw(r.Daily).Parse(),
	}
}

func (p OneCallResponse) convert() *OneCallResponseRaw {
	return &OneCallResponseRaw{
		oneCallResponseCommon: p.oneCallResponseCommon,
		Current:               p.Current.Parse(),
		Minutely:              minuteResponses(p.Minutely).convert(),
		Hourly:                hourlyResponses(p.Hourly).convert(),
		Daily:                 dailyResponses(p.Daily).convert(),
	}
}
