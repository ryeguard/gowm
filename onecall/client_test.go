package onecall

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewClient(t *testing.T) {
	client := NewClient(&ClientOptions{
		AppID: "TEST",
		Units: Units.METRIC,
	})
	_, err := client.OneCall(0, 0, nil)
	require.Error(t, err) // 401 Unauthorized
}
