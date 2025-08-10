package onecall

import "time"

// Minute forecast weather data API response
type MinuteResponseRaw struct {
	Dt            int64   `json:"dt"`            // Time of the forecasted data, unix, UTC
	Precipitation float64 `json:"precipitation"` // Precipitation, mm/h. Please note that only mm/h as units of measurement are available for this parameter
}

// Minute forecast weather data API response
type MinuteResponse struct {
	Dt            time.Time `json:"dt"`            // Time of the forecasted data
	Precipitation float64   `json:"precipitation"` // Precipitation, mm/h. Please note that only mm/h as units of measurement are available for this parameter
}

type minuteResponsesRaw []MinuteResponseRaw
type minuteResponses []MinuteResponse

func (r minuteResponsesRaw) Parse() []MinuteResponse {
	var out []MinuteResponse
	for _, m := range r {
		out = append(out, MinuteResponse{Dt: time.Unix(m.Dt, 0), Precipitation: m.Precipitation})
	}
	return out
}

func (p minuteResponses) convert() []MinuteResponseRaw {
	var out []MinuteResponseRaw
	for _, m := range p {
		out = append(out, MinuteResponseRaw{Dt: m.Dt.Unix(), Precipitation: m.Precipitation})
	}
	return out
}
