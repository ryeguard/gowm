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
	Minutely []MinuteResponseRaw `json:"minutely"`
	Daily    []DailyResponseRaw  `json:"daily"`
}

// OneCallResponse is parsed from `OneCallResponseRaw` and is a more convenient, ergonomic data structure.
type OneCallResponse struct {
	oneCallResponseCommon
	Current  CurrentResponse
	Minutely []MinuteResponse
	Daily    []DailyResponse
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
		Daily:                 dailyResponsesRaw(r.Daily).Parse(),
	}
}

func (p OneCallResponse) convert() *OneCallResponseRaw {
	return &OneCallResponseRaw{
		oneCallResponseCommon: p.oneCallResponseCommon,
		Current:               p.Current.Parse(),
		Minutely:              minuteResponses(p.Minutely).convert(),
		Daily:                 dailyResponses(p.Daily).convert(),
	}
}
