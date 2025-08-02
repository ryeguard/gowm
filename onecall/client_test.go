package onecall

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewClient(t *testing.T) {
	client, err := NewClient(&ClientOptions{
		AppID: "TEST",
		Units: Units.METRIC,
	})
	require.NoError(t, err)
	_, err = client.OneCall(0, 0, nil)
	require.Error(t, err) // 401 Unauthorized
}
