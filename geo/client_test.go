package geo

import (
	"log/slog"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestNewClient(t *testing.T) {
	var tests = []struct {
		name string
		opts *ClientOptions
	}{
		{
			name: "nil options",
			opts: nil,
		},
		{
			name: "empty options",
			opts: &ClientOptions{},
		},
		{
			name: "custom options",
			opts: &ClientOptions{
				HttpClient: &http.Client{Timeout: time.Second},
				Logger:     slog.New(&slog.TextHandler{}),
				AppID:      "123",
			},
		},
	}

	for _, tc := range tests {
		c := NewClient(tc.opts)
		require.NotNil(t, c)
	}
}
