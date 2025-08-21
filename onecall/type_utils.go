package onecall

import (
	"slices"
	"strings"
)

type PartList []Part

func (pl PartList) String() string {
	var s []string
	for _, exclude := range pl {
		s = append(s, exclude.String())
	}
	return strings.Join(s, ",")
}

func (pl PartList) Add(parts []Part) PartList {
	for _, part := range parts {
		if slices.Contains(pl, part) {
			continue
		}
		pl = append(pl, part)
	}
	return pl
}

func (pl PartList) Invert() PartList {
	var inverted PartList
	for part := range Parts.All() {
		if slices.Contains(pl, part) {
			continue
		}
		inverted = append(inverted, part)
	}
	return inverted
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
