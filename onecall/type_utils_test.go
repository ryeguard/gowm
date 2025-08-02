package onecall

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestID2WeatherCondition(t *testing.T) {
	require.Len(t, idToWeatherCondition, len(WeatherConditions.allSlice()))
}
