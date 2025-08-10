package onecall

import "strings"

type ExcludeList []Exclude

func (e ExcludeList) String() string {
	var s []string
	for _, exclude := range e {
		s = append(s, exclude.String())
	}
	return strings.Join(s, ",")
}

var idToWeatherCondition map[int64]WeatherCondition

type weathers []Weather

func init() {
	idToWeatherCondition = map[int64]WeatherCondition{}
	for w := range WeatherConditions.All() {
		idToWeatherCondition[w.Code] = w
	}
}

func (w weathers) convert() []WeatherRaw {
	var out []WeatherRaw
	for _, v := range w {
		out = append(out, WeatherRaw{
			ID:          v.Code,
			Main:        v.Group,
			Description: v.Description,
			Icon:        v.Icon,
		})
	}
	return out
}

func ptr[T any](v T) *T {
	return &v
}
