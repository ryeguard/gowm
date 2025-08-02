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

func init() {
	idToWeatherCondition = map[int64]WeatherCondition{}
	for w := range WeatherConditions.All() {
		idToWeatherCondition[w.Code] = w
	}
}
