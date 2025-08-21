package onecall

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestID2WeatherCondition(t *testing.T) {
	require.Len(t, idToWeatherCondition, len(WeatherConditions.allSlice()))
}

func TestInvert(t *testing.T) {
	var tests = []struct {
		in   PartList
		want PartList
	}{
		{
			in:   nil,
			want: Parts.allSlice(),
		},
		{
			in:   Parts.allSlice(),
			want: nil,
		},
		{
			in: PartList([]Part{
				Parts.CURRENT,
				Parts.DAILY,
				Parts.ALERTS,
			}),
			want: PartList([]Part{
				Parts.MINUTELY,
				Parts.HOURLY,
			}),
		},
		{
			in: PartList([]Part{
				Parts.CURRENT,
				Parts.MINUTELY,
				Parts.HOURLY,
			}),
			want: PartList([]Part{
				Parts.DAILY,
				Parts.ALERTS,
			}),
		},
	}

	for _, tc := range tests {
		require.Equal(t, tc.want, tc.in.Invert())
	}
}

func TestAdd(t *testing.T) {
	var tests = []struct {
		name  string
		in    PartList
		toAdd PartList
		want  PartList
	}{
		{
			name: "zero-values",
		},
		{
			name:  "nil input",
			in:    nil,
			toAdd: nil,
			want:  nil,
		},
		{
			name: "add nil",
			in: PartList([]Part{
				Parts.DAILY,
			}),
			toAdd: nil,
			want: PartList([]Part{
				Parts.DAILY,
			}),
		},
		{
			name: "add zero len",
			in: PartList([]Part{
				Parts.DAILY,
			}),
			toAdd: PartList([]Part{}),
			want: PartList([]Part{
				Parts.DAILY,
			}),
		},
		{
			name: "add to nil",
			in:   nil,
			toAdd: PartList([]Part{
				Parts.DAILY,
			}),
			want: PartList([]Part{
				Parts.DAILY,
			}),
		},
		{
			name: "add existing",
			in: PartList([]Part{
				Parts.CURRENT,
				Parts.DAILY,
				Parts.ALERTS,
			}),
			toAdd: PartList([]Part{
				Parts.DAILY,
			}),
			want: PartList([]Part{
				Parts.CURRENT,
				Parts.DAILY,
				Parts.ALERTS,
			}),
		},
		{
			name: "add not existing",
			in: PartList([]Part{
				Parts.CURRENT,
				Parts.DAILY,
				Parts.ALERTS,
			}),
			toAdd: PartList([]Part{
				Parts.HOURLY,
			}),
			want: PartList([]Part{
				Parts.CURRENT,
				Parts.DAILY,
				Parts.ALERTS,
				Parts.HOURLY,
			}),
		},
	}

	for _, tc := range tests {
		require.Equal(t, tc.want, tc.in.Add(tc.toAdd), tc.name)
	}
}
